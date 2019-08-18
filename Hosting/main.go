package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"github.com/getsentry/raven-go"
	"github.com/iris-contrib/middleware/cors"
	ravenIris "github.com/iris-contrib/middleware/raven"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type SendImage struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
	Image    string `json:"image"`
}

type SendText struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
}

type SendBulkImage struct {
	List []SendImage `json:"list"`
}

type SendBulkText struct {
	List []SendText `json:"list"`
}

type waHandler struct {
	c *whatsapp.Conn
}

type Resp struct {
	Results Results `json:"results"`
}

type Results struct {
	Messages []Messages `json:"messages"`
}

type Messages struct {
	Content string `json:"content"`
}

type Config struct {
	Sentry string `json:"sentry"`
	SAP    string `json:"SAP"`
}

var (
	wac, _         = whatsapp.NewConn(5 * time.Second)
	dir, _         = filepath.Abs(filepath.Dir(os.Args[0]))
	requestChannel chan whatsapp.TextMessage
	now            = time.Now().Unix()
	config         Config
)

func (h *waHandler) HandleError(err error) {
	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Printf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

func init() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	fmt.Println("for login error please delete whatsappSession.gob at this folder -> ", os.TempDir())
	fmt.Println("running on " + string(runtime.NumCPU()) + "cores.")

	raven.SetDSN(config.Sentry)

	wac.AddHandler(&waHandler{wac})
	err = login(wac)
	if err != nil {
		panic("Error logging in: \n" + err.Error())
	}

	<-time.After(3 * time.Second)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	requestChannel = make(chan whatsapp.TextMessage)
	var (
		temp Resp
	)
	go func() {
		for {
			request, ok := <-requestChannel
			if ok {
				url := "https://api.cai.tools.sap/build/v1/dialog"
				payload := strings.NewReader(`{"message": {"content":"` + request.Text + `","type":"text"}, "conversation_id": "` + request.Info.RemoteJid[2:12] + `"}`)

				req, _ := http.NewRequest("POST", url, payload)

				req.Header.Add("Authorization", "Token "+config.SAP)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Cache-Control", "no-cache")
				req.Header.Add("Host", "api.cai.tools.sap")
				//req.Header.Add("Accept-Encoding", "gzip, deflate")
				//req.Header.Add("Connection", "keep-alive")
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					panic(err)
				}
				defer res.Body.Close()
				//body, _ := ioutil.ReadAll(res.Body)
				//fmt.Println(res)
				//fmt.Println(string(body))
				if res.StatusCode != 200 {
					panic(err)
				}
				err = json.NewDecoder(res.Body).Decode(&temp) // Decoding the  ES JSON response
				if err != nil {
					panic(err)
				}
				if len(temp.Results.Messages) > 0 && temp.Results.Messages[0].Content != "I trigger the fallback skill because I don't understand or I don't know what I'm supposed to do..." {
					fmt.Println(temp.Results.Messages[0].Content)
					to := request.Info.RemoteJid[2:12]
					mess := temp.Results.Messages[0].Content
					//	fmt.Println(request.Info.RemoteJid[2:12])
					fmt.Println(texting(to, mess))
				}
			}
		}
	}()

	app := iris.New()
	app.Logger().SetLevel("debug")
	requestLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})
	app.Use(requestLogger)
	app.Use(ravenIris.RecoveryHandler)
	// using Cors
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})
	v1 := app.Party("/", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Get("/ping", ping) // Test url
		v1.Get("/sendText", sendText)
		v1.Post("/sendBulk", sendBulk)
		v1.Get("/sendImage", sendImage)
		v1.Post("/sendBulkImg", sendBulkImg)
	}
	err := app.Run(iris.Addr(":8081"), iris.WithOptimizations, iris.WithoutBanner, iris.WithoutStartupLog)
	if err != iris.ErrServerClosed {
		panic("Shutdown with error: " + err.Error())
	}
}

func (*waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	if len(message.Text) < 20 && len(message.Text) > 1 {
		if message.Info.Timestamp > uint64(now) {
			requestChannel <- message
		}
	}
	//fmt.Println(len(message.Text))
}

func ping(ctx iris.Context) {
	ctx.WriteString("pong")
}

func sendText(ctx iris.Context) {
	to := strings.Replace(ctx.URLParamDefault("to", "1234567890"), " ", "", -1)
	mess := ctx.URLParamDefault("msg", "testing")
	ctx.WriteString(texting(to, mess))
}

func sendImage(ctx iris.Context) {
	to := strings.Replace(ctx.URLParamDefault("to", "1234567890"), " ", "", -1)
	mess := ctx.URLParamDefault("msg", "testing")
	img := ctx.URLParamDefault("img", "testImg.jpg")
	v := SendImage{
		Receiver: to,
		Message:  mess,
		Image:    img,
	}
	ctx.WriteString(image(v))
}

func sendBulk(ctx iris.Context) {
	var data SendBulkText
	m := make(map[string]string)

	ctx.ReadJSON(&data)

	for _, each := range data.List {
		each.Receiver = strings.Replace(each.Receiver, " ", "", -1)
		if each.Receiver != "" {
			m[each.Receiver] = texting(each.Receiver, each.Message)
		}
	}
	ctx.JSON(m)
}

func sendBulkImg(ctx iris.Context) {
	var data SendBulkImage
	m := make(map[string]string)

	ctx.ReadJSON(&data)

	for _, each := range data.List {
		if each.Receiver != "" {
			each.Receiver = strings.Replace(each.Receiver, " ", "", -1)
			v := SendImage{
				Receiver: each.Receiver,
				Message:  each.Message,
				Image:    each.Image,
			}

			m[each.Receiver] = image(v)
		}
	}
	ctx.JSON(m)
}

func texting(to, mess string) string {
	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "91" + to + "@s.whatsapp.net",
		},
		Text: mess,
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		panic("Error sending message: to " + to + " " + err.Error())
		return "Error"
	}
	return "Message Sent -> " + to + " : " + msgId
}

func image(v SendImage) string {
	var folder string
	if runtime.GOOS == "windows" {
		folder = `\files\`
	} else {
		folder = "/files/"
	}
	img, err := os.Open(dir + folder + v.Image)
	if err != nil {
		panic("Error reading file: " + err.Error())
		return "Error"
	}

	msg := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "91" + v.Receiver + "@s.whatsapp.net",
		},
		Type:    "image/jpeg",
		Caption: v.Message,
		Content: img,
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		panic("Error sending message: to " + v.Receiver + " " + err.Error())
		return "Error"
	}
	return "Message Sent -> " + v.Receiver + " : " + msgId
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			//terminal := qrcodeTerminal.New()
			//terminal.Get(<-qr).Print()
			obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlue, qrcodeTerminal.ConsoleColors.BrightGreen, qrcodeTerminal.QRCodeRecoveryLevels.Low)
			obj.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
