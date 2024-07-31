package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mess/config"
	"mess/internal/api"
	"mess/internal/storage"

	"mess/internal/kafka"

	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка конфигурации из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found: %v", err)
	}

	// Чтение конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Подключение к PostgreSQL
	storage, err := storage.NewStorage(cfg.DBConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание Kafka Producer
	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, nil)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	handler := api.NewHandler(storage, producer)
	router := api.NewRouter(handler)

	// Запуск Kafka Consumer
	ctx, cancel := context.WithCancel(context.Background())
	consumerErrChan := make(chan error, 1)

	go func() {
		consumerErrChan <- kafka.StartConsumer(ctx, cfg.KafkaBrokers, cfg.KafkaGroupID, cfg.KafkaTopic, storage)
	}()

	// Graceful Shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		<-stopChan
		log.Println("Shutting down server...")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
		cancel()
	}()

	log.Printf("Server is running on port %s", cfg.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Ожидание завершения работы Consumer
	if consumerErr := <-consumerErrChan; consumerErr != nil {
		log.Fatalf("Consumer encountered an error: %v", consumerErr)
	}

	log.Println("Server stopped")
}
