package main

import (
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/edarha/kafka-golang/internals/configs"
	"github.com/edarha/kafka-golang/internals/must"
	"github.com/edarha/kafka-golang/internals/repositories"
	"github.com/edarha/kafka-golang/internals/services"
	"github.com/gin-gonic/gin"
)

var (
	brokers = "localhost:29092,localhost:29093,localhost:29094"
)

func main() {
	cfg := configs.PostgreSQL{
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
		Username: "postgres",
		Password: "admin",
	}

	db := must.ConnectPostgresql(&cfg)

	// migrate database
	if err := must.MigrateDB(db); err != nil {
		log.Fatal("something wrong while migrating database. err: %w", err)
	}

	// init kafka producer
	producer, err := SetupProducer()
	if err != nil {
		log.Fatal("something wrong while setup producer. err: %w", err)

	}

	// init repo
	studentRepo := repositories.NewStudentRepo(db)
	classStudentRepo := repositories.NewClassStudentRepo(db)

	// init service
	studentSvc := services.NewStudent(studentRepo, producer)
	classStudentSvc := services.NewClassStudent(classStudentRepo, producer)

	// setup server
	r := gin.Default()

	// student
	r.POST("/student", studentSvc.Post)
	r.PUT("/student/:id", studentSvc.Put)

	// class_student
	r.POST("/class-student/register", classStudentSvc.Register)

	r.Run(":8080")
}

// setupProducer will create a AsyncProducer and returns it
func SetupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	return sarama.NewSyncProducer(strings.Split(brokers, ","), config)
}
