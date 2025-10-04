package main

import (
	"log"
	"time"

	"github.com/Varun5711/fritzy/kafka"
	"github.com/Varun5711/fritzy/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL         string `envconfig:"DATABASE_URL"`
	AccountURL          string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL          string `envconfig:"CATALOG_SERVICE_URL"`
	KafkaBrokers        string `envconfig:"KAFKA_BROKERS"`
	KafkaConsumerGroup  string `envconfig:"KAFKA_CONSUMER_GROUP"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	var kafkaProducer *kafka.Producer
	if cfg.KafkaBrokers != "" {
		kafkaProducer = kafka.NewProducer(cfg.KafkaBrokers)
		defer kafkaProducer.Close()
		log.Println("Kafka producer initialized")
	}

	log.Println("Listening on port 8080...")
	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, kafkaProducer, cfg.AccountURL, cfg.CatalogURL, 8080))
}
