package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var app = NewApp()

func main() {
	var err error

	var (
		version = flag.Bool("version", false, "version v1.0")
		port    = flag.Int("port", 80, "listen port.")
	)

	flag.Parse()

	if *version {
		fmt.Println("v1.0")
		os.Exit(0)
	}

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	InitRoute(app.Web.HttpServer)
	log.Println("Start Wechat Index Server on ", *port)
	app.Web.StartServer(*port)
}
