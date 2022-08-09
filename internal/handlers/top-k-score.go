package handlers

import (
	"TopKScores/api"
	"TopKScores/internal/services/jobs"
	"TopKScores/internal/services/providers"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
)

type Handler struct {
	topKScoreService         providers.TopKScoreProvider
	redisMessageQueueService *jobs.RedisMessageQueueService
}

func NewHandler(topKScoreService providers.TopKScoreProvider, redisMessageQueueService *jobs.RedisMessageQueueService) *Handler {
	return &Handler{
		topKScoreService:         topKScoreService,
		redisMessageQueueService: redisMessageQueueService,
	}
}

// AddSingleScore adds a new GameScore object to the TopK Score data stream.
func (h Handler) AddSingleScore(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var addGameScoreReq *api.GameScore
	err := decoder.Decode(&addGameScoreReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// logging, metrics to be handled here.
	err = h.topKScoreService.AddScore(addGameScoreReq)
	if err != nil {
		// Todo: Use error type to decide on bad request and internal error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := &api.BaseResponse{
		StatusMessage: "success",
		StatusCode:    http.StatusOK,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) AddSingleScoreToQueue(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var addGameScoreReq *api.GameScore
	err := decoder.Decode(&addGameScoreReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// logging, metrics to be handled here.
	err = h.redisMessageQueueService.AddToQueue(addGameScoreReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := &api.BaseResponse{
		StatusMessage: "success",
		StatusCode:    http.StatusOK,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetTopKScores returns top K GameScores with respective User.
// The number of supported records can be configured through config file.
func (h Handler) GetTopKScores(w http.ResponseWriter, r *http.Request) {
	totalScoresToBeFetched := viper.GetInt64("GET_SCORE_ENTITY_COUNT")
	// logging, metrics to be handled here.
	topKScores, err := h.topKScoreService.GetTopKScores(totalScoresToBeFetched)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := &api.BaseResponse{
		Data:          topKScores,
		StatusMessage: "success",
		StatusCode:    http.StatusOK,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
