package redis_sorted_sets

import (
	"TopKScores/api"
	"TopKScores/internal/services/providers"
	"TopKScores/internal/services/subjects/gamescore"
	"fmt"
)

const observerId = "ID-REDIS-SORTED-SETS-IMPL"

type gameServiceObserverImpl struct {
	topKScoreService providers.TopKScoreProvider
}

func NewRedisSortedSetObserver(topKScoreService providers.TopKScoreProvider) gamescore.Observer {
	return &gameServiceObserverImpl{
		topKScoreService: topKScoreService,
	}
}

func (g gameServiceObserverImpl) Update(gameScore *api.GameScore) {
	go func() {
		err := g.topKScoreService.AddScore(gameScore)
		if err != nil {
			// Need to add retry logic here. Printing errors for simplicity purpose.
			fmt.Println(err)
		}
	}()
}

func (g gameServiceObserverImpl) GetId() string {
	return observerId
}
