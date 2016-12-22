package main

import "github.com/lucapette/deluminator/api"
import "log"

func main() {
	err := api.Start(3000)

	if err != nil {
		log.Panic(err)
	}
}
