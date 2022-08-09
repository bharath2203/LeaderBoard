package redis_sorted_sets

import (
	"TopKScores/api"
	redisService "TopKScores/internal/redis"
	"TopKScores/internal/services/providers"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

// RedisSortedSets implements TopKScoreProvider.
// It uses redis sorted sets data-structure to store and organize the GameScore data.
// https://redis.io/commands/?group=sorted_set
type redisSortedSetsImpl struct {
	redisClient        *redis.Client
	redisKey           string
	maxNumberOfRecords int64
}

func NewRedisSortedSetService(redisKey string, maxNumberOfRecords int64) providers.TopKScoreProvider {
	return &redisSortedSetsImpl{
		redisClient:        redisService.NewRedisClient(viper.GetString("REDIS_HOST")),
		redisKey:           redisKey,
		maxNumberOfRecords: maxNumberOfRecords,
	}
}

// AddScore adds a new GameScore to the stream.
// If the GameId and Username are same, it replaces the existing entry  making this method idem-potent.
// It also removes the extra entries if the size of the data structure exceeds limit.
func (c *redisSortedSetsImpl) AddScore(gameScore *api.GameScore) error {

	// Todo: Implement Retryable and Non Retryable errors for client to handle.
	// Assumption - Score cannot be negative.
	if err := c.validateGameScoreObject(gameScore); err != nil {
		// non retryable error
		return err
	}
	err := c.redisClient.ZAdd(c.redisKey, &redis.Z{
		Score: gameScore.UserScore,
		Member: api.GameInstance{
			GameId:   gameScore.GameId,
			UserName: gameScore.UserName,
		},
	}).Err()
	if err != nil {
		return err
	}
	// Trim extra records
	err = c.redisClient.ZRemRangeByRank(c.redisKey, 0, -1*c.maxNumberOfRecords).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetTopKScores function returns at most top k GameScore objects ordered by score descending.
func (c *redisSortedSetsImpl) GetTopKScores(k int64) ([]*api.GameScore, error) {

	// validate
	if err := c.validateRecordCount(k); err != nil {
		return nil, err
	}

	// Query top k records from sorted set.
	queryResponse, err := c.redisClient.ZRevRangeWithScores(c.redisKey, 0, k-1).Result()
	if err != nil {
		return nil, err
	}

	return c.generateResponse(queryResponse)
}

func (c *redisSortedSetsImpl) generateResponse(queryResponse []redis.Z) ([]*api.GameScore, error) {
	var topScores []*api.GameScore
	for _, queryObj := range queryResponse {
		gameInstance, err := c.extractGameInstance(queryObj)
		if err != nil {
			return nil, err
		}
		gameScore := &api.GameScore{
			GameId:    gameInstance.GameId,
			UserName:  gameInstance.UserName,
			UserScore: queryObj.Score,
		}
		topScores = append(topScores, gameScore)
	}
	return topScores, nil
}

func (c *redisSortedSetsImpl) extractGameInstance(queryObj redis.Z) (*api.GameInstance, error) {
	member := queryObj.Member.(string)
	gameInstance := &api.GameInstance{}
	err := gameInstance.UnMarshalBinary([]byte(member))
	return gameInstance, err
}

func (c *redisSortedSetsImpl) validateRecordCount(k int64) error {
	if k > c.maxNumberOfRecords {
		return fmt.Errorf("only %v entries supported", c.maxNumberOfRecords)
	}
	return nil
}

func (c *redisSortedSetsImpl) validateGameScoreObject(gameScore *api.GameScore) error {
	if gameScore == nil {
		return fmt.Errorf("gamescore object nil")
	}
	if gameScore.UserScore <= 0 {
		return fmt.Errorf("gamescore cannot be zero/negative value")
	}
	if len(gameScore.GameId) == 0 {
		return fmt.Errorf("game identifier cannot be nil")
	}
	if len(gameScore.UserName) == 0 {
		return fmt.Errorf("user name cannot by empty")
	}
	return nil
}
