package main

import (
	"log"
	"strings"

	"github.com/TwiN/lazywebhooks/discord"
)

var lastMessageSent string

func SendNotification(report Report) {
	if len(report.Problems) == 0 {
		return
	}
	sb := strings.Builder{}
	for _, problem := range report.Problems {
		sb.WriteString("**" + problem.Summary + "**")
		sb.WriteString("\n")
		if len(problem.Description) > 0 {
			sb.WriteString(problem.Description)
		}
		sb.WriteString("\n")
	}
	message := sb.String()
	if message == lastMessageSent {
		if debug {
			log.Print("[SendNotification] Skipping because it's the same as the last message sent")
		}
		return
	}
	log.Print("[SendNotification] Sending notification to Discord")
	discord.Send(message, webhookURL)
	lastMessageSent = message
}
