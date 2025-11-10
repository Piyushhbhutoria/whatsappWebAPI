package main

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func text(ctx context.Context, args []string) {
	check, item := findUsers(ctx, args)
	if check && item.IsIn {
		msg := &waE2E.Message{Conversation: proto.String(strings.Join(args[1:], " "))}
		ts, err := cli.SendMessage(ctx, item.JID, msg)
		if err != nil {
			log.Errorf("Error sending message: %v", err)
		} else {
			log.Infof("Message sent (server timestamp: %s)", ts)
		}
	} else {
		log.Errorf("User %s doesn't exist", item.Query)
	}
}
