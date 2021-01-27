package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(time.Now().String())) ; if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/panic", func(writer http.ResponseWriter, request *http.Request) {
		go func() {
			panic(1)
		}()
	})
	addr := ":3000"
	serve := http.Server{
		Addr: addr,
		Handler: mux,
	}
	go func() {
		log.Print("listen: http://127.0.0.1" + addr)
		err := serve.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	log.Print("Shuting down server...")
	if err := serve.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
	go func() {
		<-exit
	}()
}
