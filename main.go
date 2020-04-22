package main

import (
	"bufio"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

type SendText struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
}

type SendImage struct {
	Receiver string `json:"to"`
	Message  string `json:"text"`
	Image    string `json:"image"`
}

var (
	wac, _       = whatsapp.NewConn(20 * time.Second)
	dir, _       = filepath.Abs(filepath.Dir(os.Args[0]))
	folder       string
	textChannel  chan SendText
	imageChannel chan SendImage
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if runtime.GOOS == "windows" {
		folder = `\files\`
	} else {
		folder = "/files/"
	}

	fmt.Println("running on " + strconv.Itoa(runtime.NumCPU()) + " cores.")

	textChannel = make(chan SendText)
	imageChannel = make(chan SendImage)
	wac.SetClientVersion(0, 4, 1307)

	err := login(wac)
	if err != nil {
		panic("Error logging in: \n" + err.Error())
	}

	<-time.After(3 * time.Second)
}

func main() {

	go func() {
		for {
			request, ok := <-textChannel
			if ok {
				log.Println(texting(request))
			}
		}
	}()

	go func() {
		for {
			request, ok := <-imageChannel
			if ok {
				log.Println(image(request))
			}
		}
	}()

	for {
		fmt.Println("Press 0. Test")
		fmt.Println("Press 1. Send Text")
		fmt.Println("Press 2. Send Image")
		fmt.Println("Press 3. Send Bulk Text")
		fmt.Println("Press 4. Send Bulk Image")
		fmt.Println("Press 5. Exit")

		var s int
		fmt.Scanln(&s)
		if s == 0 {
			v := SendText{
				Receiver: "1234567890",
				Message:  "testing",
			}
			log.Println(texting(v))
		} else if s == 1 {
			var v SendText
			fmt.Print("Enter the number: ")
			fmt.Scanln(&v.Receiver)
			fmt.Print("Enter the message: ")
			fmt.Scanln(&v.Message)
			log.Println(texting(v))
		} else if s == 2 {
			var v SendImage
			fmt.Print("Enter the number: ")
			fmt.Scanln(&v.Receiver)
			fmt.Print("Enter the message: ")
			fmt.Scanln(&v.Message)
			fmt.Print("Enter the image name: ")
			fmt.Scanln(&v.Image)
			log.Println(image(v))
		} else if s == 3 {
			var file string
			fmt.Print("Enter the csv name: ")
			fmt.Scanln(&file)
			log.Println(sendBulk(file + ".csv"))
		} else if s == 4 {
			var file string
			fmt.Print("Enter the csv name: ")
			fmt.Scanln(&file)
			log.Println(sendBulkImg(file + ".csv"))
		} else if s == 5 {
			break
		}
	}

}

func sendBulk(file string) string {
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
			v := SendText{
				Receiver: each[0],
				Message:  each[1],
			}
			textChannel <- v
		}
	}

	return "Done"
}

func sendBulkImg(file string) string {
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
			imageChannel <- v
		}
	}

	return "Done"
}

func texting(v SendText) string {
	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "91" + v.Receiver + "@s.whatsapp.net",
		},
		Text: v.Message,
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		log.Printf("Error sending message: to %v --> %v\n", v.Receiver, err)
		return "Error"
	}

	return "Message Sent -> " + v.Receiver + " : " + msgId
}

func image(v SendImage) string {
	var folder string

	img, err := os.Open(dir + folder + v.Image)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
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
		log.Printf("Error sending message: to %v --> %v\n", v.Receiver, err)
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
