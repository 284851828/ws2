package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"gtest/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	// go ws.WebsocketManager.Start()
	// go ws.WebsocketManager.SendService()
	// go ws.WebsocketManager.SendService()
	// go ws.WebsocketManager.SendGroupService()
	// go ws.WebsocketManager.SendGroupService()
	// go ws.WebsocketManager.SendAllService()
	// go ws.WebsocketManager.SendAllService()
	wsIO := ws.WsIoBase{}
	wsIO.Init("chan", "")
	ws.WebsocketManager.ReceiveMsgHandle = wsIO.NewMsg
	go ws.WebsocketManager.StartService()
	go ws.TestSendGroup()
	go ws.TestSendAll()

	go wsIO.ReadFromws(doMsg)
	var node ws.NodeWrite
	wsIO.Write2ws(node)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	wsGroup := router.Group("/ws")
	{
		wsGroup.GET("/:channel", ws.WebsocketManager.WsClient)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Start Error: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Error:", err)
	}
	log.Println("Server Shutdown")
}

// server  do msg
var __packcount int64

func doMsg(msg ws.NodeRead) {
	__packcount += 1
	if __packcount%1000 == 0 {
		fmt.Printf("receive %d package \n", __packcount)
	}
	// fmt.Printf("doMsg id:%s, group:%s, msglen:%d\n", msg.Uid, msg.Gid, len(string(msg.Msg.([]byte))))
}

/////////////////////////////////////////////////////////////////////////////////////
// client

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func testServer() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/leffss"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 5)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}
