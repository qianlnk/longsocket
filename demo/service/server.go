package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/qianlnk/longsocket"
	"golang.org/x/net/websocket"
)

func testHander(ws *websocket.Conn) {
	req := ws.Request()
	fmt.Println(req)
	u, err := url.Parse(req.Header.Get("Origin"))
	if err != nil {
		ws.Close()
		return
	}

	user := u.Query().Get("user")
	password := u.Query().Get("password")
	fmt.Println(user, password)
	mysocket := longsocket.NewConn("", "", "", true, 128*1024)
	mysocket.SetSocket(ws)
	defer mysocket.Close()
	go mysocket.WriteLoop()
	go mysocket.ReadLoop()
	mysocket.Read(testdealmsg)
}

func testdealmsg(msg []byte, l *longsocket.Longsocket) error {
	fmt.Println(string(msg))
	return nil
}

func main() {
	http.Handle("/test", websocket.Handler(testHander))
	// initialize server
	srv := &http.Server{
		Addr:           ":1234",
		Handler:        nil,
		ReadTimeout:    time.Duration(5) * time.Minute,
		WriteTimeout:   time.Duration(5) * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	// start listen
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
