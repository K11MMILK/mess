package api

import (
	"encoding/json"
	"log"
	"net/http"

	"mess/internal/models"
	"mess/internal/storage"

	"github.com/Shopify/sarama"
)

type Handler struct {
	Storage  *storage.Storage
	Producer sarama.SyncProducer
}

func NewHandler(storage *storage.Storage, producer sarama.SyncProducer) *Handler {
	return &Handler{
		Storage:  storage,
		Producer: producer,
	}
}

func (h *Handler) PostMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg.Status = "received"
	err = h.Storage.SaveMessage(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка сообщения в Kafka
	_, _, err = h.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "messages",
		Value: sarama.StringEncoder(msg.Content),
	})
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		http.Error(w, "Failed to send message to Kafka", http.StatusInternalServerError)
		return
	}

	log.Printf("Message sent to Kafka: %s", msg.Content)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "success"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	count, err := h.Storage.GetProcessedMessagesCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"processed_messages": count}
	json.NewEncoder(w).Encode(response)
}
