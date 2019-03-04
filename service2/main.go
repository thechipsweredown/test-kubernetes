package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Server struct {
	mux     map[string]func(http.ResponseWriter, *http.Request)
	httpSer *http.Server
}
func GetTopic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hihih")
	response,err := http.Get("http://127.0.0.1:9099")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		w.Write([]byte(string(contents)))
	}
}
func (server *Server) InitMux() {
	server.mux = make(map[string]func(http.ResponseWriter, *http.Request))
	server.mux["/"] = GetTopic
}

func (server *Server) Run() {
	server.InitMux()
	withGz := server
	server.httpSer = &http.Server{
		Addr:         "127.0.0.1:3000",
		Handler:      withGz,
		ReadTimeout:  500 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Server is listening at 3000...")
	server.httpSer.ListenAndServe()
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := server.mux[r.URL.Path]; ok {
		h(w, r)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("<b><font size=\"6\">Bad request</font></b>"))
}

func main() {
	runtime.GOMAXPROCS(5)
	server := &Server{}
	chan_sync := make(chan bool)
	go server.Run()
	<- chan_sync
}