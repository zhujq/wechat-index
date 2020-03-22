package main

import (
	"flag"
	"fmt"
	"log"
	"os"
//	"io"
)

var app = NewApp()

/*
func init() {                                         //初始，日志文件生成
	file := "./" +"logindex"+ ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
			panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile) // 将文件设置为log输出的文件
//	mw := io.MultiWriter(os.Stdout,logFile) //同时输出到文件和控制台
//  log.SetOutput(mw)
	log.SetPrefix("[wechat-index]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}
*/

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
