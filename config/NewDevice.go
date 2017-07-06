package config

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func CreateDevice() []string {
	args := os.Args[1:]
	log.Warningln("Type:"+"["+args[0]+"];", "Name:"+"["+args[1]+"];", "MAC:"+"["+args[2]+"]")
	if len(args) < 3 {
		panic("Incorrect devices's information")
	}
	return args
}