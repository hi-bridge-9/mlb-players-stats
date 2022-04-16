package function

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	secret = os.Getenv("LINE_BOT_CHANNEL_SECRET")
	token  = os.Getenv("LINE_BOT_CHANNEL_TOKEN")
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) MakeMessage(ps []*Profile) (msg string) {
	for _, p := range ps {
		msg += fmt.Sprintf("----- %s -----\n", p.PlayersName)
		if p.Pitching.DailyResult != nil {
			msg += makePitchingSummary(p)
		}

		if p.Batting.DailyResult != nil {
			msg += makeBattingSummary(p)
		}
	}
	return msg
}

func makeBattingSummary(p *Profile) (msg string) {
	r := p.Batting.DailyResult
	msg += fmt.Sprintln("【打撃】")

	// 今日の記録
	msg += fmt.Sprintln("<本日>")
	msg += fmt.Sprintf("%s打数 %s安打 %s本塁打\n%s打点 %s盗塁 %s三振 %s四球\n",
		r["AB"], r["H"], r["HR"], r["RBI"], r["SB"], r["SO"], r["BB"])

	// シーズンの記録
	sum := p.Batting.YearSummary
	log.Println(p.Batting.YearSummary)
	msg += fmt.Sprintln("<年間>")
	msg += fmt.Sprintf("打率%s %s安打 %s本塁打\n%s打点 %s盗塁 %s三振\n",
		sum["G"], sum["H"], sum["HR"], sum["RBI"], sum["SB"], sum["SO"])

	return
}

func makePitchingSummary(p *Profile) (msg string) {
	r := p.Pitching.DailyResult
	msg += fmt.Sprintln("【投球】")

	// 今日の記録
	msg += fmt.Sprintln("<本日>")
	msg += fmt.Sprintf("%s回 %s失点 %s奪三振\n%s四球 %s被安打\n",
		r["IP"], r["ER"], r["SO"], r["BB"], r["H"])

	// 勝敗
	if r["W"] == "1" {
		msg += fmt.Sprintln("  --> 勝利投手")
	} else if r["L"] == "1" {
		msg += fmt.Sprintln("  --> 敗戦投手")
	} else {
		msg += fmt.Sprintln("  --> 勝ち負けつかず")
	}

	// シーズンの記録
	sum := p.Pitching.YearSummary
	wl := strings.Split(sum["W-L"], "-")
	msg += fmt.Sprintln("<年間>")
	msg += fmt.Sprintf("%s試合 %s勝%s敗 防御率%s\n",
		sum["G"], wl[0], wl[1], sum["ERA"])

	return
}

func (s *Sender) Send(msgStr string) error {
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


