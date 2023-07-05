package shorten_bulk

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/types"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

func TestGet(
	hash string,
	repo ShortenBulkRepository,
	exp *repositories.RepositoryDTO[entities.ShortenBulkEntity],
	t *testing.T,
) {

	res, err := repo.Get(hash)
	if err != nil {
		t.Errorf("Get element at hash '%s': %s", hash, err.Error())
	}

	if !reflect.DeepEqual(res, exp) {
		t.Errorf("Result and expected are different")
	}
}

func TestGetOldest(
	repo ShortenBulkRepository,
	exp map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	t *testing.T,
) {
	res, err := repo.GetOldest(1)
	if err != nil {
		t.Errorf("Test error at get oldest elements: %s", err.Error())
	}

	if len(res) != len(exp) {
		t.Errorf("Test error at compare: Not the same size")
	}

	if !reflect.DeepEqual(res, exp) {
		fmt.Println("result:")
		for key, val := range res {
			if val == nil {
				fmt.Println(key, ": nil")
				continue
			}
			fmt.Println(key, ":", types.NewShortenBulkFlattenDTO(*val))
		}

		fmt.Println("expected:")
		for key, val := range exp {
			if val == nil {
				fmt.Println(key, ": nil")
				continue
			}
			fmt.Println(key, ":", types.NewShortenBulkFlattenDTO(*val))
		}
		t.Errorf("Test error at compare the get result")
	}
}

func createTest(
	hash string,
	repo ShortenBulkRepository,
	executor func(ShortenBulkRepository) error,
	update func(*repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	),
) error {
	prev, err := repo.Get(hash)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = executor(repo)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	exp, err := update(prev)
	if err != nil {
		return err
	}

	res, err := repo.Get(hash)
	if err != nil {
		return err
	}

	if res == nil && exp == nil {
		return nil
	}

	if res == nil || exp == nil {
		if res == nil {
			return errors.New("Result and expected are different: result is nil, expected is not")
		} else {
			return errors.New("Result and expected are different: expected is nil, result is nil")
		}
	}

	resFlatten := types.NewShortenBulkFlattenDTO(*res)
	expFlatten := types.NewShortenBulkFlattenDTO(*exp)

	if !reflect.DeepEqual(resFlatten, expFlatten) {
		log.Println("res:", resFlatten)
		log.Println("exp:", expFlatten)

		return errors.New("Result and expected are different")
	}

	return nil
}

func TestIncrementClicks(
	hash string,
	repo ShortenBulkRepository,
	t *testing.T,
) {
	updatedAt := time.Now()
	executor := func(repo ShortenBulkRepository) error {
		err := repo.IncrementClicks(hash, updatedAt)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		return nil
	}

	update := func(prev *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		if prev == nil {
			return nil, errors.New("No element available")
		}
		prev.Entity.Clicks += 1
		prev.UpdatedAt = updatedAt
		return prev, nil
	}

	err := createTest(hash, repo, executor, update)
	if err != nil {
		t.Errorf("Increment Clicks error: %s", err.Error())
	}
}

func TestPost(
	hash string,
	repo ShortenBulkRepository,
	dto *repositories.RepositoryDTO[entities.ShortenBulkEntity],
	t *testing.T,
) {
	executor := func(repo ShortenBulkRepository) error {
		err := repo.Post(hash, *dto)
		if err != nil {
			return err
		}

		return nil
	}

	update := func(_ *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		return dto, nil
	}

	err := createTest(hash, repo, executor, update)
	if err != nil {
		t.Errorf("Post error: %s", err.Error())
	}
}

func TestLockUnlock(
	hash string,
	repo ShortenBulkRepository,
	t *testing.T,
) {
	updatedAt := time.Now()
	executorLock := func(repo ShortenBulkRepository) error {
		return repo.Lock(hash, updatedAt)
	}

	updateLock := func(prev *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		if prev == nil {
			return nil, errors.New("Empty dto")
		}

		if prev.Locked {
			return nil, errors.New("Already locked")
		}

		prev.Locked = true
		prev.UpdatedAt = updatedAt
		return prev, nil
	}

	err := createTest(hash, repo, executorLock, updateLock)
	if err != nil {
		t.Errorf("Lock error: %s", err.Error())
		return
	}

	updatedAt = time.Now()
	executorUnlock := func(repo ShortenBulkRepository) error {
		return repo.Unlock(hash, updatedAt)
	}

	updateUnlock := func(prev *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		if prev == nil {
			return nil, errors.New("Empty dto")
		}

		if !prev.Locked {
			return nil, errors.New("Already unlocked")
		}

		prev.Locked = false
		prev.UpdatedAt = updatedAt
		return prev, nil
	}

	err = createTest(hash, repo, executorUnlock, updateUnlock)
	if err != nil {
		t.Errorf("Unlock error: %s", err.Error())
	}
}
