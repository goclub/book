package cd_user

import (
	cd_message "github.com/goclub/book/di/import cycle/message"
)

func UserName(userID string) string {
	return "nimoc"
}

func MyMessageList() []string {
	return cd_message.MessageListByUserID("a")
}