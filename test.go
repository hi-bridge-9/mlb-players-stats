package function

import (
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func FunctionTest(w http.ResponseWriter, r *http.Request) {
	// bytes, err := ioutil.ReadFile("../data/players_list.json")
	// if err != nil {
	// 	log.Fatalf("downloadFileIntoMemory: %w", err)
	// }

	// var ps []*Player
	// if err := json.Unmarshal(bytes, &ps); err != nil {
	// 	log.Fatalf("error in unmarshal json: %w", err)
	// }

	var c *Crawler
	// newsList, err := c.GetUpdateInfo(ps)
	// if err != nil {
	// 	log.Fatalf("error in unmarshal json: %w", err)
	// }

	// if newsList == nil {
	// 	log.Println("not exist new information")
	// 	return
	// }

	s := NewSender()
	// msg := s.MakeMessage(newsList)

	// if err := s.SendTest(msg); err != nil {
	// 	log.Fatalf("error in send message by line bot: %w", err)
	// }

	titleList, err := c.GetTitleCompetitor("https://baseball.yahoo.co.jp/mlb/stats/")
	if err != nil {
		log.Fatalf("error in get title competitor: %w", err)
	}

	msg := s.MakeTitleCompetitorSummary(titleList)

	log.Println(msg)
	// log.Println(msg)

}

func (s *Sender) SendTest(msgStr string) error {
	bot, err := linebot.New(
		secret,
		token)

	if err != nil {
		return fmt.Errorf("error in new line bot: %w", err)
	}

	msg := linebot.NewTextMessage(msgStr)
	log.Println(msgStr)
	if _, err := bot.BroadcastMessage(msg).Do(); err != nil {
		return fmt.Errorf("error in send message by line bot: %w", err)
	}

	return nil
}
