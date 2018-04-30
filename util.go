package podstalk

import (
	"log"
	"os"
)

// GetEnvOr returns the value uf the env var named `name`, or `fallback` if said
// env var is empty
func GetEnvOr(name, fallback string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return fallback
}

// Check ensures `err` is `nil`, and outputs an error message and exits if it is
// not
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
