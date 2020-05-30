package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

var (
	wac, _       = whatsapp.NewConn(20 * time.Second)
	dir, _       = filepath.Abs(filepath.Dir(os.Args[0]))
	textChannel  chan sendText
	imageChannel chan sendImage
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("running on " + strconv.Itoa(runtime.NumCPU()) + " cores.")

	textChannel = make(chan sendText)
	imageChannel = make(chan sendImage)
	wac.SetClientVersion(2, 2021, 4)

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
			v := sendText{
				Receiver: "1234567890",
				Message:  "testing",
			}
			log.Println(texting(v))
		} else if s == 1 {
			var v sendText
			fmt.Print("Enter the reciever number: ")
			fmt.Scanln(&v.Receiver)
			fmt.Print("Enter the message to be sent: ")
			fmt.Scanln(&v.Message)
			log.Println(texting(v))
		} else if s == 2 {
			var v sendImage
			fmt.Print("Enter the reciever number: ")
			fmt.Scanln(&v.Receiver)
			fmt.Print("Enter the message to be sent: ")
			fmt.Scanln(&v.Message)
			fmt.Print("Enter the image name: ")
			fmt.Scanln(&v.Image)
			log.Println(image(v))
		} else if s == 3 {
			var file string
			fmt.Print("Enter the csv file name: ")
			fmt.Scanln(&file)
			log.Println(sendBulk(file + ".csv"))
		} else if s == 4 {
			var file string
			fmt.Print("Enter the csv name: ")
			fmt.Scanln(&file)
			log.Println(sendBulkImg(file + ".csv"))
		} else if s == 5 {
			log.Println("Application exiting...")
			break
		}
	}
}
