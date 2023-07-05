package GAM

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	config "github.com/jei-el/vuo.be-backend/src/config"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk_gateway "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/interfaces"
	randutil "go.step.sm/crypto/randutil"
)

const (
	getHash = "+get_getway"
)

func getGAM() (*GAMShortenBulkGateway, error) {
	config.Load()
	envName := os.Getenv("TEST_ENV")
	return newGAMShortenBulkGateway(envName)
}

func getRandomEntity() (
	*entities.ShortenBulkEntity,
	error,
) {
	rnd := rand.New(
		rand.NewSource(
			time.Now().UTC().UnixNano(),
		),
	)

	clicks := int64(rnd.Uint32())
	url, err := randutil.ASCII(101)
	if err != nil {
		return nil, err
	}

	entity := entities.NewShortenBulkEntity(url, clicks)

	return entity, nil
}

func TestGet(t *testing.T) {
	gam, err := getGAM()
	if err != nil {
		t.Errorf("Creating a time: %s", err.Error())
	}

	shorten_bulk_gateway.TestGet(getHash, gam, t)
}

func TestPost(t *testing.T) {
	gam, err := getGAM()
	if err != nil {
		t.Errorf("Creating a time: %s", err.Error())
	}

	exp, err := getRandomEntity()
	if err != nil {
		t.Errorf("Error creating a entity: %s", err.Error())
	}
	if exp == nil {
		t.Errorf("Error creating a entity: nil pointer")
	}
	shorten_bulk_gateway.TestPost(gam, *exp, t)
}

func TestPostAtOldHash(t *testing.T) {
	gam, err := getGAM()
	if err != nil {
		t.Errorf("Creating a time: %s", err.Error())
	}

	exp, err := getRandomEntity()
	if err != nil {
		t.Errorf("Error creating a entity: %s", err.Error())
	}
	if exp == nil {
		t.Errorf("Error creating a entity: nil pointer")
	}

	hash, err := gam.postAtOldHash(exp)
	if err != nil {
		t.Errorf("Test error at post entity: %s", err.Error())
		return
	}

	res, err := gam.Get(hash)
	if err != nil {
		t.Errorf("Test error at get entity at hash '%s': %s", hash, err.Error())
		return
	}

	if (exp.Clicks+1) != res.Clicks || exp.URL != res.URL {
		log.Println("res:", res)
		log.Println("exp:", exp)

		t.Errorf("Result and expected are different")
	}
}
