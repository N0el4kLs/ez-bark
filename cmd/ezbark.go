package main

import (
	"log"

	"github.com/N0el4kLs/ez-bark/pkg/util"
)

func main() {
	options, err := util.InitOptions()
	if err != nil {
		return
	}

	if options == nil && err == nil {
		log.Println("Init config successfully")
		return
	}

	if options.TestSwitch {
		options.Title = "Test"
		options.Body = "This is Test case"
	}

	err = options.Notice()
	if err != nil {
		log.Fatalln("Send notification err.Check your server and config key")
	}
}
