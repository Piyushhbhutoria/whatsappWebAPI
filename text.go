package main

import (
	"log"

	"github.com/Rhymen/go-whatsapp"
)

func texting(v sendText) string {
	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "91" + v.Receiver + "@s.whatsapp.net",
		},
		Text: v.Message,
	}

	msgID, err := wac.Send(msg)
	if err != nil {
		log.Printf("Error sending message: to %v --> %v\n", v.Receiver, err)
		return "Error"
	}

	return "Message Sent -> " + v.Receiver + " : " + msgID
}
