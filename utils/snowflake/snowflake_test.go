package snowflake

import (
	"log"
	"testing"
)

func TestGenID(t *testing.T) {
	if err := Init("2022-12-19", 1); err != nil {
		log.Fatalf("init err: %s", err)
	}
	log.Println(GenID())
}
