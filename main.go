package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	// 環境変数からトークンとシークレットを取得
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	// LINE Bot SDK クライアント作成
	bot, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		log.Fatal("LINE bot client creation error:", err)
	}

	// 通常の確認用ページ
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!! Goだよ!")
	})

	// LINE Webhook エンドポイント
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		// 署名の検証とイベントのパース
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				http.Error(w, "Invalid signature", http.StatusBadRequest)
			} else {
				http.Error(w, "Parse error", http.StatusInternalServerError)
			}
			return
		}

		// すべてのイベントを処理
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				if _, ok := event.Message.(*linebot.TextMessage); ok {
					// メッセージが来たら「test」と返信
					if _, err := bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("test"),
					).Do(); err != nil {
						log.Println("Reply error:", err)
					}
					log.Printf("Replied 'test' to user: %s\n", event.Source.UserID)
				}
			}
		}
		// 200 OK を返す
		io.WriteString(w, "OK")
	})

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル用デフォルト
	}
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
