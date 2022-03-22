package main

import (
	"flag"
	"fmt"
	"github.com/uanid/fakenews-server/application"
	"log"
)

func main() {
	port := flag.Int("port", 8080, "Rest Api Server Port")
	flag.Parse()

	app := application.NewApplication(*port)
	err := app.Start()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Stop Application.")
}
