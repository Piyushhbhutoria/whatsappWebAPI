package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func image(ctx context.Context, args []string) {
	check, item := findUsers(ctx, args)
	if check && item.IsIn {
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Errorf("Failed to read %s: %v", args[1], err)
			return
		}
		uploaded, err := cli.Upload(ctx, data, whatsmeow.MediaImage)
		if err != nil {
			log.Errorf("Failed to upload file: %v", err)
			return
		}
		msg := &waE2E.Message{ImageMessage: &waE2E.ImageMessage{
			Caption:       proto.String(strings.Join(args[2:], " ")),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
		}}
		ts, err := cli.SendMessage(ctx, item.JID, msg)
		if err != nil {
			log.Errorf("Error sending image message: %v", err)
		} else {
			log.Infof("Image message sent (server timestamp: %s)", ts)
		}
	} else {
		log.Errorf("User %s doesn't exist", item.Query)
	}
}
