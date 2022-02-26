package main

import (
	"strings"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

func text(args []string) {
	recipient, ok := parseJID(args[0])
	if !ok {
		return
	}
	check := checkuser(args)
	if check {
		msg := &waProto.Message{Conversation: proto.String(strings.Join(args[1:], " "))}
		ts, err := cli.SendMessage(recipient, "", msg)
		if err != nil {
			log.Errorf("Error sending message: %v", err)
		} else {
			log.Infof("Message sent (server timestamp: %s)", ts)
		}
	} else {
		log.Errorf("User doesn't exist: %v", args[0])
	}
}
