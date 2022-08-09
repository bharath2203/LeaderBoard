package services

import (
	"TopKScores/internal/services/providers"
	redis_sorted_sets "TopKScores/internal/services/strategies/redis-sorted-sets"
	"github.com/spf13/viper"
)

const redisSortedSetKey = "game-service"

func NewTopKScoreProvider() providers.TopKScoreProvider {
	maxNumberOfRecordsToSupport := viper.GetInt64("MAX_SIZE_OF_SORTED_SET")
	return redis_sorted_sets.NewRedisSortedSetService(redisSortedSetKey, maxNumberOfRecordsToSupport)
}
