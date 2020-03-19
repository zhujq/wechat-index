package main

import (
	"crypto/sha1"
	"encoding/xml"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"os"
	"bytes"
	"flag"
	"strconv"
)

const (
	token = "wechat4go"
)

const GetTokenUrl = "http://35.230.121.24:10527/token?appid=wxf183d5e1fe4d5204"
const GetMaterialSum = "https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token="
const GetMaterial = "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token="
const GetMediainfo = "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token="
const GetIndexUrl = "http://35.230.121.24:5901/?keyword="
const WelcomeMsg =  "谢谢您的关注！[微笑]\n“断简遗编”个人公众号主要用来记录本人用简单、理性的态度体验这大千世界的所见、所听、所想、所思，希望花费您宝贵时间的关注能够让您有所得。内容逐步完善中，您可以输入 help 或 帮助 获得本公众号使用帮助，也可以试试输入 视频 云计算 mtv ctu 看看有没好玩的。\n       Best Wishes!\n                                                                                                                Zhujq [猪头]"

type TextRequestBody struct {                    //请求结构，需要解析xml后才能赋值给它
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
	Event		 string
}

type TextResponseBody struct {                   //文本响应结构，需要用xml编码后才能http发送
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   time.Duration
	MsgType      CDATAText
	Content      CDATAText
}

type ImageResponseBody struct {                   //图片响应结构，需要用xml编码后才能http发送
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   time.Duration
	MsgType      CDATAText
	ImageMediaid CDATAText   `xml:"Image>MediaId"`
}

type VoiceResponseBody struct {                   //音频响应结构，需要用xml编码后才能http发送
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   time.Duration
	MsgType      CDATAText
	VoiceMediaid CDATAText   `xml:"Voice>MediaId"`
}


type VideoResponseBody struct {                   //视频响应结构，需要用xml编码后才能http发送
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   time.Duration
	MsgType      CDATAText
	VideoMediaid CDATAText   `xml:"Video>MediaId"`
	VideoTitle CDATAText   `xml:"Video>Title"`
	VideoDesc CDATAText   `xml:"Video>Description"`
}

type NewsResponseBody struct {                   //图文响应结构，需要用xml编码后才能http发送
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   time.Duration
	MsgType      CDATAText
	ArticleCount  int
	NewsTitle  CDATAText `xml:"Articles>item>Title"`
	NewsDesc  CDATAText `xml:"Articles>item>Description"`
	NewsPicurl  CDATAText `xml:"Articles>item>PicUrl"`
	NewsUrl  CDATAText `xml:"Articles>item>Url"`
}


type Token struct {
	AccessToken string `json:"access_token"`
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


type CDATAText struct {
	Text string `xml:",innerxml"`
}

type MediaVideoinfo struct {
	Mediaid string
	Title string 
	Desc  string 
	Url   string 
}

type MediaNewsinfo struct {
	Mediaid string
	Title string 
	Desc  string 
	Picurl string
	Url   string 
}

func init() {                                         //初始，日志文件生成
	file := "./" +"log"+ ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
			panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[wechat]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}

func makeSignature(timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func validateUrl(w http.ResponseWriter, r *http.Request) bool {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signatureGen := makeSignature(timestamp, nonce)

	signatureIn := strings.Join(r.Form["signature"], "")
	if signatureGen != signatureIn {
		return false
	}
	log.Println("signature check pass!")                        //日志记录签名通过
	echostr := strings.Join(r.Form["echostr"], "")
	fmt.Fprintf(w, echostr)                                    //echostr作为body返回给微信公众服务器，只在接入鉴权时带echostr
	return true
}

func parseTextRequestBody(r *http.Request) *TextRequestBody {   //读取http请求中的body部分赋值给TextRequestBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
//	fmt.Println(string(body))
	log.Println(string(body))                                  //收到的body写入日志文件
	requestBody := &TextRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}

func value2CDATA(v string) CDATAText {
	//return CDATAText{[]byte("<![CDATA[" + v + "]]>")}
	return CDATAText{"<![CDATA[" + v + "]]>"}
}

func makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {    //赋值TextResponseBody后用xml编码
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = value2CDATA(fromUserName)
	textResponseBody.ToUserName = value2CDATA(toUserName)
	textResponseBody.MsgType = value2CDATA("text")
	textResponseBody.Content = value2CDATA(content)
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func makeImageResponseBody(fromUserName, toUserName, imageid string) ([]byte, error) {    //赋值ImageResponseBody后用xml编码
	imageResponseBody := &ImageResponseBody{}
	imageResponseBody.FromUserName = value2CDATA(fromUserName)
	imageResponseBody.ToUserName = value2CDATA(toUserName)
	imageResponseBody.MsgType = value2CDATA("image")
	imageResponseBody.ImageMediaid = value2CDATA(imageid)
	imageResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(imageResponseBody, " ", "  ")
}

func makeVoiceResponseBody(fromUserName, toUserName, voiceid string) ([]byte, error) {    //赋值VoiceResponseBody后用xml编码
	voiceResponseBody := &VoiceResponseBody{}
	voiceResponseBody.FromUserName = value2CDATA(fromUserName)
	voiceResponseBody.ToUserName = value2CDATA(toUserName)
	voiceResponseBody.MsgType = value2CDATA("voice")
	voiceResponseBody.VoiceMediaid = value2CDATA(voiceid)
	voiceResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(voiceResponseBody, " ", "  ")
}

func makeVideoResponseBody(fromUserName string, toUserName string, videoinfo MediaVideoinfo) ([]byte, error) {    //赋值VideoResponseBody后用xml编码
	videoResponseBody := &VideoResponseBody{}
	videoResponseBody.FromUserName = value2CDATA(fromUserName)
	videoResponseBody.ToUserName = value2CDATA(toUserName)
	videoResponseBody.MsgType = value2CDATA("video")
	videoResponseBody.VideoMediaid = value2CDATA(videoinfo.Mediaid)
	videoResponseBody.VideoTitle = value2CDATA(videoinfo.Title)
	videoResponseBody.VideoDesc = value2CDATA(videoinfo.Desc)
	videoResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(videoResponseBody, " ", "  ")
}

func makeNewsResponseBody(fromUserName string, toUserName string, newsinfo MediaNewsinfo) ([]byte, error) {    //赋值NewsResponseBody后用xml编码
	newsResponseBody := &NewsResponseBody{}
	newsResponseBody.FromUserName = value2CDATA(fromUserName)
	newsResponseBody.ToUserName = value2CDATA(toUserName)
	newsResponseBody.MsgType = value2CDATA("news")
	newsResponseBody.ArticleCount = 1
	newsResponseBody.NewsTitle = value2CDATA(newsinfo.Title)
	newsResponseBody.NewsDesc = value2CDATA(newsinfo.Desc)
	newsResponseBody.NewsPicurl = value2CDATA(newsinfo.Picurl)
	newsResponseBody.NewsUrl = value2CDATA(newsinfo.Url)
	newsResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(newsResponseBody, " ", "  ")

}

//HTTPGet get 请求
func HTTPGet(uri string) ([]byte, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

func httpClient() *http.Client {
	return &http.Client{ }
}

//HTTPPost post 请求
func HTTPPost(uri string, data string) ([]byte, error) {
	body := bytes.NewBuffer([]byte(data))
	response, err := http.Post(uri, "", body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//PostJSON post json 数据请求
func PostJson(uri string, obj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient().Post(uri, "application/json;charset=utf-8", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http post error : uri=%v , statusCode=%v", uri, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func procRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !validateUrl(w, r) {
		log.Println("Wechat Service: this http request is not from Wechat platform!")
		return
	}

	if r.Method == "GET" {  //get方法只有接入鉴权时用，所以第一步vlidateUrl后无需再处理
		return
	}

	if r.Method !=  "POST" {  //按照规范，应该收到POST信息，如果不是，直接返回SUCCESS
		fmt.Fprintf(w, string("success"))
		return
	}
	
	textRequestBody := parseTextRequestBody(r)
	var rsp ResBody
	responseBody := make([]byte, 0)

	if textRequestBody.Event == "subscribe" {  //订阅事件处理

		responseBody, _ = makeTextResponseBody(textRequestBody.ToUserName,textRequestBody.FromUserName,WelcomeMsg)
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprintf(w, string(responseBody))
		return

	}

	if textRequestBody.MsgType != "text" {   //收到非文本消息

		responseBody, _ = makeTextResponseBody(textRequestBody.ToUserName,textRequestBody.FromUserName,"您好，目前只能识别文本消息，请重新输入文本。")
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprintf(w, string(responseBody))
		return

	}

	if textRequestBody == nil || textRequestBody.Content ==""{  //空内容时直接返回
		fmt.Fprintf(w, string("success"))
		return
	}

	fmt.Printf("Wechat Service: Recv text msg [%s] from user [%s]!",textRequestBody.Content,textRequestBody.FromUserName)

	if textRequestBody.Content == " " {
		fmt.Fprintf(w, string("success"))
		return
	}	

	if strings.Contains(textRequestBody.Content,`/`) {  //接收到表情符号时回填模式
		responseBody, _ = makeTextResponseBody(textRequestBody.ToUserName,textRequestBody.FromUserName,"我也"+textRequestBody.Content)
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprintf(w, string(responseBody))
		return
	}
	if strings.Contains(textRequestBody.Content,"\n") {  //换行符替换成空格

		textRequestBody.Content = strings.Replace(textRequestBody.Content, "\n", " ", -1)

	}
	

	buff, err := HTTPGet(GetIndexUrl+textRequestBody.Content)
	if err != nil{

		log.Println("error:", err)

	}else{
		err := json.Unmarshal(buff,&rsp)
		
		if err != nil {
			
			log.Println("error:", err)

			}

	}

	if rsp.Status == "success" {
		switch rsp.Mediatype{
			case "image":
				responseBody, err = makeImageResponseBody(textRequestBody.ToUserName,textRequestBody.FromUserName,rsp.Mediaid)
				log.Println(string(responseBody))
				if err != nil {
					log.Println("Wechat Service: makeimageResponseBody error: ", err)
								
				}
			case "voice":
				responseBody, err = makeVoiceResponseBody(textRequestBody.ToUserName,textRequestBody.FromUserName,rsp.Mediaid)
				log.Println(string(responseBody))
				if err != nil {
						log.Println("Wechat Service: makevoiceResponseBody error: ", err)
								
				}
						
			case "video":
				var video MediaVideoinfo
				video.Mediaid = rsp.Mediaid
				video.Title = rsp.Mediatitle
				video.Desc = rsp.Mediadigest
				video.Url = rsp.Mediaurl 

				responseBody, err = makeVideoResponseBody(textRequestBody.ToUserName,textRequestBody.FromUserName,video)
				log.Println(string(responseBody))
				if err != nil {
					log.Println("Wechat Service: makevideoResponseBody error: ", err)
								
				}

			case "news":
				var news MediaNewsinfo
				news.Mediaid = rsp.Mediaid
				news.Title = rsp.Mediatitle
				news.Desc = rsp.Mediadigest
				news.Url = rsp.Mediaurl  
				news.Picurl = rsp.Mediathumb

				responseBody, err = makeNewsResponseBody(textRequestBody.ToUserName,textRequestBody.FromUserName,news)
				log.Println(string(responseBody))
				if err != nil {
					log.Println("Wechat Service: makenewsResponseBody error: ", err)
								
				}

				default:
						
		}
			
	}else{    									// 关键字查询失败(包括不能命中或者其他失败）时直接回显Hello,后续与IBM联调
		
		responseBody, err = makeTextResponseBody(textRequestBody.ToUserName,textRequestBody.FromUserName,WelcomeMsg)
						
		log.Println(string(responseBody))
		
		if err != nil {
			log.Println("Wechat Service: makeTextResponseBody error: ", err)
			fmt.Fprintf(w, string("success"))
			return
		}

	}
	w.Header().Set("Content-Type", "text/xml")
	log.Println(string(responseBody))
	fmt.Fprintf(w, string(responseBody))
			
}

func main() {                                         //主函数入口

	var (
		version = flag.Bool("version", false, "version v1.0")
		port    = flag.Int("port", 8080, "listen port.")
	)

	flag.Parse()

	if *version {
		fmt.Println("v1.0")
		os.Exit(0)
	}

	log.Println("Wechat Service Starting")
	http.HandleFunc("/", procRequest)
	err := http.ListenAndServe((":"+strconv.Itoa(*port)), nil)
	if err != nil {
		log.Fatal("Wechat Service: ListenAndServe failed, ", err)
	}
	log.Println("Wechat Service: Stop!")
}