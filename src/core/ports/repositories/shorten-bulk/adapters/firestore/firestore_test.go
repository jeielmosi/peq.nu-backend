package firestore_shorten_bulk

import (
	"math/rand"
	"os"
	"testing"
	"time"

	config "github.com/jei-el/vuo.be-backend/src/config"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/helpers"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
	randutil "go.step.sm/crypto/randutil"
)

const (
	getHash             = "+get"
	postHash            = "+post"
	lockHash            = "+lock"
	incrementClicksHash = "+increment_clicks"
)

func getFirestore() *ShortenBulkFirestore {
	config.Load()
	envName := os.Getenv("TEST_ENV")
	return newShortenBulkFirestore(envName)
}

func getExpectedGetDTO() (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	timestamp := "0000-01-01T00:00:00"
	timeObj, err := repository_helpers.NewTimeFromTimestamp1e8(&timestamp)
	if err != nil {
		return nil, err
	}
	locked := false
	entity := entities.NewShortenBulkEntity("firebase.google.com", 0)

	exp := &repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		Entity:    entity,
		Locked:    locked,
		CreatedAt: *timeObj,
		UpdatedAt: *timeObj,
	}

	return exp, nil
}

func getRandomDTO() (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	rnd := rand.New(
		rand.NewSource(
			time.Now().UTC().UnixNano(),
		),
	)

	lock := false
	clicks := rnd.Int63()
	url, err := randutil.ASCII(101)
	if err != nil {
		return nil, err
	}

	entity := entities.NewShortenBulkEntity(url, clicks)

	timestamp := "9999-12-31T23:59:59.99999999"
	date, err := repository_helpers.NewTimeFromTimestamp1e8(&timestamp)
	if err != nil || date == nil {
		return nil, err
	}

	exp := &repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		Entity:    entity,
		CreatedAt: *date,
		Locked:    lock,
		UpdatedAt: *date,
	}

	return exp, nil
}

func TestGet(t *testing.T) {
	firestore := getFirestore()

	exp, err := getExpectedGetDTO()
	if err != nil {
		t.Errorf("Creating a time: %s", err.Error())
	}

	shorten_bulk.TestGet(getHash, firestore, exp, t)
}

func TestGetOldest1(t *testing.T) {
	firestore := getFirestore()

	dto, err := getExpectedGetDTO()
	if err != nil {
		t.Errorf("Test error at creating a time: %s", err.Error())
	}

	exp := map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity]{}
	exp[getHash] = dto

	shorten_bulk.TestGetOldest(firestore, exp, t)
}

func TestGetOldest31(t *testing.T) {
	firestore := getFirestore()

	dto, err := getExpectedGetDTO()
	if err != nil {
		t.Errorf("Test error at creating a time: %s", err.Error())
	}

	exp := map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity]{}
	exp[getHash] = dto

	size := 31
	res, err := firestore.GetOldest(size)
	if err != nil {
		t.Errorf("Test error at get oldest elements: %s", err.Error())
	}

	if len(res) != size {
		t.Errorf("Test error at compare: Not the same size")
	}

	delete(res, getHash)
	for key, val := range res {
		if len(key) != helpers.HASH_SIZE {
			t.Errorf("Expected a hash with size %d, but got hash: '%s' with size %d", helpers.HASH_SIZE, key, len(key))
		}

		if val == nil {
			t.Errorf("Error got a nil DTO")
		}
	}
}

func TestIncrementClicks(t *testing.T) {
	firestore := getFirestore()
	shorten_bulk.TestIncrementClicks(incrementClicksHash, firestore, t)
}

func TestPost(t *testing.T) {
	dto, err := getRandomDTO()
	if err != nil {
		t.Errorf("Error creating a dto: %s", err.Error())
	}

	firestore := getFirestore()
	shorten_bulk.TestPost(postHash, firestore, dto, t)
}

func TestLockUnlock(t *testing.T) {
	firestore := getFirestore()
	shorten_bulk.TestLockUnlock(lockHash, firestore, t)
}
