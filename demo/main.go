package main

import (
    "bytes"
	"time"
	"context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
	"strings"
	"log"
    "os"
    "github.com/tencent-connect/botgo"
    "github.com/tencent-connect/botgo/dto"
    "github.com/tencent-connect/botgo/openapi"
    "github.com/tencent-connect/botgo/token"
    "github.com/tencent-connect/botgo/websocket"
    "github.com/tencent-connect/botgo/event"
   // yaml "gopkg.in/yaml.v2"
)
var api openapi.OpenAPI
var ctx context.Context

type Type string

// TokenType
const (
	TypeBot    Type = "Bot"
	TypeNormal Type = "Bearer"
)
//type Token struct {
//	AppID       uint64
//	AccessToken string
//	Type        Type
//}
// 定义API请求的结构体
type APIRequest struct {
    Action string `json:"action"`
    Data   struct {
        OpenID string `json:"open_id"`
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"data"`
}
func atMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
    if strings.HasSuffix(data.Content, "> hello") { // 如果@机器人并输入 hello 则回复 你好。
        api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: "你好"})
    }
    return nil
}
// 发送消息到QQ频道的函数
func SendMessageToQQChannel(url, openID, messageContent, token1 string) error {
    // 创建API请求结构体
    request := APIRequest{
        Action: "send_message",
        Data: struct {
            OpenID string `json:"open_id"`
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        }{
            OpenID: openID,
            Message: struct {
                Content string `json:"content"`
            }{
                Content: messageContent,
            },
        },
    }

    // 将请求结构体编码为JSON
    jsonBytes, err := json.Marshal(request)
    if err != nil {
        return err
    }

    // 创建HTTP客户端
    client := &http.Client{}

    // 创建HTTP请求
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
    if err != nil {
        return err
    }

    // 设置请求头部，包含令牌
    req.Header.Set("Authorization", "Bearer "+token1)
    req.Header.Set("Content-Type", "application/json")

    // 发送请求
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 读取响应
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // 打印响应内容
    fmt.Println(string(body))

    // 检查响应状态码
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
    }

    return nil
}

func main() {
    apiURL := "wss://api.sgroup.qq.com/websocket"
    openID := "102098741"
    messageContent := "Hello, QQ Channel!"
    token1 := "S6XyPjeUyfw2anXPpYX3XYm9qj5Pf0ya" // 替换为你的令牌

    if err := SendMessageToQQChannel(apiURL, openID, messageContent, token1); err != nil {
        fmt.Printf("Error sending message to QQ Channel: %s\n", err)
    } else {
        fmt.Println("Message sent successfully!")
    }
	token := token.BotToken(102098741, "S6XyPjeUyfw2anXPpYX3XYm9qj5Pf0ya") 
    //第三步：获取操作机器人的API对象
    api = botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
    //获取context
    ctx = context.Background()
	ws, err := api.WS(ctx, nil, "") 
    if err != nil {
        log.Fatalln("websocket错误， err = ", err)
        os.Exit(1)
    }
	////var token1 token.Token
	//type Token struct {
//	AppID       uint64
//	AccessToken string
//	Type        Type
//}
// 定义API请求的结构体
//	(*token1).AppID = 102098741
//	(*token1).AccessToken = "S6XyPjeUyfw2anXPpYX3XYm9qj5Pf0ya"
//	(*token1).Type = TypeBot
	
	var atMessage event.ATMessageEventHandler = atMessageEventHandler
	intent := websocket.RegisterHandlers(atMessage)     // 注册socket消息处理
    botgo.NewSessionManager().Start(ws, token, &intent) // 启动socket监听
}