package gamescore

import (
	"TopKScores/api"
	"github.com/go-redis/redis/v7"
)

type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	NotifyAll(gameScore *api.GameScore)
}

type subject struct {
	redisClient  *redis.Client
	observersMap map[string]Observer
}

func NewSubject() Subject {
	return &subject{
		observersMap: make(map[string]Observer),
	}
}

func (s *subject) Register(observer Observer) {
	s.observersMap[observer.GetId()] = observer
}

func (s *subject) Deregister(observer Observer) {
	delete(s.observersMap, observer.GetId())
}

func (s *subject) NotifyAll(gameScore *api.GameScore) {
	for _, observer := range s.observersMap {
		observer.Update(gameScore)
	}
}
