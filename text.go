package main

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func text(args []string) {
	recipient, ok := parseJID(args[0])
	if !ok {
		return
	}
	check := checkuser(args)
	if check {
		msg := &waE2E.Message{Conversation: proto.String(strings.Join(args[1:], " "))}
		ts, err := cli.SendMessage(context.Background(), recipient, msg)
		if err != nil {
			log.Errorf("Error sending message: %v", err)
		} else {
			log.Infof("Message sent (server timestamp: %s)", ts)
		}
	} else {
		log.Errorf("User doesn't exist: %v", args[0])
	}
}
