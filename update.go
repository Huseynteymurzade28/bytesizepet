package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.screen {
		case screenCharSelect:
			return updateCharSelect(m, msg)
		case screenMain:
			return updateMain(m, msg)
		case screenMinigame:
			return updateMinigame(m, msg)
		}

	case tickMsg:
		if m.screen != screenCharSelect {
			return handleTick(m)
		}
		return m, tickCmd()

	case reactionSignalMsg:
		if m.mg.phase == mgReactionWait {
			m.mg.phase = mgReactionReady
			m.mg.signalAt = time.Now()
		}
		return m, nil
	}

	return m, nil
}

// ── Character select ──────────────────────────────────────────────────────────

func updateCharSelect(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		m.selectedChar = (m.selectedChar - 1 + len(AllChars)) % len(AllChars)
	case "right", "l":
		m.selectedChar = (m.selectedChar + 1) % len(AllChars)
	case "enter", " ":
		m.pet = newPet(m.selectedChar)
		m.screen = screenMain
		return m, tickCmd()
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

// ── Main screen ───────────────────────────────────────────────────────────────

func updateMain(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	p := &m.pet

	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "f": // feed
		if p.sleeping {
			p.statusMsg = p.name + " uyuyor, onu rahatsız etme..."
		} else {
			p.hunger = clamp(p.hunger-30, 0, 100)
			p.happiness = clamp(p.happiness+5, 0, 100)
			p.statusMsg = p.name + " beslendi! Nom nom~ ♪"
		}

	case "p": // play
		if p.sleeping {
			p.statusMsg = p.name + " uyuyor, oynamak istemez..."
		} else if p.energy < 15 {
			p.statusMsg = p.name + " çok yorgun! Önce uyut (s)"
		} else {
			p.happiness = clamp(p.happiness+20, 0, 100)
			p.energy = clamp(p.energy-15, 0, 100)
			p.statusMsg = p.name + " seninle oynadı! ✨"
		}

	case "s": // sleep toggle
		p.sleeping = !p.sleeping
		if p.sleeping {
			p.statusMsg = p.name + " uyuyor... zZz ♪"
		} else {
			p.statusMsg = p.name + " uyandı! Günaydın! ☀"
		}

	case "1": // start guess game
		m.screen = screenMinigame
		m.mg = minigameModel{
			phase:       mgGuessPlaying,
			target:      1 + rand.Intn(9),
			maxAttempts: 5,
			hint:        "1-9 arasında bir sayı düşündüm!",
		}

	case "2": // start reaction game
		m.screen = screenMinigame
		m.mg = minigameModel{phase: mgReactionWait}
		return m, reactionDelayCmd()
	}

	return m, nil
}

// ── Minigame screen ───────────────────────────────────────────────────────────

func updateMinigame(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.screen = screenMain
		m.mg = minigameModel{}
		return m, nil
	}

	switch m.mg.phase {
	case mgGuessPlaying:
		return updateGuess(m, msg)

	case mgGuessDone, mgReactionDone:
		if msg.String() == "enter" || msg.String() == " " {
			m.screen = screenMain
			m.mg = minigameModel{}
		}

	case mgReactionWait:
		if msg.String() == " " {
			m.mg.phase = mgReactionDone
			m.mg.resultMsg = "Çok erken bastın! 😅"
			m.pet.happiness = clamp(m.pet.happiness-5, 0, 100)
		}

	case mgReactionReady:
		if msg.String() == " " {
			elapsed := time.Since(m.mg.signalAt).Milliseconds()
			m.mg.reactionMs = elapsed
			m.mg.phase = mgReactionDone
			switch {
			case elapsed < 300:
				m.mg.resultMsg = fmt.Sprintf("İnanılmaz! %dms — +30 mutluluk 🚀", elapsed)
				m.pet.happiness = clamp(m.pet.happiness+30, 0, 100)
			case elapsed < 600:
				m.mg.resultMsg = fmt.Sprintf("Harika! %dms — +20 mutluluk ⚡", elapsed)
				m.pet.happiness = clamp(m.pet.happiness+20, 0, 100)
			case elapsed < 1000:
				m.mg.resultMsg = fmt.Sprintf("İyi! %dms — +10 mutluluk 👍", elapsed)
				m.pet.happiness = clamp(m.pet.happiness+10, 0, 100)
			default:
				m.mg.resultMsg = fmt.Sprintf("Yavaş... %dms. Pratik yap! 🐢", elapsed)
			}
		}
	}

	return m, nil
}

func updateGuess(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	k := msg.String()
	if len(k) != 1 || k[0] < '1' || k[0] > '9' {
		return m, nil
	}

	guess := int(k[0] - '0')
	m.mg.attempts++

	switch {
	case guess == m.mg.target:
		m.mg.phase = mgGuessDone
		bonus := max(5, 30-m.mg.attempts*4)
		m.pet.happiness = clamp(m.pet.happiness+bonus, 0, 100)
		m.mg.resultMsg = fmt.Sprintf("Doğru! Sayı %d'ydi! +%d mutluluk 🎉", m.mg.target, bonus)

	case m.mg.attempts >= m.mg.maxAttempts:
		m.mg.phase = mgGuessDone
		m.mg.resultMsg = fmt.Sprintf("Bitti! Sayı %d'ydi. 😢", m.mg.target)
		m.pet.happiness = clamp(m.pet.happiness-5, 0, 100)

	case guess < m.mg.target:
		m.mg.hint = fmt.Sprintf("Daha büyük! (%d/%d deneme)", m.mg.attempts, m.mg.maxAttempts)

	default:
		m.mg.hint = fmt.Sprintf("Daha küçük! (%d/%d deneme)", m.mg.attempts, m.mg.maxAttempts)
	}

	return m, nil
}

// ── Tick logic ────────────────────────────────────────────────────────────────

func handleTick(m model) (model, tea.Cmd) {
	m.tick++
	p := &m.pet

	if p.sleeping {
		p.energy = clamp(p.energy+8, 0, 100)
		p.hunger = clamp(p.hunger+3, 0, 100)
		if p.energy >= 100 {
			p.sleeping = false
			p.statusMsg = p.name + " tam dinlendi, uyandı! ☀"
		}
	} else {
		p.hunger = clamp(p.hunger+5, 0, 100)
		p.happiness = clamp(p.happiness-3, 0, 100)
		p.energy = clamp(p.energy-4, 0, 100)
		p.age++

		if m.tick%3 == 0 {
			switch {
			case p.energy <= 0:
				p.statusMsg = p.name + " çok yorgun! (s) ile uyut"
			case p.hunger >= 80:
				p.statusMsg = p.name + " çok acıktı! (f) ile besle"
			case p.happiness <= 20:
				p.statusMsg = p.name + " mutsuz, oyna! (p)"
			}
		}
	}

	return m, tickCmd()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
