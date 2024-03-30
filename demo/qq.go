package main

import (
 "crypto/tls"
 //"fmt"
 "log"
 "os"
 "os/signal"
 "syscall"
 "time"

 "github.com/gorilla/websocket"
)

const (
 wsURL = "wss://api.sgroup.qq.com/websocket=S6XyPjeUyfw2anXPpYX3XYm9qj5Pf0ya&t=YOUR_TIMESTAMP"
)

func main() {
 interrupt := make(chan os.Signal, 1)
 signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

 dialer := websocket.Dialer{
 TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 注意：在生产环境中不应该跳过证书验证
 }

 var conn *websocket.Conn
 var err error
 for {
 conn, _, err = dialer.Dial(wsURL, nil)
 if err != nil {
 log.Printf("Failed to connect to WebSocket: %v. Retrying after 5 seconds...\n", err)
 time.Sleep(5 * time.Second)
 continue
 }
 log.Println("Connected to WebSocket.")
 break
 }
 defer conn.Close()

 go func() {
 for {
 mt, message, err := conn.ReadMessage()
 if err != nil {
 log.Printf("Error reading message: %v. Retrying after 5 seconds...\n", err)
 time.Sleep(5 * time.Second)
 continue
 }
 log.Printf("Received message: %s\n", message)

 // 在这里处理接收到的消息
 // ...

 // 响应消息
 err = conn.WriteMessage(mt, []byte("Hello from bot!"))
 if err != nil {
 log.Printf("Error writing message: %v\n", err)
 break
 }
 }

 if err := conn.Close(); err != nil {
 log.Printf("Error closing WebSocket connection: %v\n", err)
 }
 }()

 <-interrupt
 log.Println("Interrupt received, closing WebSocket connection...")
}