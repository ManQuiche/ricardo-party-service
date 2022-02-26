package errors

import "log"

const (
	cannotConnectToDb = "Cannot connect to database %s:%s ! Aborting ..."
)

func CannotConnectToDb(host, port string) {
	log.Fatalf(cannotConnectToDb, host, port)
}
