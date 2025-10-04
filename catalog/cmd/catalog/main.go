package main

import (
	"log"
	"time"

	"github.com/Varun5711/fritzy/catalog"
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

	var r catalog.Repository
	retry.ForeverSleep(1*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
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
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, kafkaProducer, 8080))
}
