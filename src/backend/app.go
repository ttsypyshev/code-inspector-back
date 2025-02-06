package backend

import (
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "rip/docs/swagger"
	"rip/lib/config"
)

type App struct {
	db          *Db
	minioClient *minio.Client
	redisClient *redis.Client
	secret      string
}

func Run() error {
	log.Println("Server starting up")

	r := gin.Default()

	app, err := NewDB(config.FromEnvDB())
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
		return err
	}

	app.minioClient, err = InitializeMinIO()
	if err != nil {
		log.Fatalf("Error initializing MinIO: %v", err)
		return err
	}

	app.redisClient, err = InitializeRedis()
	if err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
		return err
	}

	app.secret, err = config.LoadSecret()
	if err != nil {
		log.Fatalf("Error load secret phrase: %v", err)
		return err
	}

	rand.Seed(time.Now().UnixNano())

	app.SetupRoutes(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := r.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		return err
	}

	log.Println("Server stopped")
	return nil
}
