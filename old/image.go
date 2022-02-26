package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Rhymen/go-whatsapp"
)

func image(v sendImage) string {
	img, err := os.Open(filepath.Join(dir, v.Image))
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

	msgID, err := wac.Send(msg)
	if err != nil {
		log.Printf("Error sending message: to %v --> %v\n", v.Receiver, err)
		return "Error"
	}

	return "Message Sent -> " + v.Receiver + " : " + msgID
}
