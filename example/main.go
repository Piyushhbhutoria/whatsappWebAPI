package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type SendText struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
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
	requestChannel chan whatsapp.TextMessage
	now            = time.Now().Unix()
	config         Config
)

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

	fmt.Println("for session issue  remove whatsapp.gob from --> " + os.TempDir())
	fmt.Println("running on " + strconv.Itoa(runtime.NumCPU()) + " cores.")

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.Sentry,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	requestChannel = make(chan whatsapp.TextMessage, runtime.NumCPU())

	wac.AddHandler(&waHandler{wac})
	if err = login(wac); err != nil {
		panic("Error logging in: \n" + err.Error())
	}

	<-time.After(3 * time.Second)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var temp Resp

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
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("error request: %v\n", err)
				}
				defer res.Body.Close()
				if res.StatusCode != 200 {
					log.Printf("Bad request: %v\n", err)
				}
				err = json.NewDecoder(res.Body).Decode(&temp) // Decoding the  ES JSON response
				if err != nil {
					log.Printf("Error decoding body: %v\n", err)
				}
				if len(temp.Results.Messages) > 0 && temp.Results.Messages[0].Content != "I trigger the fallback skill because I don't understand or I don't know what I'm supposed to do..." {
					fmt.Println(temp.Results.Messages[0].Content)
					to := request.Info.RemoteJid[2:12]
					mess := temp.Results.Messages[0].Content
					log.Printf("%v --> %v\nBot --> %v", to, request.Text, mess)
					//	fmt.Println(request.Info.RemoteJid[2:12])
					fmt.Println(texting(to, mess))
				}
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s %s %d %s %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))
	router.Use(sentrygin.New(sentrygin.Options{}))
	router.Use(cors.Default())

	router.GET("/", helloworld)
	router.GET("/ping", ping)
	router.GET("/sendText", sendText)
	router.POST("/sendBulk", sendBulk)

	if err := router.Run(":8081"); err != nil {
		log.Printf("Shutdown with error: %v\n", err)
	}

}

func (*waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	//if message.Info.Timestamp > uint64(now) && !message.Info.FromMe && len(message.Text) < 17 && len(message.Text) > 1 {
	if message.Info.Timestamp > uint64(now) && !message.Info.FromMe && !strings.Contains(message.Text, "@g.us") && len(message.Text) < 21 {
		fmt.Printf("%v from %v\n", message.Text, message.Info)
		requestChannel <- message
	}
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func helloworld(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func sendText(c *gin.Context) {
	to := strings.Replace(c.DefaultQuery("to", "1234567890"), " ", "", -1)
	mess := c.DefaultQuery("msg", "testing")
	c.String(http.StatusOK, texting(to, mess))
}

func sendBulk(c *gin.Context) {
	var data SendBulkText
	m := make(map[string]string)

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, each := range data.List {
		each.Receiver = strings.Replace(each.Receiver, " ", "", -1)
		if each.Receiver != "" {
			m[each.Receiver] = texting(each.Receiver, each.Message)
		}
	}
	c.JSON(http.StatusOK, data)
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
		log.Printf("Error sending message: to %v --> %v\n", to, err)
		return "Error"
	}
	return "Message Sent -> " + to + " : " + msgId
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
		qrChan := make(chan string)
		go func() {
			obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlue, qrcodeTerminal.ConsoleColors.BrightGreen, qrcodeTerminal.QRCodeRecoveryLevels.Low)
			obj.Get(<-qrChan).Print()
		}()
		session, err = wac.Login(qrChan)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	if err = writeSession(session); err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsapp.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&session); err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsapp.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	if err = gob.NewEncoder(file).Encode(session); err != nil {
		return err
	}
	return nil
}

func (h *waHandler) HandleError(err error) {
	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v/n", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		if err := h.c.Restore(); err != nil {
			log.Printf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}
