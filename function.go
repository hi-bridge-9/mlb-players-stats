package function

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	bucket    = os.Getenv("GCS_BUCKET_NAME")
	object    = os.Getenv("GCS_OBJECT_NAME")
	titlePage = os.Getenv("MLB_TITLE_COMPETITON_PAGE_URL")
)

func Function(w http.ResponseWriter, r *http.Request) {
	bytes, err := downloadFileIntoMemory(ioutil.Discard, bucket, object)
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

	if err := s.Send(msg); err != nil {
		log.Fatalf("error in send message by line bot: %w", err)
	}

	if time.Now().Weekday() == time.Sunday {
		titleList, err := c.GetTitleCompetitor(titlePage)
		if err != nil {
			log.Fatalf("error in get title competitor: %w", err)
		}

		msg := s.MakeTitleCompetitorSummary(titleList)
		if err := s.Send(msg); err != nil {
			log.Fatalf("error in send message by line bot: %w", err)
		}
	}

	log.Println("Success")

}
