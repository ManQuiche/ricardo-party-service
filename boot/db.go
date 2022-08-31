package boot

import (
	"fmt"
	_ "github.com/lib/pq"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
	"gitlab.com/ricardo134/party-service/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

var (
	client *gorm.DB
)

func LoadDb() {
	// TODO: finish that shiiit
	//connectionString := fmt.Sprintf("")

	var err error
	client, err = gorm.Open(postgres.Open(
		fmt.Sprint("postgres://", dbUser, ":", dbPassword, "@", dbHost, ":", dbPort, "/", dbDatabase, "?sslmode=disable")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: fmt.Sprint(dbSchema, "."),
		},
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		errors.CannotConnectToDb(dbHost, dbPort)
	}

	err = client.AutoMigrate(&entities.Party{}, &entities.User{})
	if err != nil {
		log.Fatal("could not migrate db, exiting...")
	}
}
