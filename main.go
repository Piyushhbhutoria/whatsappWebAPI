package main

import (
	"bufio"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

var (
	wac, _ = whatsapp.NewConn(5 * time.Second)
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	folder string
)

func init() {
	if runtime.GOOS == "windows" {
		folder = `\files\`
	} else {
		folder = "/files/"
	}

	fmt.Println("running on " + string(runtime.NumCPU()) + "cores.")

	err := login(wac)
	if err != nil {
		panic("Error logging in: \n" + err.Error())
		return
	}

	<-time.After(3 * time.Second)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	router := gin.New()
	router.Use(gin.Recovery())
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
	router.Use(cors.Default())

	router.GET("/ping", ping)
	router.GET("/sendText", sendText)
	router.GET("/sendBulk", sendBulk)
	router.GET("/sendImage", sendImage)
	router.GET("/sendBulkImg", sendBulkImg)

	if err := router.Run(":8081"); err != nil {
		log.Printf("Shutdown with error: %v\n", err)
	}
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func sendText(c *gin.Context) {
	to := strings.Replace(c.DefaultQuery("to", "1234567890"), " ", "", -1)
	mess := c.DefaultQuery("msg", "testing")
	c.String(http.StatusOK, texting(to, mess))
}

func sendBulk(c *gin.Context) {
	file := c.DefaultQuery("file", "test.csv")
	var folder string
	m := make(map[string]string)

	if runtime.GOOS == "windows" {
		folder = `\files\`
	} else {
		folder = "/files/"
	}

	csvFile, err := os.Open(dir + folder + file)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, each := range csvData {
		each[0] = strings.Replace(each[0], " ", "", -1)
		if each[0] != "" {
			m[each[0]] = texting(each[0], each[1])
		}
	}

	c.JSON(http.StatusOK, m)
}

func sendImage(c *gin.Context) {
	to := strings.Replace(c.DefaultQuery("to", "1234567890"), " ", "", -1)
	mess := c.DefaultQuery("msg", "testing")
	img := c.DefaultQuery("img", "testImg.jpg")
	v := SendImage{
		Receiver: to,
		Message:  mess,
		Image:    img,
	}

	c.String(http.StatusOK, image(v))
}

func sendBulkImg(c *gin.Context) {
	file := c.DefaultQuery("file", "testImg.csv")
	var folder string
	m := make(map[string]string)

	csvFile, err := os.Open(dir + folder + file)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, each := range csvData {
		if each[0] != "" {
			each[0] = strings.Replace(each[0], " ", "", -1)
			v := SendImage{
				Receiver: each[0],
				Message:  each[1],
				Image:    each[2],
			}
			m[each[0]] = image(v)
		}
	}

	c.JSON(http.StatusOK, m)
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
	fmt.Print("Enter your number: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	fmt.Println("Logging in -> " + text)
	//load saved session
	session, err := readSession(text)
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
	err = writeSession(session, text)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession(s string) (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(dir + folder + s + ".gob")
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

func writeSession(session whatsapp.Session, s string) error {
	file, err := os.Create(dir + folder + s + ".gob")
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
