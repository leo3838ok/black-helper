package main

import (
	"bytes"
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
					payload := &bytes.Buffer{}
					writer := multipart.NewWriter(payload)
					_ = writer.WriteField("stc_candidate_id", "42")
					_ = writer.WriteField("fb_id", "3600467336633412")
					_ = writer.WriteField("fb_name", "詹立誠")
					_ = writer.WriteField("fb_email", "leo3838ok@yahoo.com.tw")
					err := writer.Close()
					if err != nil {
						log.Println(err)
					}

					client := &http.Client {
					}
					req, err := http.NewRequest("POST", "https://www.mtv.com.tw/api/stc/vote/3", payload)

					if err != nil {
						log.Println(err)
					}
					req.Header.Set("Content-Type", writer.FormDataContentType())
					res, err := client.Do(req)
					defer res.Body.Close()
					body, err := ioutil.ReadAll(res.Body)

					content := string(body)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}