package main

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"os"
	"path/filepath"
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
)

func init() {
	err := login(wac)
	if err != nil {
		panic("Error logging in: %v\n" + err.Error())
		return
	}

	<-time.After(3 * time.Second)
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(logger.New())
	// using Cors
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})
	v1 := app.Party("/", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Get("/ping", ping) // Test url
		v1.Get("/sendText", sendText)
		v1.Get("/sendBulk", sendBulk)
		v1.Get("/sendBulkImg", sendBulkImg)
		v1.Get("/sendImage", sendImage)
		v1.Get("/testText", testText)
		v1.Get("/testImage", testImage)
	}
	err := app.Run(iris.Addr(":8080"), iris.WithOptimizations, iris.WithoutBanner, iris.WithoutStartupLog)
	if err != iris.ErrServerClosed {
		panic("Shutdown with error: " + err.Error())
	}
}

func ping(ctx iris.Context) {
	ctx.WriteString("pong")
}

func sendText(ctx iris.Context) {
	to := strings.Replace(ctx.URLParam("to"), " ", "", -1)
	mess := ctx.URLParam("msg")
	ctx.WriteString(texting(to, mess))
}

func sendBulk(ctx iris.Context) {
	file := ctx.URLParamDefault("file", "test.csv")

	csvFile, err := os.Open(dir + "/files/" + file)
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
	m := make(map[string]string)

	for _, each := range csvData {
		each[0] = strings.Replace(each[0], " ", "", -1)
		if each[0] != "" {
			m[each[0]] = texting(each[0], each[1])
		}
	}

	ctx.JSON(m)
}

func sendBulkImg(ctx iris.Context) {
	file := ctx.URLParamDefault("file", "testImg.csv")

	csvFile, err := os.Open(dir + "/files/" + file)
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
	m := make(map[string]string)

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

	ctx.JSON(m)
}

func sendImage(ctx iris.Context) {
	to := strings.Replace(ctx.URLParam("to"), " ", "", -1)
	mess := ctx.URLParam("msg")
	img := ctx.URLParam("img")

	v := SendImage{
		Receiver: to,
		Message:  mess,
		Image:    img,
	}

	ctx.WriteString(image(v))
}

func testText(ctx iris.Context) {
	to := strings.Replace(ctx.URLParamDefault("to", "1234567890"), " ", "", -1)
	mess := ctx.URLParamDefault("msg", "testing")
	ctx.WriteString(texting(to, mess))
}

func testImage(ctx iris.Context) {
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

func texting(to, mess string) string {
	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "91" + to + "@s.whatsapp.net",
		},
		Text: mess,
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		panic("Error sending message: %v" + err.Error())
		return "Error"
	}
	return "Message Sent -> " + to + " : " + msgId
}

func image(v SendImage) string {
	img, err := os.Open(dir + "/files/" + v.Image)
	if err != nil {
		panic("Error reading file: %v" + err.Error())
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
		panic("Error sending message: %v" + err.Error())
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
