package main

import (
	"log"
	"net/http"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/devfeel/dotweb"
)

const Dbconn = "zhujq:Juju1234@tcp(35.230.121.24:3316)/wechat"

type App struct {
	Web      *dotweb.DotWeb
}



type ResBody struct {
	Status      string `json:"status"`
	Mediatype   string `json:"mediatype"`
	Mediaid     string `json:"mediaid"`
}

func NewApp() *App {
	var a = &App{}
	a.Web = dotweb.New()
	return a
}

func indexHandler(ctx dotweb.Context) error {
	keyword := ctx.QueryString("keyword")
	log.Println(keyword)

	var message = ResBody{
		Status:      "failed",
		Mediatype: "",
		Mediaid: "",
	}

	if keyword == "" {
		log.Println("ERROR: 没有提供keyword")
		return ctx.WriteJsonC(http.StatusOK, message)
	}

	db, err := sql.Open("mysql",Dbconn)
	defer db.Close()
	err = db.Ping()
	if err != nil{
		log.Println("error:", err)	
		return ctx.WriteJsonC(http.StatusOK, message)
	}

	for {                                            //去掉Keyword首尾空格
		if strings.HasPrefix(keyword," ") || strings.HasSuffix(keyword," "){
			keyword = strings.TrimPrefix(keyword," ")
			keyword = strings.TrimSuffix(keyword," ")		 
		}else{
			break
		}

	}

	keyword = strings.ReplaceAll(keyword,` `,`%" and title like "%`)
	sqlstr := `select mediatype,mediaid from media where title like "%` + keyword + `%"  order by rand() limit 1; `
	log.Println(sqlstr)

	row, err := db.Query(sqlstr)
	defer row.Close()
	if err != nil {
		log.Println("error:", err)	
		return ctx.WriteJsonC(http.StatusOK, message)
	}

	if err = row.Err(); err != nil {
		log.Println("error:", err)	
		return ctx.WriteJsonC(http.StatusOK, message)
	}
	
	count := 0
	for row.Next() {
		if err := row.Scan(&message.Mediatype,&message.Mediaid); err != nil {
			log.Println("error:", err)	
			return ctx.WriteJsonC(http.StatusOK, message)
		}	
		count += 1;
		message.Status = "success"
	}

	if count ==0 {
		message.Status = "failed"
		return ctx.WriteJsonC(http.StatusOK, message)
	}	
	return ctx.WriteJson(message)	
}

func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", indexHandler)
}