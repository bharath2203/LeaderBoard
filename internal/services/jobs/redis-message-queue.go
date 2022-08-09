package jobs

import (
	"TopKScores/api"
	"TopKScores/internal/services/subjects/gamescore"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

// Redis key that holds the message queue
const jobKey = "game_service_queue"

// RedisMessageQueueService exposes two functionalities.
// 1. Add a GameScore to the queue.
// 2. Start Listening to the queue and notify all the subscribers.
type RedisMessageQueueService struct {
	redisClient *redis.Client
	subject     gamescore.Subject
}

func NewRedisMessageQueueService(redisClient *redis.Client, subject gamescore.Subject) *RedisMessageQueueService {
	return &RedisMessageQueueService{
		redisClient: redisClient,
		subject:     subject,
	}
}

func (s *RedisMessageQueueService) AddToQueue(score *api.GameScore) error {
	err := s.redisClient.LPush(jobKey, score).Err()
	return err
}

func (s *RedisMessageQueueService) StartListening() {
	go func() {
		for {
			message, err := s.redisClient.BRPop(0*time.Second, jobKey).Result()
			if err != nil {
				// Todo: Add log, metrics and also configure some alert to notify the POC.
				fmt.Printf("err in listing to message queue: %v", err)
				continue
			}

			gameScore := &api.GameScore{}
			err = gameScore.UnMarshalBinary([]byte(message[1]))
			if err != nil {
				fmt.Printf("unmarshall error: %v", err)
				continue
			}

			s.subject.NotifyAll(gameScore)
		}
	}()
}
