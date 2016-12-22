package main

import (
	"os"
	"os/signal"

	"github.com/lucapette/deluminator/api"
)

func main() {
	api.Start(3000)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
