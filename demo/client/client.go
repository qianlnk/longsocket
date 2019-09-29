package main

import (
	"fmt"
	"time"

	"github.com/qianlnk/longsocket"
)

var mysocket *longsocket.Longsocket

func testdealmsg(msg []byte, l *longsocket.Longsocket) error {
	fmt.Println(string(msg))
	return nil
}

func keeplongconnect() {
	for {
		wsAddr := fmt.Sprintf("ws://127.0.0.1:1234/test")
		httpAddr := fmt.Sprintf("http://127.0.0.1:1234/test?user=%s&pwd=%s", "qianlnk", "123456")
		mysocket = longsocket.NewConn(wsAddr, "", httpAddr, true, 128*1024)
		err := mysocket.Dial(true)
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		defer mysocket.Close()
		go mysocket.WriteLoop()
		go mysocket.ReadLoop()
		mysocket.Read(testdealmsg)
		time.Sleep(2 * time.Second)
	}
}

func main() {
	go keeplongconnect()
	time.Sleep(2 * time.Second) //wait for connect
	for {
		mysocket.Write([]byte("test"))
		time.Sleep(2 * time.Second)
	}
}
