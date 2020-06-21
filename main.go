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
	"sort"
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
					file, _ := ioutil.ReadFile("fb.json")
					fbInfos := FBInfos{}
					_ = json.Unmarshal([]byte(file), &fbInfos)

					content := vote(fbInfos)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
						log.Println(err)
					}
				case strings.EqualFold(message.Text, "報戰況"):
					content, err := list()
					if err != nil {
						log.Println(err)
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func vote(infos FBInfos) string {
	var content string
	var voteCount int

	for _, info := range infos {
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("stc_candidate_id", "42")
		_ = writer.WriteField("fb_id", info.FbID)
		_ = writer.WriteField("fb_name", info.FbName)
		_ = writer.WriteField("fb_email", info.FbEmail)
		err := writer.Close()
		if err != nil {
			log.Println(err)
			continue
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://www.mtv.com.tw/api/stc/vote/3", payload)

		if err != nil {
			log.Println(err)
			continue
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		for {
			res, err := client.Do(req)
			if err != nil {
				log.Println(err)
				break
			}

			body, err := ioutil.ReadAll(res.Body)

			var vote *Vote
			if err = json.Unmarshal(body, &vote); err != nil {
				log.Println(err)
				break
			}

			if !vote.Success {
				content += info.FbName + "已完成投票，" + vote.Msg + "\n"
				if err = res.Body.Close(); err != nil {
					log.Println(err)
				}
				break
			}

			voteCount++

			if err = res.Body.Close(); err != nil {
				log.Println(err)
				break
			}
		}
	}

	if voteCount > 0 {
		content += fmt.Sprintf("感謝大恩大德多賜了%v票\n", voteCount)
	}

	return content
}

func list() (string, error) {
	var content string

	client := &http.Client {}
	req, err := http.NewRequest("GET", "https://www.mtv.com.tw/api/stc/listall/3", nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var candidates *Candidates
	if err = json.Unmarshal(body, &candidates); err != nil {
		return "", err
	}

	sort.SliceStable(candidates.Data, func(i, j int) bool {
		return candidates.Data[i].VotesCount > candidates.Data[j].VotesCount
	})

	for i, v := range candidates.Data {
		if v.ID == 42 {
			content += fmt.Sprintf("目前暫居第%v名，", i + 1)
			if i == 0 {
				diff := v.VotesCount - candidates.Data[i + 1].VotesCount
				content += fmt.Sprintf("贏過第2名%v票\n", diff)
			} else {
				diff := candidates.Data[i - 1].VotesCount - v.VotesCount
				content += fmt.Sprintf("距離第%v名%v票\n", i, diff)
			}
		}
	}

	return content, nil
}