package errors

import (
	"log"
)

const (
	missingEnvVarF = "Missing %s env variables. Aborting..."
)

func MissingEnvVarF(name string) {
	log.Fatalf(missingEnvVarF, name)
}
