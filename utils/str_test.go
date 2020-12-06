package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestStr(t *testing.T) {
	log.Println(fmt.Sprintf("%02d", 1))
	log.Println(fmt.Sprintf("%02d", 2))
	log.Println(fmt.Sprintf("%02d", 12))
	log.Println(fmt.Sprintf("%02d", 123))
}
