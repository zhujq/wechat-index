package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"database/sql"
//	"io"
)

const Dbconn = "zhujq:Juju1234@tcp(wechat-mysql:3306)/wechat"
//const Dbconn = "zhujq:Juju1234@tcp(35.230.121.24:3316)/wechat"
//const Dbconn = "aW1JQvFFJD:9qN7iS4Ro6@tcp(remotemysql.com:3306)/aW1JQvFFJD"

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

var db *sql.DB

func main() {
	var err error

	var (
		version = flag.Bool("version", false, "version v1.0")
		port    = flag.Int("port", 8080, "listen port.")
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

	db, err = sql.Open("mysql",Dbconn)
	db.SetConnMaxLifetime(0)
	defer db.Close()
	err = db.Ping()
	if err != nil{
		log.Println("error:", err)	
		return 
	}

	InitRoute(app.Web.HttpServer)
	log.Println("Start Wechat Index Server on ", *port)
	app.Web.StartServer(*port)
}
