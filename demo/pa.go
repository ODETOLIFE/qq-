package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

// MessageRequest 是发送到QQ频道机器人的消息请求结构
type MessageRequest struct {
    Action string `json:"action"`
    Data   struct {
        OpenID string `json:"open_id"` // 接收消息的用户的OpenID
        Msg    struct {
            Type int    `json:"type"` // 消息类型，例如0表示文本消息
            Text string `json:"text"` // 文本消息内容
        } `json:"msg"`
    } `json:"data"`
}

// SendMessage 发送消息到QQ频道
func SendMessage(token, openID, message string) error {
    requestBody := MessageRequest{
        Action: "send_msg",
        Data: struct {
            OpenID string `json:"open_id"`
            Msg    struct {
                Type int    `json:"type"`
                Text string `json:"text"`
            } `json:"msg"`
        }{
            OpenID: openID,
            Msg: struct {
                Type int    `json:"type"`
                Text string `json:"text"`
            }{
                Type: 0,
                Text: message,
            },
        },
    }

    // 将请求体转换为JSON
    requestBodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        return err
    }

    // 创建请求
    req, err := http.NewRequest("POST", "https://api.q.qq.com/cgi-bin/ws/send_msg?access_token="+token, bytes.NewBuffer(requestBodyBytes))
    if err != nil {
        return err
    }

    // 设置请求头
    req.Header.Add("Content-Type", "application/json")

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 读取响应
    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // 输出响应
    fmt.Println(string(responseBody))

    return nil
}

func main() {
    // TODO: 替换为您的访问令牌和接收消息的用户的OpenID
    token := "S6XyPjeUyfw2anXPpYX3XYm9qj5Pf0ya"
    openID := "102098741"
    message := "Hello, QQ Channel Bot!"

    // 发送消息
    err := SendMessage(token, openID, message)
    if err != nil {
        fmt.Println("Error sending message:", err)
    } else {
        fmt.Println("Message sent successfully!")
    }
}