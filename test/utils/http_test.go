package test

import (
	"math/rand"
	"time"

	"log"
	"testing"

	"github.com/bitontop/gored/utils"
)

func Test_DigOutboundIP(t *testing.T) {
	str, err := utils.DigOutboundIP()
	log.Printf("IP:[%s]  ERROR:%s", str, err)
}

func Test_WebOutboundIP(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	str := utils.WebOutboundIP()
	log.Printf("IP:[%s]  ", str)
}

func Test_GetExternalIP(t *testing.T) {

	str := utils.GetExternalIP()
	log.Printf("IP:[%s]", str)
}
