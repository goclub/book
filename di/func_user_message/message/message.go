package func_message

import (
	"log"
)

func SendMessage(userID string, UserName func(userID string) string ) {
	userName := UserName(userID)
	log.Print("SendMessage: Welcome " + userName + "!")
}

func MessageListByUserID(userID string) []string {
	return []string{"Friend request.", "Welcome to join!"}
}