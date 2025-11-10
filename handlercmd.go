package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"
)

func handleCmd(cmd string, args []string) {
	ctx := context.Background()
	switch cmd {
	case "reconnect":
		cli.Disconnect()
		err := cli.Connect()
		if err != nil {
			log.Errorf("Failed to connect: %v", err)
		}
	case "logout":
		err := cli.Logout(ctx)
		if err != nil {
			log.Errorf("Error logging out: %v", err)
		} else {
			log.Infof("Successfully logged out")
		}
	case "appstate":
		if len(args) < 1 {
			log.Errorf("Usage: appstate <types...>")
			return
		}
		names := []appstate.WAPatchName{appstate.WAPatchName(args[0])}
		if args[0] == "all" {
			names = []appstate.WAPatchName{appstate.WAPatchRegular, appstate.WAPatchRegularHigh, appstate.WAPatchRegularLow, appstate.WAPatchCriticalUnblockLow, appstate.WAPatchCriticalBlock}
		}
		resync := len(args) > 1 && args[1] == "resync"
		for _, name := range names {
			err := cli.FetchAppState(ctx, name, resync, false)
			if err != nil {
				log.Errorf("Failed to sync app state: %v", err)
			}
		}
	case "checkuser":
		if len(args) < 1 {
			log.Errorf("Usage: checkuser <phone numbers...>")
			return
		}
		findUsers(ctx, args)
	case "subscribepresence":
		if len(args) < 1 {
			log.Errorf("Usage: subscribepresence <jid>")
			return
		}
		jid, err := types.ParseJID(args[0])
		if err != nil {
			log.Errorf("Invalid JID %s: %v", args[0], err)
			return
		}
		err = cli.SubscribePresence(ctx, jid)
		if err != nil {
			fmt.Println(err)
		}
	case "presence":
		fmt.Println(cli.SendPresence(ctx, types.Presence(args[0])))
	case "chatpresence":
		jid, _ := types.ParseJID(args[1])
		fmt.Println(cli.SendChatPresence(ctx, jid, types.ChatPresence(args[0]), types.ChatPresenceMedia(args[2])))
	case "privacysettings":
		resp, err := cli.TryFetchPrivacySettings(ctx, false)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%+v\n", resp)
		}
	case "getuser":
		if len(args) < 1 {
			log.Errorf("Usage: getuser <jids...>")
			return
		}
		var jids []types.JID
		for _, arg := range args {
			jid, err := types.ParseJID(arg)
			if err != nil {
				log.Errorf("Invalid JID %s: %v", arg, err)
				return
			}
			jids = append(jids, jid)
		}
		resp, err := cli.GetUserInfo(ctx, jids)
		if err != nil {
			log.Errorf("Failed to get user info: %v", err)
		} else {
			for jid, info := range resp {
				log.Infof("%s: %+v", jid, info)
			}
		}
	case "getavatar":
		if len(args) < 1 {
			log.Errorf("Usage: getavatar <jid>")
			return
		}
		jid, err := types.ParseJID(args[0])
		if err != nil {
			log.Errorf("Invalid JID %s: %v", args[0], err)
			return
		}
		params := &whatsmeow.GetProfilePictureParams{
			Preview: len(args) > 1 && args[1] == "preview",
		}
		pic, err := cli.GetProfilePictureInfo(ctx, jid, params)
		if err != nil {
			log.Errorf("Failed to get avatar: %v", err)
		} else if pic != nil {
			log.Infof("Got avatar ID %s: %s", pic.ID, pic.URL)
		} else {
			log.Infof("No avatar found")
		}
	case "getgroup":
		if len(args) < 1 {
			log.Errorf("Usage: getgroup <jid>")
			return
		}
		group, err := types.ParseJID(args[0])
		if err != nil {
			log.Errorf("Invalid JID %s: %v", args[0], err)
			return
		}
		resp, err := cli.GetGroupInfo(ctx, group)
		if err != nil {
			log.Errorf("Failed to get group info: %v", err)
		} else {
			log.Infof("Group info: %+v", resp)
		}
	case "listgroups":
		groups, err := cli.GetJoinedGroups(ctx)
		if err != nil {
			log.Errorf("Failed to get group list: %v", err)
		} else {
			for _, group := range groups {
				log.Infof("%+v", group)
			}
		}
	case "getinvitelink":
		if len(args) < 1 {
			log.Errorf("Usage: getinvitelink <jid> [--reset]")
			return
		}
		group, err := types.ParseJID(args[0])
		if err != nil {
			log.Errorf("Invalid JID %s: %v", args[0], err)
			return
		} else if group.Server != types.GroupServer {
			log.Errorf("Input must be a group JID (@%s)", types.GroupServer)
			return
		}
		resp, err := cli.GetGroupInviteLink(ctx, group, len(args) > 1 && args[1] == "--reset")
		if err != nil {
			log.Errorf("Failed to get group invite link: %v", err)
		} else {
			log.Infof("Group invite link: %s", resp)
		}
	case "queryinvitelink":
		if len(args) < 1 {
			log.Errorf("Usage: queryinvitelink <link>")
			return
		}
		resp, err := cli.GetGroupInfoFromLink(ctx, args[0])
		if err != nil {
			log.Errorf("Failed to resolve group invite link: %v", err)
		} else {
			log.Infof("Group info: %+v", resp)
		}
	case "querybusinesslink":
		if len(args) < 1 {
			log.Errorf("Usage: querybusinesslink <link>")
			return
		}
		resp, err := cli.ResolveBusinessMessageLink(ctx, args[0])
		if err != nil {
			log.Errorf("Failed to resolve business message link: %v", err)
		} else {
			log.Infof("Business info: %+v", resp)
		}
	case "joininvitelink":
		if len(args) < 1 {
			log.Errorf("Usage: acceptinvitelink <link>")
			return
		}
		groupID, err := cli.JoinGroupWithLink(ctx, args[0])
		if err != nil {
			log.Errorf("Failed to join group via invite link: %v", err)
		} else {
			log.Infof("Joined %s", groupID)
		}
	case "send":
		if len(args) < 2 {
			log.Errorf("Usage: send <jid> <text>")
			return
		}
		text(ctx, args)
	case "sendbulk":
		if len(args) < 1 {
			log.Errorf("Usage: sendbulk <csv file>")
			return
		}
		csvFile, err := os.Open(filepath.Join(args[0]))
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
				text(ctx, each)
			}
		}
	case "sendimg":
		if len(args) < 2 {
			log.Errorf("Usage: sendimg <jid> <image path> [caption]")
			return
		}
		image(ctx, args)
	case "sendbulkimg":
		if len(args) < 1 {
			log.Errorf("Usage: sendbulk <csv file>")
			return
		}
		csvFile, err := os.Open(filepath.Join(args[0]))
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
			image(ctx, each)
		}
	}
}
