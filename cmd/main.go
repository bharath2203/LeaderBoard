package main

import (
	"TopKScores/internal"
	"TopKScores/internal/handlers"
	"TopKScores/internal/redis"
	"TopKScores/internal/services"
	"TopKScores/internal/services/jobs"
	redis_sorted_sets "TopKScores/internal/services/strategies/redis-sorted-sets"
	"TopKScores/internal/services/subjects/gamescore"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	// Load config variables from environment
	viper.SetConfigFile("cmd/conf/local.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	// Initialize the infra, services etc
	redisClient := redis.NewRedisClient(viper.GetString("REDIS_HOST"))

	// Init TopK Service provider
	topKScoreService := services.NewTopKScoreProvider()

	// Init Subject and Observer for GameScores subscription.
	redisGameScoreMessageSubject := gamescore.NewSubject()
	redisSortedSetsObserver := redis_sorted_sets.NewRedisSortedSetObserver(topKScoreService)

	// Register the server
	redisGameScoreMessageSubject.Register(redisSortedSetsObserver)

	// Initialize job queue process
	redisMessageQueueService := jobs.NewRedisMessageQueueService(redisClient, redisGameScoreMessageSubject)

	handlerService := handlers.NewHandler(topKScoreService, redisMessageQueueService)
	router := internal.NewRouter(handlerService)

	// Start listening to the message queue
	redisMessageQueueService.StartListening()

	// Initialize a http server
	srv := http.Server{
		Addr:    viper.GetString("LISTEN_HOST"),
		Handler: router,
	}
	if httpErr := srv.ListenAndServe(); httpErr != nil {
		log.Fatalf("http error: %v", httpErr)
	}
}
