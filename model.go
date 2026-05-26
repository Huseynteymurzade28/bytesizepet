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
	mgReactionWait
	mgReactionReady
	mgReactionDone
	mgRPSPlaying
	mgRPSDone
	mgMathPlaying
	mgMathDone
)

// ── Sub-models ────────────────────────────────────────────────────────────────

type petModel struct {
	charIdx   int
	name      string
	hunger    int // 0-100; higher = hungrier
	happiness int // 0-100
	energy    int // 0-100
	sleeping  bool
	statusMsg string
	age       int
}

type minigameModel struct {
	phase     mgPhase
	resultMsg string

	target      int
	hint        string
	attempts    int
	maxAttempts int

	signalAt   time.Time
	reactionMs int64

	rpsPetChoice int

	mathExpr       string
	mathOptions    [3]int
	mathCorrectOpt int
}

// ── Root model ────────────────────────────────────────────────────────────────

type model struct {
	screen       screenID
	pet          petModel
	mg           minigameModel
	selectedChar int
	tick         int
	width        int
	height       int
	quitting     bool

	// Animation state
	animFrame   int // 0 or 1 — idle frame toggle
	sleepFrame  int // 0-3 — z / zZ / zZz cycle
	actionAnim  int // 0=none  1=feed  2=play
	actionFrame int // 0-3 — action progress
	envFrame    int // 0 or 1 — habitat toggle
}

// ── Messages ──────────────────────────────────────────────────────────────────

type tickMsg time.Time
type animTickMsg time.Time
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
		statusMsg: ch.Name + " was born! Welcome! " + ch.Emoji,
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

func animTickCmd() tea.Cmd {
	return tea.Tick(700*time.Millisecond, func(t time.Time) tea.Msg {
		return animTickMsg(t)
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
