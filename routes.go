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
	Mediatitle  string `json:"mediatitle"`
	Mediaurl    string `json:"mediaurl"`
	Mediadigest string `json:"mediadigest"`
	Mediathumb  string `json:"mediathumb"`
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
		Mediaurl: "",
		Mediadigest: "",
		Mediathumb: "",
	}

	if keyword == "" {
		log.Println("ERROR: 没有提供keyword")
		return ctx.WriteJsonC(http.StatusNotFound, message)
	}

	db, err := sql.Open("mysql",Dbconn)
	defer db.Close()
	err = db.Ping()
	if err != nil{
		log.Println("error:", err)	
		return ctx.WriteJsonC(http.StatusNotFound, message)
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
	sqlstr := `select mediatype,mediaid,title,url,digest,thumbmedia from media where title like "%` + keyword + `%"  order by rand() limit 1; `
	log.Println(sqlstr)

	row, err := db.Query(sqlstr)
	defer row.Close()
	if err != nil {
		log.Println("error:", err)	
		return ctx.WriteJsonC(http.StatusNotFound, message)
	}

	if err = row.Err(); err != nil {
		log.Println("error:", err)	
		return ctx.WriteJsonC(http.StatusNotFound, message)
	}
	
	count := 0
	for row.Next() {
		if err := row.Scan(&message.Mediatype,&message.Mediaid,&message.Mediatitle,&message.Mediaurl,&message.Mediadigest,&message.Mediathumb); err != nil {
			log.Println("error:", err)	
			return ctx.WriteJsonC(http.StatusNotFound, message)
		}	
		count += 1;
		message.Status = "success"
	}

	if count ==0 {
		message.Status = "failed"
		return ctx.WriteJsonC(http.StatusNotFound, message)
	}	

	if message.Mediatype == "news"{    //图文类型时把封面图片的mediaid转换为Picurl
		sqlstr := `select url from media where mediaid = "` + message.Mediaid + `"; `
		rows, _ := db.Query(sqlstr)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&message.Mediathumb)
		}
	}
	
	return ctx.WriteJson(message)	
}

func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", indexHandler)
}
