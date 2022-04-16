package function

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var filter = map[string]string{
	"table":        "div.responsive-datatable__scrollable > div > table",
	"year-summary": `div.player-stats-yearbyyear > div.yearbyyear[data-summary-view="%s"]`,
	"daily-result": `div.player-splits--last.player-splits--last-3.has-xgames[data-split-view="%s"]`,
}

type Crawler struct{}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) GetUpdateInfo(ps []*Player) ([]*Profile, error) {
	var profileList []*Profile
	for _, p := range ps {
		profile, err := c.Fetch(p)
		if err != nil {
			return nil, err
		}
		if p.PlayedThisYear(profile) {
			if p.PlayedToday(profile) {
				// if p.IsTwoWayPlayer(profile) {
				// 	excludeNotLatesStats(profile)
				// }
				profileList = append(profileList, profile)
			}
		}
	}
	return profileList, nil
}

func (p *Player) IsTwoWayPlayer(pro *Profile) bool {
	return pro.Pitching.Date != "" && pro.Batting.Date != ""
}

func (p *Player) PlayedThisYear(pro *Profile) bool {
	return pro.Pitching.Date != "" || pro.Batting.Date != ""
}

func (p *Player) PlayedToday(pro *Profile) bool {
	now := getTimeJST()
	latestGame := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)

	if pro.Pitching.Date != "" {
		pitDate := convertToDate(pro.Pitching.Date, now)
		// 取得対象の試合時刻は日本時刻で考えると前日
		if pitDate.Unix() >= latestGame.Unix() {
			return true
		}
	}

	if pro.Batting.Date != "" {
		batDate := convertToDate(pro.Batting.Date, now)
		// 取得対象の試合時刻は日本時刻で考えると前日
		if batDate.Unix() >= latestGame.Unix() {
			return true
		}
	}
	return false
}

func getTimeJST() time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.Now().In(jst)
}

func convertToDate(str string, now time.Time) time.Time {
	li := strings.Split(str, " ")
	monthAndDay := strings.Split(li[0], "/")
	month, day := monthAndDay[0], monthAndDay[1]

	m, _ := strconv.Atoi(month)
	d, _ := strconv.Atoi(day)

	return time.Date(now.Year(), time.Month(m), d, 0, 0, 0, 0, time.Local)
}

func (c *Crawler) Fetch(p *Player) (*Profile, error) {
	doc, err := goquery.NewDocument(p.URL)
	if err != nil {
		return nil, fmt.Errorf("error in get site resource: %w", err)
	}

	return pickProfile(doc.Find("sectiopro.section-container"), p), nil
}

func pickProfile(data *goquery.Selection, p *Player) *Profile {
	var pro = &Profile{PlayersName: p.Name}
	if p.IsPitcher {
		pro.Pitching.YearSummary = extractSummary(data.Find(fmt.Sprintf(filter["year-summary"], "pitching")))
		pro.Pitching.Date, pro.Pitching.DailyResult = extractDailyResult(data.Find(fmt.Sprintf(filter["daily-result"], "pitching")))
	}

	if p.IsBatter {
		pro.Batting.YearSummary = extractSummary(data.Find(fmt.Sprintf(filter["year-summary"], "hitting")))
		pro.Batting.Date, pro.Batting.DailyResult = extractDailyResult(data.Find(fmt.Sprintf(filter["daily-result"], "hitting")))
	}

	return pro
}

func excludeNotLatesStats(pro *Profile) {
	now := getTimeJST()
	latestGame := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)

	if pro.Pitching.Date != "" {
		pitDate := convertToDate(pro.Pitching.Date, now)
		if pitDate.Unix() >= latestGame.Unix() {
			pro.Pitching = Stats{}
		}
	}

	if pro.Batting.Date != "" {
		batDate := convertToDate(pro.Batting.Date, now)
		if batDate.Unix() >= latestGame.Unix() {
			pro.Batting = Stats{}
		}
	}
}

func extractSummary(data *goquery.Selection) map[string]string {

	table := data.Find(filter["table"])
	col := table.Find("thead > tr")
	row := table.Find("tbody > tr")

	topicCount := 0
	col.Find("th > span").Each(func(_ int, v *goquery.Selection) {
		topicCount += 1
	})

	var summary = map[string]string{}
	for i := 0; i < topicCount; i++ {
		topic := col.Find(fmt.Sprintf("th.no-sort.col-%d > span", i)).Text()
		value := row.Find(fmt.Sprintf("td.col-%d.row-0 > span", i)).Text()
		summary[topic] = value
	}
	log.Println(summary)

	return summary
}

func extractDailyResult(data *goquery.Selection) (string, map[string]string) {
	table := data.Find(filter["table"])
	col := table.Find("thead > tr")
	row := table.Find("tbody > tr")

	var topicCount int
	col.Find("th > span").Each(func(_ int, v *goquery.Selection) {
		topicCount += 1
	})

	var result = map[string]string{}

	// 最初の一列目は試合の日程
	date := row.Find("td.col-0.row-0.td--text > span").Text()
	for i := 1; i < topicCount; i++ {
		topic := col.Find(fmt.Sprintf("th.no-sort.col-%d > span", i)).Text()
		value := row.Find(fmt.Sprintf("td.col-%d.row-0 > span", i)).Text()
		result[topic] = value
	}

	log.Println(result)
	return date, result
}
