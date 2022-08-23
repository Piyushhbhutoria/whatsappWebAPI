package main

func checkuser(args []string) bool {
	resp, err := cli.IsOnWhatsApp(args)
	if err != nil {
		log.Errorf("Failed to check if users are on WhatsApp:", err)
		return false
	}
	for _, item := range resp {
		if item.VerifiedName != nil {
			log.Infof("%s: on whatsapp: %t, JID: %s, business name: %s", item.Query, item.IsIn, item.JID, item.VerifiedName.Details.GetVerifiedName())
		} else {
			log.Infof("%s: on whatsapp: %t, JID: %s", item.Query, item.IsIn, item.JID)
		}
		return item.IsIn
	}

	return false
}
