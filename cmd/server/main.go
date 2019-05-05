package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/prophesional/intuit"

	"github.com/gravitational/configure"
)

func main() {
	var config intuit.SQLConfig

	config.Type = "mysql" // add to config
	fmt.Println(configure.ParseEnv(&config))
	fmt.Println(config)

	p, err := intuit.NewServer(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.Start()
	defer fmt.Println("closing")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

}
