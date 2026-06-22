package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// weatherKind is the simplified bucket we drive visuals from.
type weatherKind int

const (
	weatherUnknown weatherKind = iota
	weatherClear
	weatherClouds
	weatherFog
	weatherRain
	weatherSnow
	weatherThunder
)

// weatherState is what the UI keeps around between refreshes.
type weatherState struct {
	kind    weatherKind
	desc    string
	tempC   int
	area    string
	ok      bool
	fetched time.Time
}

func (k weatherKind) icon() string {
	switch k {
	case weatherClear:
		return "☀️"
	case weatherClouds:
		return "☁️"
	case weatherFog:
		return "🌫️"
	case weatherRain:
		return "🌧️"
	case weatherSnow:
		return "❄️"
	case weatherThunder:
		return "⛈️"
	default:
		return "❔"
	}
}

type weatherMsg struct{ state weatherState }
type weatherRefreshMsg time.Time

// fetchWeatherCmd hits wttr.in (no API key required, auto-locates by IP) on a
// background goroutine and never blocks the UI. Failure is silent: the game
// just keeps its previous weather.
func fetchWeatherCmd() tea.Cmd {
	return func() tea.Msg {
		city := os.Getenv("TAMAGOTCHI_CITY") // optional override
		url := "https://wttr.in/" + city + "?format=j1"

		client := &http.Client{Timeout: 6 * time.Second}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return weatherMsg{weatherState{}}
		}
		req.Header.Set("User-Agent", "tamagotchi-tui")
		resp, err := client.Do(req)
		if err != nil {
			return weatherMsg{weatherState{}}
		}
		defer resp.Body.Close()

		var payload struct {
			CurrentCondition []struct {
				TempC       string `json:"temp_C"`
				WeatherCode string `json:"weatherCode"`
				WeatherDesc []struct {
					Value string `json:"value"`
				} `json:"weatherDesc"`
			} `json:"current_condition"`
			NearestArea []struct {
				AreaName []struct {
					Value string `json:"value"`
				} `json:"areaName"`
			} `json:"nearest_area"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil ||
			len(payload.CurrentCondition) == 0 {
			return weatherMsg{weatherState{}}
		}

		cc := payload.CurrentCondition[0]
		st := weatherState{ok: true, fetched: time.Now()}
		st.kind = codeToKind(cc.WeatherCode)
		if len(cc.WeatherDesc) > 0 {
			st.desc = cc.WeatherDesc[0].Value
		}
		if t, err := jatoi(cc.TempC); err == nil {
			st.tempC = t
		}
		if len(payload.NearestArea) > 0 && len(payload.NearestArea[0].AreaName) > 0 {
			st.area = payload.NearestArea[0].AreaName[0].Value
		}
		return weatherMsg{st}
	}
}

// weatherRefreshCmd schedules the next refresh ~15 minutes out.
func weatherRefreshCmd() tea.Cmd {
	return tea.Tick(15*time.Minute, func(t time.Time) tea.Msg {
		return weatherRefreshMsg(t)
	})
}

// codeToKind maps WWO weather codes (as returned by wttr.in) to our buckets.
func codeToKind(code string) weatherKind {
	switch code {
	case "113":
		return weatherClear
	case "116", "119", "122":
		return weatherClouds
	case "143", "248", "260":
		return weatherFog
	case "200", "386", "389", "392", "395":
		return weatherThunder
	case "179", "227", "230", "323", "326", "329", "332", "335", "338", "368", "371", "374", "377":
		return weatherSnow
	case "176", "182", "185", "263", "266", "281", "284", "293", "296", "299",
		"302", "305", "308", "311", "314", "317", "320", "350", "353", "356",
		"359", "362", "365":
		return weatherRain
	default:
		return weatherClouds
	}
}

// jatoi is a tiny atoi that tolerates leading '+'/'-' and ignores junk.
func jatoi(s string) (int, error) {
	neg := false
	i := 0
	if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
		neg = s[0] == '-'
		i = 1
	}
	n := 0
	got := false
	for ; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			break
		}
		n = n*10 + int(s[i]-'0')
		got = true
	}
	if !got {
		return 0, errNoNumber
	}
	if neg {
		n = -n
	}
	return n, nil
}

var errNoNumber = &weatherErr{"no number"}

type weatherErr struct{ s string }

func (e *weatherErr) Error() string { return e.s }
