package boot

import (
	"log"
	"os"
	"ricardo/party-service/pkg/errors"
	"strconv"
)

var (
	dbHost       string
	dbPort       string
	dbUser       string
	dbPassword   string
	dbSchema     string
	port         string
	url          string
	accessSecret string

	natsURL         string
	natsUsr         string
	natsPwd         string
	natsUserDeleted string

	debug bool
)

func LoadEnv() {
	dbHost = env("DB_HOST")
	dbPort = env("DB_PORT")
	dbUser = env("DB_USER")
	dbSchema = env("DB_SCHEMA")
	dbPassword = env("DB_PASSWORD")

	port = env("PORT")
	url = env("URL")
	accessSecret = env("ACCESS_SECRET")
	debug = envBool("DEBUG")

	natsURL = env("NATS_URL")
	natsUsr = env("NATS_USR")
	natsPwd = env("NATS_PWD")
	natsUserDeleted = env("NATS_USER_DELETED")
}

func envBool(name string) bool {
	res, err := strconv.ParseBool(env(name))
	if err != nil {
		log.Fatalf("env var %s needs to be of boolean type", name)
	}

	return res
}

func env(name string) string {
	str, ok := os.LookupEnv(name)
	if !ok {
		errors.MissingEnvVarF(name)
	}

	return str
}
