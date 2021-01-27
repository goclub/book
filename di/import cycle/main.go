package main

import (
	cd_message "github.com/goclub/book/di/import cycle/message"
	cd_user "github.com/goclub/book/di/import cycle/user"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		cd_message.SendMessage("a")
		_, err := writer.Write([]byte(`<a href="/message" >message list</a>`)) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	http.HandleFunc("/message", func(writer http.ResponseWriter, request *http.Request) {
		messageList := cd_user.MyMessageList()
		data := "message list" + strings.Join(messageList, "\n")
		_, err := writer.Write([]byte(data)) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	err := http.ListenAndServe(":1219", nil) ; if err != nil {
		panic(err)
	}
}
