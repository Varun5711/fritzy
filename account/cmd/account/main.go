package main

import (
	"log"
	"time"

	"github.com/Varun5711/fritzy/account"
	"github.com/Varun5711/fritzy/kafka"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL         string `envconfig:"DATABASE_URL"`
	KafkaBrokers        string `envconfig:"KAFKA_BROKERS"`
	KafkaConsumerGroup  string `envconfig:"KAFKA_CONSUMER_GROUP"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
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

	log.Println("Listening on port 8080 ...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, kafkaProducer, 8080))
}
