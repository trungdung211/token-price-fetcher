package main

import (
	"log"

	"github.com/gtuk/discordwebhook"
)

func main() {
	var username = "BotUser"
	var content = "☀ This is a test <strong>message</strong>"
	var url = "https://discord.com/api/webhooks/1083782248412753941/0S5RPsN0jR00cUa7ZlHL_zHzR01QTeqoZHqCJ-bv_on9jOSzOU6X26pKqsuSe8hvwoFH"

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(url, message)
	if err != nil {
		log.Fatal(err)
	}
}
