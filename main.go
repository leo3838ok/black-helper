package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

type Vote struct {
	Success    bool   `json:"success"`
	CreatedAt  string `json:"created_at"`
	Left       int    `json:"left"`
	Msg        string `json:"msg"`
	Datetime   string `json:"datetime"`
	Timestamps int    `json:"timestamps"`
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				switch {
				case strings.EqualFold(message.Text, "投起來"):
					var content string

					payload := &bytes.Buffer{}
					writer := multipart.NewWriter(payload)
					fbName := "詹立誠"
					_ = writer.WriteField("stc_candidate_id", "42")
					_ = writer.WriteField("fb_id", "3600467336633412")
					_ = writer.WriteField("fb_name", fbName)
					_ = writer.WriteField("fb_email", "leo3838ok@yahoo.com.tw")
					err := writer.Close()
					if err != nil {
						log.Println(err)
					}

					client := &http.Client{}
					req, err := http.NewRequest("POST", "https://www.mtv.com.tw/api/stc/vote/3", payload)

					if err != nil {
						log.Println(err)
					}
					req.Header.Set("Content-Type", writer.FormDataContentType())

					for {
						res, err := client.Do(req)
						body, err := ioutil.ReadAll(res.Body)

						var vote *Vote
						if err = json.Unmarshal(body, &vote); err != nil {
							log.Println(err)
							break
						}

						if !vote.Success {
							content += fbName + "已完成投票，" + vote.Msg + "\n"
							break
						}

						if err = res.Body.Close(); err != nil {
							log.Println(err)
							break
						}
					}

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
