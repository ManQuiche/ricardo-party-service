package boot

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"ricardo/party-service/internal/core/entities"
	"ricardo/party-service/pkg/errors"
)

var (
	client *gorm.DB
)

func LoadDb() {
	// TODO: finish that shiiit
	//connectionString := fmt.Sprintf("")

	var err error
	client, err = gorm.Open(postgres.Open(
		fmt.Sprint("postgres://", dbUser, ":", dbPassword, "@", dbHost, ":", dbPort, "?sslmode=disable")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "ricardo.",
		},
	})
	if err != nil {
		errors.CannotConnectToDb(dbHost, dbPort)
	}

	_ = client.Migrator().DropTable(entities.Party{})
	err = client.AutoMigrate(&entities.Party{})
	if err != nil {
		log.Fatal("could not migrate db, exiting...")
	}

	if debug {
	}
}
