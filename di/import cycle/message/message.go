package cd_message

import (
	cd_user "github.com/goclub/book/di/import cycle/user"
	"log"
)

func SendMessage(userID string) {
	userName := cd_user.UserName(userID)
	log.Print("SendMessage: Welcome " + userName + "!")
}

func MessageListByUserID(userID string) []string {
	return []string{"Friend request.", "Welcome to join!"}
}