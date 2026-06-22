package main

import "testing"

// TestRenderScreens is a smoke test: it drives every screen through a few
// animation frames and asserts nothing panics and output is produced.
func TestRenderScreens(t *testing.T) {
	m := initialModel()
	m.loaded = false
	m.width, m.height = 90, 40

	for _, sc := range []screenID{screenCharSelect, screenMain, screenMinigame} {
		m.screen = sc
		if sc != screenCharSelect {
			m.pet = newPet(0)
		}
		m.weather = weatherState{ok: true, kind: weatherRain, desc: "Light rain", tempC: 14, area: "Istanbul"}
		for i := 0; i < 12; i++ {
			m, _ = handleAnimTick(m)
		}
		if out := m.View(); len(out) == 0 {
			t.Fatalf("empty render for screen %d", sc)
		}
	}
}

// TestEveryCharRenders ensures all sprites/moods compose without index panics.
func TestEveryCharRenders(t *testing.T) {
	moods := []petMood{moodNormal, moodBlink, moodHappy, moodHungry, moodTired, moodSleep, moodDead}
	for i := range AllChars {
		base := AllChars[i].spriteLines(moodNormal)
		for _, mood := range moods {
			lines := AllChars[i].spriteLines(mood)
			if len(lines) != len(base) {
				t.Fatalf("char %d mood %d: %d lines, want %d", i, mood, len(lines), len(base))
			}
		}
	}
}
