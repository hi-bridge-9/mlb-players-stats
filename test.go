package function

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func FunctionTest(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadFile("../data/players_list.json")
	if err != nil {
		log.Fatalf("downloadFileIntoMemory: %w", err)
	}

	var ps []*Player
	if err := json.Unmarshal(bytes, &ps); err != nil {
		log.Fatalf("error in unmarshal json: %w", err)
	}

	var c *Crawler
	newsList, err := c.GetUpdateInfo(ps)
	if err != nil {
		log.Fatalf("error in unmarshal json: %w", err)
	}

	if newsList == nil {
		log.Println("not exist new information")
		return
	}

	s := NewSender()
	msg := s.MakeMessage(newsList)

	// if err := s.SendTest(msg); err != nil {
	// 	log.Fatalf("error in send message by line bot: %w", err)
	// }
	log.Println(msg)

}

// func (s *Sender) SendTest(msgStr string) error {
// 	bot, err := linebot.New(
// 		secret,
// 		token)

// 	if err != nil {
// 		return fmt.Errorf("error in new line bot: %w", err)
// 	}

// 	msg := linebot.NewTextMessage(msgStr)
// 	log.Println(msgStr)
// 	if _, err := bot.BroadcastMessage(msg).Do(); err != nil {
// 		return fmt.Errorf("error in send message by line bot: %w", err)
// 	}

// 	return nil
// }
