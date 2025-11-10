package main

import (
	"context"

	"go.mau.fi/whatsmeow/types"
)

func findUsers(ctx context.Context, args []string) (bool, types.IsOnWhatsAppResponse) {
	resp, err := cli.IsOnWhatsApp(ctx, args)
	if err != nil {
		log.Errorf("Failed to check if users are on WhatsApp:", err)
		return false, types.IsOnWhatsAppResponse{}
	}
	if len(resp) == 0 {
		log.Infof("No results")
		return false, types.IsOnWhatsAppResponse{}
	}

	item := resp[0]
	if item.VerifiedName != nil {
		log.Infof("%s: on whatsapp: %t, JID: %s, business name: %s", item.Query, item.IsIn, item.JID, item.VerifiedName.Details.GetVerifiedName())
	} else {
		log.Infof("%s: on whatsapp: %t, JID: %s", item.Query, item.IsIn, item.JID)
	}
	return true, item
}
