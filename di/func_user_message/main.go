package main

import (
	func_message "github.com/goclub/book/di/func_user_message/message"
	func_user "github.com/goclub/book/di/func_user_message/user"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		dep := func_user.UserName
		func_message.SendMessage("a", dep)
		_, err := writer.Write([]byte(`<a href="/message" >message list</a>`)) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	http.HandleFunc("/message", func(writer http.ResponseWriter, request *http.Request) {
		dep := func_message.MessageListByUserID
		messageList := func_user.MyMessageList(dep)
		data := "message list" + strings.Join(messageList, "\n")
		_, err := writer.Write([]byte(data)) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	log.Print("http://127.0.0.1:1219")
	err := http.ListenAndServe(":1219", nil) ; if err != nil {
		panic(err)
	}
}
