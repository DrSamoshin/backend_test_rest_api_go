package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	todo "github.com/ssamoshin-ms/golang_test"
	"github.com/ssamoshin-ms/golang_test/pkg/handler"
	"github.com/ssamoshin-ms/golang_test/pkg/repository"
	"github.com/ssamoshin-ms/golang_test/pkg/service"
	"github.com/subosito/gotenv"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatalf("error init configs: %s", err.Error())
	}

	if err := gotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmade"),
	})
	if err != nil {
		log.Fatalf("failed to init db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
