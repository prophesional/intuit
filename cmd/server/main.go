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

	err := configure.ParseEnv(&config)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Debug:  Config is: ", config)

	p, err := intuit.NewServer(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.Start()
	defer fmt.Println("closing")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

}
