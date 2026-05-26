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

	case animTickMsg:
		return handleAnimTick(m)

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

const charSelectCols = 3

func updateCharSelect(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	n := len(AllChars)
	switch msg.String() {
	case "left", "h":
		m.selectedChar = (m.selectedChar - 1 + n) % n
	case "right", "l":
		m.selectedChar = (m.selectedChar + 1) % n
	case "up", "k":
		if m.selectedChar >= charSelectCols {
			m.selectedChar -= charSelectCols
		}
	case "down", "j":
		if m.selectedChar+charSelectCols < n {
			m.selectedChar += charSelectCols
		}
	case "enter", " ":
		m.pet = newPet(m.selectedChar)
		m.screen = screenMain
		return m, tea.Batch(tickCmd(), animTickCmd())
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
			p.statusMsg = p.name + " is sleeping, don't disturb them..."
		} else {
			p.hunger = clamp(p.hunger-30, 0, 100)
			p.happiness = clamp(p.happiness+5, 0, 100)
			p.statusMsg = p.name + " nom nom~ ♪"
			m.actionAnim = 1
			m.actionFrame = 0
		}

	case "p": // play
		if p.sleeping {
			p.statusMsg = p.name + " is sleeping, they don't want to play..."
		} else if p.energy < 15 {
			p.statusMsg = p.name + " is too tired! Put them to sleep first [s]"
		} else {
			p.happiness = clamp(p.happiness+20, 0, 100)
			p.energy = clamp(p.energy-15, 0, 100)
			p.statusMsg = p.name + " played with you! ✨"
			m.actionAnim = 2
			m.actionFrame = 0
		}

	case "s": // sleep toggle
		p.sleeping = !p.sleeping
		if p.sleeping {
			p.statusMsg = p.name + " is sleeping... zZz ♪"
		} else {
			p.statusMsg = p.name + " woke up! Good morning! ☀"
		}

	case "1": // guess game
		m.screen = screenMinigame
		m.mg = minigameModel{
			phase:       mgGuessPlaying,
			target:      1 + rand.Intn(9),
			maxAttempts: 5,
			hint:        "I'm thinking of a number between 1 and 9!",
		}

	case "2": // reaction game
		m.screen = screenMinigame
		m.mg = minigameModel{phase: mgReactionWait}
		return m, reactionDelayCmd()

	case "3": // rock paper scissors
		m.screen = screenMinigame
		m.mg = minigameModel{
			phase:        mgRPSPlaying,
			rpsPetChoice: rand.Intn(3),
		}

	case "4": // math quiz
		m.screen = screenMinigame
		m.mg = newMathGame()
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

	case mgRPSPlaying:
		return updateRPS(m, msg)

	case mgMathPlaying:
		return updateMath(m, msg)

	case mgGuessDone, mgReactionDone, mgRPSDone, mgMathDone:
		if msg.String() == "enter" || msg.String() == " " {
			m.screen = screenMain
			m.mg = minigameModel{}
		}

	case mgReactionWait:
		if msg.String() == " " {
			m.mg.phase = mgReactionDone
			m.mg.resultMsg = "Too early! 😅 -5 happiness"
			m.pet.happiness = clamp(m.pet.happiness-5, 0, 100)
		}

	case mgReactionReady:
		if msg.String() == " " {
			elapsed := time.Since(m.mg.signalAt).Milliseconds()
			m.mg.reactionMs = elapsed
			m.mg.phase = mgReactionDone
			switch {
			case elapsed < 300:
				m.mg.resultMsg = fmt.Sprintf("Incredible! %dms — +30 happiness 🚀", elapsed)
				m.pet.happiness = clamp(m.pet.happiness+30, 0, 100)
			case elapsed < 600:
				m.mg.resultMsg = fmt.Sprintf("Amazing! %dms — +20 happiness ⚡", elapsed)
				m.pet.happiness = clamp(m.pet.happiness+20, 0, 100)
			case elapsed < 1000:
				m.mg.resultMsg = fmt.Sprintf("Good! %dms — +10 happiness 👍", elapsed)
				m.pet.happiness = clamp(m.pet.happiness+10, 0, 100)
			default:
				m.mg.resultMsg = fmt.Sprintf("Too slow... %dms. Keep practicing! 🐢", elapsed)
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
		bonus := maxInt(5, 30-m.mg.attempts*4)
		m.pet.happiness = clamp(m.pet.happiness+bonus, 0, 100)
		m.mg.resultMsg = fmt.Sprintf("Correct! The number was %d! +%d happiness 🎉", m.mg.target, bonus)

	case m.mg.attempts >= m.mg.maxAttempts:
		m.mg.phase = mgGuessDone
		m.mg.resultMsg = fmt.Sprintf("Out of tries! The number was %d. 😢", m.mg.target)
		m.pet.happiness = clamp(m.pet.happiness-5, 0, 100)

	case guess < m.mg.target:
		m.mg.hint = fmt.Sprintf("Higher! (%d/%d attempts)", m.mg.attempts, m.mg.maxAttempts)

	default:
		m.mg.hint = fmt.Sprintf("Lower! (%d/%d attempts)", m.mg.attempts, m.mg.maxAttempts)
	}

	return m, nil
}

var rpsNames = [3]string{"Rock", "Paper", "Scissors"}
var rpsEmoji = [3]string{"🪨", "📄", "✂️"}

func updateRPS(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	var player int
	switch msg.String() {
	case "r":
		player = 0
	case "p":
		player = 1
	case "s":
		player = 2
	default:
		return m, nil
	}

	pet := m.mg.rpsPetChoice
	result := (player - pet + 3) % 3

	youLine := fmt.Sprintf("You:   %s %s", rpsNames[player], rpsEmoji[player])
	petLine := fmt.Sprintf("%s: %s %s", m.pet.name, rpsNames[pet], rpsEmoji[pet])

	switch result {
	case 0:
		m.mg.resultMsg = youLine + "\n" + petLine + "\n\n🤝 It's a draw! +10 happiness"
		m.pet.happiness = clamp(m.pet.happiness+10, 0, 100)
	case 1:
		m.mg.resultMsg = youLine + "\n" + petLine + "\n\n🎉 You win! +20 happiness"
		m.pet.happiness = clamp(m.pet.happiness+20, 0, 100)
	case 2:
		m.mg.resultMsg = youLine + "\n" + petLine + "\n\n😅 You lose! -5 happiness"
		m.pet.happiness = clamp(m.pet.happiness-5, 0, 100)
	}

	m.mg.phase = mgRPSDone
	return m, nil
}

func updateMath(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	var chosen int
	switch msg.String() {
	case "1":
		chosen = 0
	case "2":
		chosen = 1
	case "3":
		chosen = 2
	default:
		return m, nil
	}

	if chosen == m.mg.mathCorrectOpt {
		m.mg.resultMsg = "✅ Correct! +25 happiness!"
		m.pet.happiness = clamp(m.pet.happiness+25, 0, 100)
	} else {
		correct := m.mg.mathOptions[m.mg.mathCorrectOpt]
		m.mg.resultMsg = fmt.Sprintf("❌ Wrong! The answer was %d.  -5 happiness", correct)
		m.pet.happiness = clamp(m.pet.happiness-5, 0, 100)
	}

	m.mg.phase = mgMathDone
	return m, nil
}

func newMathGame() minigameModel {
	type cfg struct {
		a, b int
		op   string
	}
	var c cfg

	switch rand.Intn(3) {
	case 0:
		c = cfg{1 + rand.Intn(12), 1 + rand.Intn(12), "+"}
	case 1:
		a := 6 + rand.Intn(10)
		c = cfg{a, 1 + rand.Intn(a-1), "-"}
	case 2:
		c = cfg{2 + rand.Intn(7), 2 + rand.Intn(4), "×"}
	}

	var answer int
	switch c.op {
	case "+":
		answer = c.a + c.b
	case "-":
		answer = c.a - c.b
	case "×":
		answer = c.a * c.b
	}

	off1 := 1 + rand.Intn(4)
	off2 := -(1 + rand.Intn(4))
	if answer+off2 <= 0 {
		off2 = 5 + rand.Intn(4)
	}

	opts := [3]int{answer, answer + off1, answer + off2}
	for i := 2; i > 0; i-- {
		j := rand.Intn(i + 1)
		opts[i], opts[j] = opts[j], opts[i]
	}

	correctIdx := 0
	for i, v := range opts {
		if v == answer {
			correctIdx = i
			break
		}
	}

	return minigameModel{
		phase:          mgMathPlaying,
		mathExpr:       fmt.Sprintf("%d %s %d = ?", c.a, c.op, c.b),
		mathOptions:    opts,
		mathCorrectOpt: correctIdx,
	}
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
			p.statusMsg = p.name + " fully rested, waking up! ☀"
		}
	} else {
		p.hunger = clamp(p.hunger+5, 0, 100)
		p.happiness = clamp(p.happiness-3, 0, 100)
		p.energy = clamp(p.energy-4, 0, 100)
		p.age++

		if m.tick%3 == 0 {
			switch {
			case p.energy <= 0:
				p.statusMsg = p.name + " is exhausted! Press [s] to sleep"
			case p.hunger >= 80:
				p.statusMsg = p.name + " is starving! Press [f] to feed"
			case p.happiness <= 20:
				p.statusMsg = p.name + " is unhappy, play with them! [p]"
			}
		}
	}

	return m, tickCmd()
}

// handleAnimTick advances all animation frame counters.
func handleAnimTick(m model) (model, tea.Cmd) {
	m.animFrame = 1 - m.animFrame
	m.sleepFrame = (m.sleepFrame + 1) % 4
	m.envFrame = 1 - m.envFrame

	if m.actionAnim != 0 {
		m.actionFrame++
		if m.actionFrame >= 4 {
			m.actionAnim = 0
			m.actionFrame = 0
		}
	}

	return m, animTickCmd()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
