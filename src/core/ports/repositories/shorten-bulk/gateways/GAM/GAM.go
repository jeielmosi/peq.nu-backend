package GAM

import (
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/helpers"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	firestore_shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/firestore"
	pigeonhole_shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/pigeonhole"
	shorten_bulk_gateway "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/interfaces"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type GAMShortenBulkGateway struct {
	repository shorten_bulk.ShortenBulkRepository
}

func (g *GAMShortenBulkGateway) Get(hash string) (*entities.ShortenBulkEntity, error) {
	now := time.Now()
	res, err := g.repository.Get(hash)
	if err != nil {
		return res.Entity, err
	}

	if res == nil {
		return nil, err
	}

	err = g.repository.IncrementClicks(hash, now)
	if err != nil {
		return res.Entity, errors.New("Error at count access processing")
	}

	res.Entity.Clicks += 1
	return res.Entity, nil
}

func (g *GAMShortenBulkGateway) unlock(hash string) {
	now := time.Now()
	err := g.repository.Unlock(hash, now)
	if err != nil {
		log.Fatalf("Defer unlock '%s': %s", hash, err.Error())
	}
}

func (g *GAMShortenBulkGateway) post(
	hash string,
	shortenBulk *entities.ShortenBulkEntity,
	stopFunc func(*repositories.RepositoryDTO[entities.ShortenBulkEntity]) error,
) error {
	now := time.Now()
	if shortenBulk == nil {
		return errors.New("Empty Shorten Bulk")
	}

	err := g.repository.Lock(hash, now)
	if err != nil {
		log.Printf("Lock '%s': %s", hash, err.Error())
		return err
	}
	defer g.unlock(hash)

	backup, err := g.repository.Get(hash)
	if err != nil {
		log.Printf("Get backup '%s': %s", hash, err.Error())
		return err
	}

	err = stopFunc(backup)
	if err != nil {
		log.Printf("StopFunc '%s': %s", hash, err.Error())
		return err
	}

	dto := repositories.NewRepositoryDTO(shortenBulk, true)
	if dto == nil {
		return errors.New("Empty DTO")
	}
	err = g.repository.Post(hash, *dto)
	if err == nil {
		log.Printf("Posted at '%s'", hash)
		return err
	}

	dto = backup.Update()
	err = g.repository.Post(hash, *dto)
	if err != nil {
		log.Printf("Post backup '%s': %s", hash, err.Error())
	}

	return err
}

func (g *GAMShortenBulkGateway) postAtNewHash(shortenBulk *entities.ShortenBulkEntity) (string, error) {
	const TRY_SIZE int = 37

	stopFunc := func(dto *repositories.RepositoryDTO[entities.ShortenBulkEntity]) error {
		if dto != nil {
			return errors.New("Hash: is used")
		}
		return nil
	}

	for t := 0; t < TRY_SIZE; t++ {
		hash := helpers.NewRandomHash(helpers.HASH_SIZE)
		err := g.post(hash, shortenBulk, stopFunc)
		if err == nil {
			return hash, err
		}
	}

	return "", errors.New("Internal error: Not found empty hash")
}

func (g *GAMShortenBulkGateway) postAtOldHash(
	shortenBulk *entities.ShortenBulkEntity,
) (string, error) {
	const OLDESTS_SIZE int = 101

	mp, err := g.repository.GetOldest(OLDESTS_SIZE)
	if err != nil {
		return "", err
	}

	keys := helpers.GetKeys(mp)

	lastTimestamp := ""
	for _, key := range keys {
		if mp[key] == nil {
			continue
		}
		timestamp := repository_helpers.TimeToTimestamp1e8(mp[key].CreatedAt)
		if lastTimestamp < timestamp {
			lastTimestamp = timestamp
		}
	}

	stopFunc := func(dto *repositories.RepositoryDTO[entities.ShortenBulkEntity]) error {
		if dto == nil {
			return errors.New("Empty hash")
		}

		timestamp := repository_helpers.TimeToTimestamp1e8(dto.CreatedAt)
		if timestamp > lastTimestamp {
			return errors.New("Element is not old")
		}
		return nil
	}

	rnd := rand.New(
		rand.NewSource(
			time.Now().UTC().UnixNano(),
		),
	)
	perm := rnd.Perm(len(mp))

	for _, idx := range perm {
		hash := keys[idx]
		if strings.HasPrefix(hash, "+") {
			continue
		}

		err = g.post(hash, shortenBulk, stopFunc)
		if err == nil {
			return hash, err
		}
	}

	return "", errors.New("Internal error: Hash not found")
}

func (g *GAMShortenBulkGateway) Post(shortenBulk entities.ShortenBulkEntity) (string, error) {
	hash, err := g.postAtNewHash(&shortenBulk)
	if err == nil {
		return hash, err
	}
	return g.postAtOldHash(&shortenBulk)
}

func NewGAMShortenBulkGateway(envName string) (shorten_bulk_gateway.ShortenBulkGateway, error) {
	return newGAMShortenBulkGateway(envName)
}

func newGAMShortenBulkGateway(envName string) (*GAMShortenBulkGateway, error) {
	firestore := firestore_shorten_bulk.NewShortenBulkFirestore(envName)
	var repos = &[]*shorten_bulk.ShortenBulkRepository{
		&firestore,
		//TODO: Create A.M
	}
	pigeonhole, err := pigeonhole_shorten_bulk.NewPigeonholeShortenBulkRepository(repos)
	if err != nil {
		return nil, err
	}

	return &GAMShortenBulkGateway{
		repository: pigeonhole,
	}, nil
}
