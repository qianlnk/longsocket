# longsocket
websocket long connect
# use
### server
do some thing in func `testdealmsg` and return message by func `Write`.
```golang
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
```

### client
keep long connect by func `keeplongconnect` with gorouting. and send message to server by func `Write`. deal server's res message in func `testdealmsg`, at the same time, you can send message to server also.
```golang
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
```
