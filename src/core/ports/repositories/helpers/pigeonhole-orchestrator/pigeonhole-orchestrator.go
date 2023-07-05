package repository_helpers

import (
	"errors"
	"log"
	"strconv"
	"sync"

	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type PigeonholeOrchestrator[T any, K any] struct {
	worksSize int
	repos     *[]*T
}

func (o *PigeonholeOrchestrator[T, K]) ExecuteSingleFunc(
	singleFunc func(*T) (*repositories.RepositoryDTO[K], error),
) (*repositories.RepositoryDTO[K], error) {
	size := len(*o.repos)
	idxsCh := helpers.NewRandChanIdxs(uint(size))

	if len(idxsCh) < o.worksSize {
		log.Fatalf("ExecuteSingleFunc: len(idxsCh)[%d] < o.worksSize[%d]", len(idxsCh), o.worksSize)
		return nil, errors.New("Internal error: Not enough repositories")
	}

	var wg sync.WaitGroup
	resCh := make(chan *repositories.RepositoryDTO[K], o.worksSize)

	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for idx, ok := <-idxsCh; ok; idx, ok = <-idxsCh {
				repo := (*o.repos)[idx]
				res, err := singleFunc(repo)
				if err == nil {
					resCh <- res
					break
				}
				log.Println("Go func error:", err.Error())
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(resCh)

	if len(resCh) != o.worksSize {
		log.Printf("len(resCh)[%d] != o.worksSize[%d]", len(resCh), o.worksSize)
		return nil, errors.New("Internal error: Not enough successful workers")
	}

	valMp := map[string]*repositories.RepositoryDTO[K]{}
	countMp := map[string]int{}
	for result := range resCh {
		timestamp := ""
		if result != nil {
			timestamp = helpers.TimeToTimestamp1e8(result.UpdatedAt)
		}

		count, ok := countMp[timestamp]
		if !ok {
			valMp[timestamp] = result
			count = 0
		}
		countMp[timestamp] = count + 1
	}

	bestTimestamp := ""
	bestCount := 0
	if count, ok := countMp[bestTimestamp]; ok {
		bestCount = count
	}
	for timestamp, count := range countMp {
		if (bestCount < count) || (bestCount == count && bestTimestamp < timestamp) {
			bestCount = count
			bestTimestamp = timestamp
		}
	}

	return valMp[bestTimestamp], nil
}

func (o *PigeonholeOrchestrator[T, K]) ExecuteMultipleFunc(
	multipleFunc func(*T) (map[string]*repositories.RepositoryDTO[K], error),
) (map[string]*repositories.RepositoryDTO[K], error) {
	size := len(*o.repos)
	idxsCh := helpers.NewRandChanIdxs(uint(size))
	res := map[string]*repositories.RepositoryDTO[K]{}

	if len(*o.repos) < o.worksSize {
		return res, errors.New("Internal error: Not enough repositories")
	}

	var wg sync.WaitGroup
	resCh := make(chan map[string]*repositories.RepositoryDTO[K], o.worksSize)
	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for idx, ok := <-idxsCh; ok; idx, ok = <-idxsCh {
				repo := (*o.repos)[idx]
				res, err := multipleFunc(repo)
				if err == nil {
					resCh <- res
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(resCh)

	if len(resCh) != o.worksSize {
		log.Printf("len(resCh)[%d] != o.worksSize[%d]", len(resCh), o.worksSize)
		return res, errors.New("Internal error: Not enough successful workers")
	}

	valMp := map[string]map[string]*repositories.RepositoryDTO[K]{}
	countMp := map[string]map[string]int{}
	for mp := range resCh {
		for hash, dto := range mp {
			_, ok := valMp[hash]
			if !ok {
				valMp[hash] = map[string]*repositories.RepositoryDTO[K]{}
				countMp[hash] = map[string]int{}
			}

			timestamp := ""
			if dto != nil {
				timestamp = helpers.TimeToTimestamp1e8(dto.UpdatedAt)
			}

			count, _ := countMp[hash][timestamp]
			if !ok {
				valMp[hash][timestamp] = dto
				count = 0
			}
			countMp[hash][timestamp] = count + 1
		}
	}

	ansMp := map[string]*repositories.RepositoryDTO[K]{}
	for hash, mp := range countMp {
		countSum := 0
		bestTimestamp := ""

		bestCount := 0
		if count, ok := mp[bestTimestamp]; ok {
			bestCount = count
		}

		for timestamp, count := range mp {
			if (bestCount < count) || (bestCount == count && bestTimestamp < timestamp) {
				bestCount = count
				bestTimestamp = timestamp
			}
			countSum += count
		}
		if countSum <= o.worksSize {
			ansMp[hash] = valMp[hash][bestTimestamp]
		}
	}

	return ansMp, nil
}

func NewPigeonholeOrchestrator[T any, K any](
	repos *[]*T,
) (*PigeonholeOrchestrator[T, K], error) {
	if repos == nil {
		return nil, errors.New("Internal error: repos is a nil pointer")
	}

	size := len(*repos)
	if size == 0 {
		return nil, errors.New("Internal error: repositories array is empty")
	}

	for idx, repo := range *repos {
		if repo == nil {
			return nil, errors.New("Internal error: Repos[" + strconv.Itoa(idx) + "] is a nil pointer")
		}
	}

	return &PigeonholeOrchestrator[T, K]{
		worksSize: size/2 + 1,
		repos:     repos,
	}, nil
}
