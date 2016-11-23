package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func configureLoginCommand(app *kingpin.Application) {
	app.Command("login", "Login").
		Action(executeLoginCommand)
}

func executeLoginCommand(context *kingpin.ParseContext) error {
	var config = GetConfig()
	var url = "http://localhost:" + config.AuthServerPort

	ch := make(chan interface{})
	go startAuthServer(ch, config)

	fmt.Println("Opening url: " + url)
	openWebsite(url)

	select {
	case <-ch:
		fmt.Println("Closing server...")
	}

	return nil
}
