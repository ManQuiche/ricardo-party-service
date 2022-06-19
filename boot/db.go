package boot

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
			TablePrefix: fmt.Sprint(dbSchema, "."),
		},
	})
	if err != nil {
		errors.CannotConnectToDb(dbHost, dbPort)
	}

	//_ = client.Migrator().DropTable(entities.Party{})
	//_ = client.Migrator().DropTable(entities.User{})
	//_ = client.Migrator().DropTable("ricardo.party_members")
	//err = client.AutoMigrate(&entities.Party{})
	//if err != nil {
	//	log.Fatal("could not migrate db, exiting...")
	//}

	if debug {
	}
}
