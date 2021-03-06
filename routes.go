package main

import (
	"log"
	"net/http"
	"strings"
	_ "github.com/go-sql-driver/mysql"
//	"database/sql"
	"github.com/devfeel/dotweb"
)


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


	for {                                            //去掉Keyword首尾空格
		if strings.HasPrefix(keyword," ") || strings.HasSuffix(keyword," "){
			keyword = strings.TrimPrefix(keyword," ")
			keyword = strings.TrimSuffix(keyword," ")		 
		}else{
			break
		}

	}

	var sqlstr string = ""

	switch keyword{
	case "help","帮助":
		sqlstr = `select mediatype,mediaid,title,url,digest,thumbmedia from media where title = "公众号使用帮助" and mediatype = "news"  order by rand() limit 1; `
	case "about me","关于我","aboutme":
		sqlstr = `select mediatype,mediaid,title,url,digest,thumbmedia from media where title = "about me" and mediatype = "news" order by rand() limit 1; `
	case "list","文章","文章列表","ls":
		sqlstr = `select mediatype,mediaid,title,url,digest,thumbmedia from media where title = "原创文章列表" and mediatype = "news" order by rand() limit 1; `
	default:
		keyword = strings.ReplaceAll(keyword,` `,`%" and title like "%`)
		if strings.LastIndex(keyword,"+") == (len(keyword) - 2) && strings.LastIndex(keyword,"+") != -1 && len(keyword) >= 3 {        //关键字含有+ 且 倒数第二字节是+ 
			lastword := string(keyword[len(keyword)-1:])
			prefix := string(keyword[0:len(keyword)-2])
			switch lastword{
			case "V","v":
				sqlstr = `select a.mediatype,a.mediaid,a.title,a.url,a.digest,a.thumbmedia from media a inner join (select id from media  where title like "%` + prefix + `%"  and mediatype = "video" order by rand() limit 1) b on a.id=b.id; ` //20200330优化随机返回结果
			case "A","a":
				sqlstr = `select a.mediatype,a.mediaid,a.title,a.url,a.digest,a.thumbmedia from media a inner join (select id from media  where title like "%` + prefix + `%"  and mediatype = "news" order by rand() limit 1) b on a.id=b.id; ` //20200330优化随机返回结果
			case "I","i":
				sqlstr = `select a.mediatype,a.mediaid,a.title,a.url,a.digest,a.thumbmedia from media a inner join (select id from media  where title like "%` + prefix + `%"  and mediatype = "image" order by rand() limit 1) b on a.id=b.id; ` //20200330优化随机返回结果		
			default:
				sqlstr = `select a.mediatype,a.mediaid,a.title,a.url,a.digest,a.thumbmedia from media a inner join (select id from media  where title like "%` + keyword + `%"  order by rand() limit 1) b on a.id=b.id; ` //20200330优化随机返回结果

			}

		}else{
			sqlstr = `select a.mediatype,a.mediaid,a.title,a.url,a.digest,a.thumbmedia from media  a inner join (select id from media  where title like "%` + keyword + `%"  order by rand() limit 1) b on a.id=b.id; ` //20200330优化随机返回结果
	
		}
	}
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
		sqlstr := `select url from media where mediaid = "` + message.Mediathumb+ `"; `
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
