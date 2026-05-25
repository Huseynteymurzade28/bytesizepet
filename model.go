package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ── Screens ───────────────────────────────────────────────────────────────────

type screenID int

const (
	screenCharSelect screenID = iota
	screenMain
	screenMinigame
)

// ── Minigame phases ───────────────────────────────────────────────────────────

type mgPhase int

const (
	mgNone mgPhase = iota
	mgGuessPlaying
	mgGuessDone
	mgReactionWait  // waiting for the signal to fire
	mgReactionReady // signal is visible, waiting for keypress
	mgReactionDone
)

// ── Sub-models ────────────────────────────────────────────────────────────────

type petModel struct {
	charIdx   int
	name      string
	hunger    int // 0-100; higher means hungrier
	happiness int // 0-100
	energy    int // 0-100
	sleeping  bool
	statusMsg string
	age       int // number of ticks lived
}

type minigameModel struct {
	phase       mgPhase
	target      int    // guess game: secret number
	hint        string // guess game: feedback line
	attempts    int
	maxAttempts int
	signalAt    time.Time // reaction game: when signal fired
	reactionMs  int64
	resultMsg   string
}

// ── Root model ────────────────────────────────────────────────────────────────

type model struct {
	screen       screenID
	pet          petModel
	mg           minigameModel
	selectedChar int // cursor on char-select screen
	tick         int
	width        int
	height       int
	quitting     bool
}

// ── Messages ──────────────────────────────────────────────────────────────────

type tickMsg time.Time
type reactionSignalMsg struct{}

// ── Constructors ──────────────────────────────────────────────────────────────

func newPet(charIdx int) petModel {
	ch := AllChars[charIdx]
	return petModel{
		charIdx:   charIdx,
		name:      ch.Name,
		hunger:    20,
		happiness: 80,
		energy:    90,
		statusMsg: ch.Name + " dünyaya geldi! " + ch.Emoji,
	}
}

func initialModel() model {
	return model{
		screen:       screenCharSelect,
		selectedChar: 0,
		width:        80,
		height:       24,
	}
}

// ── Commands ──────────────────────────────────────────────────────────────────

func tickCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func reactionDelayCmd() tea.Cmd {
	delay := time.Duration(1500+rand.Intn(3000)) * time.Millisecond
	return tea.Tick(delay, func(time.Time) tea.Msg {
		return reactionSignalMsg{}
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
