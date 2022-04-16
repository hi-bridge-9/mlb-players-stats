package function

type Player struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	IsPitcher bool   `json:"pitching,omitempty"`
	IsBatter  bool   `json:"batting,omitempty"`
}

type Profile struct {
	PlayersName string
	Pitching    Stats
	Batting     Stats
}

type Stats struct {
	YearSummary map[string]string
	DailyResult map[string]string
	Date        string
}
