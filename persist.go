package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// saveData is the on-disk snapshot of a pet. It lives in the user's config dir
// so the creature genuinely persists between runs and ages while you're away.
type saveData struct {
	CharIdx   int       `json:"char_idx"`
	Name      string    `json:"name"`
	Hunger    int       `json:"hunger"`
	Happiness int       `json:"happiness"`
	Energy    int       `json:"energy"`
	Sleeping  bool      `json:"sleeping"`
	Age       int       `json:"age"`
	LastSeen  time.Time `json:"last_seen"`
}

func savePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tamagotchi", "save.json"), nil
}

// savePet writes the pet to disk. Errors are returned but callers treat saving
// as best-effort.
func savePet(p petModel) error {
	path, err := savePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data := saveData{
		CharIdx: p.charIdx, Name: p.name,
		Hunger: p.hunger, Happiness: p.happiness, Energy: p.energy,
		Sleeping: p.sleeping, Age: p.age, LastSeen: time.Now(),
	}
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, buf, 0o644)
}

// loadPet restores a saved pet and applies offline progression. The second
// return is false when there's no valid save to load.
func loadPet() (petModel, string, bool) {
	path, err := savePath()
	if err != nil {
		return petModel{}, "", false
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		return petModel{}, "", false
	}
	var d saveData
	if err := json.Unmarshal(buf, &d); err != nil || d.CharIdx < 0 || d.CharIdx >= len(AllChars) {
		return petModel{}, "", false
	}

	p := petModel{
		charIdx: d.CharIdx, name: d.Name,
		hunger: d.Hunger, happiness: d.Happiness, energy: d.Energy,
		sleeping: d.Sleeping, age: d.Age,
	}

	// Offline progression: gentle, 1 tick per 30 real seconds, capped so a long
	// absence stings but rarely kills outright.
	away := time.Since(d.LastSeen)
	ticks := int(away.Seconds() / 30)
	if ticks > 240 {
		ticks = 240
	}
	msg := p.name + " missed you~ welcome back! 🐾"
	if ticks > 0 {
		if p.sleeping {
			p.energy = clamp(p.energy+ticks*4, 0, 100)
			p.hunger = clamp(p.hunger+ticks*2, 0, 100)
			if p.energy >= 100 {
				p.sleeping = false
			}
		} else {
			p.hunger = clamp(p.hunger+ticks*2, 0, 100)
			p.happiness = clamp(p.happiness-ticks, 0, 100)
			p.energy = clamp(p.energy-ticks, 0, 100)
		}
		p.age += ticks
		msg = p.name + " waited " + humanizeDur(away) + " for you 🐾"
	}
	p.statusMsg = msg
	return p, msg, true
}

func humanizeDur(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "a moment"
	case d < time.Hour:
		return itoa(int(d.Minutes())) + "m"
	case d < 24*time.Hour:
		return itoa(int(d.Hours())) + "h"
	default:
		return itoa(int(d.Hours()/24)) + "d"
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
