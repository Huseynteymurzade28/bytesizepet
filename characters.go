package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// init normalises every sprite into ONE shared box size (width × height),
// centred horizontally and bottom-aligned vertically. This makes the
// character-select grid perfectly uniform and keeps every pet's feet on the
// same ground line in the diorama, while expressions never jump sideways.
func init() {
	// 1. Global box dimensions across all characters (base + overrides).
	gw, gh := 0, 0
	for i := range AllChars {
		c := &AllChars[i]
		if len(c.base) > gh {
			gh = len(c.base)
		}
		for _, l := range c.base {
			if w := lipgloss.Width(l); w > gw {
				gw = w
			}
		}
		for _, ov := range c.moods {
			for _, l := range ov {
				if w := lipgloss.Width(l); w > gw {
					gw = w
				}
			}
		}
	}

	// 2. Pad each character to that box.
	for i := range AllChars {
		c := &AllChars[i]
		for j, l := range c.base {
			c.base[j] = padCenter(l, gw)
		}
		for _, ov := range c.moods {
			for r, l := range ov {
				ov[r] = padCenter(l, gw)
			}
		}
		// Bottom-align by prepending blank rows; shift override indices to match.
		if pre := gh - len(c.base); pre > 0 {
			blank := strings.Repeat(" ", gw)
			padded := make([]string, 0, gh)
			for k := 0; k < pre; k++ {
				padded = append(padded, blank)
			}
			c.base = append(padded, c.base...)

			shifted := make(map[petMood]map[int]string, len(c.moods))
			for mood, ov := range c.moods {
				no := make(map[int]string, len(ov))
				for r, l := range ov {
					no[r+pre] = l
				}
				shifted[mood] = no
			}
			c.moods = shifted
		}
	}
}

func padCenter(s string, w int) string {
	gap := w - lipgloss.Width(s)
	if gap <= 0 {
		return s
	}
	left := gap / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", gap-left)
}

// ── Moods ─────────────────────────────────────────────────────────────────────

type petMood int

const (
	moodNormal petMood = iota
	moodBlink
	moodHappy
	moodHungry
	moodTired
	moodSleep
	moodDead
)

// CharDef is one creature. Every character is hand-drawn with its own unique
// silhouette in `base`; `moods` only swaps the line(s) that carry the eyes /
// mouth so each expression stays cheap to define while the body stays distinct.
type CharDef struct {
	Name    string
	Species string
	Emoji   string

	Body  lipgloss.Color
	Belly lipgloss.Color

	Happy  lipgloss.Color
	Hungry lipgloss.Color
	Tired  lipgloss.Color
	Sleep  lipgloss.Color

	base  []string
	moods map[petMood]map[int]string // mood → {rowIndex: replacementLine}
}

func (c CharDef) spriteLines(mood petMood) []string {
	out := make([]string, len(c.base))
	copy(out, c.base)
	if ov, ok := c.moods[mood]; ok {
		for row, txt := range ov {
			if row >= 0 && row < len(out) {
				out[row] = txt
			}
		}
	}
	return out
}

func (c CharDef) moodColor(mood petMood) lipgloss.Color {
	switch mood {
	case moodHappy:
		return c.Happy
	case moodHungry:
		return c.Hungry
	case moodTired:
		return c.Tired
	case moodSleep:
		return c.Sleep
	case moodDead:
		return lipgloss.Color("#7c7c96")
	default:
		return c.Body
	}
}

// AllChars is the roster. Each sprite is hand-drawn with its own silhouette,
// ears, facial features, body and feet — kept left-right symmetric so the
// shared centring keeps every vertical lined up.
var AllChars = []CharDef{
	// ── Mochi — a cat: tufted ears, shiny eyes, whiskers, paws & tail ───────
	{
		Name: "Mochi", Species: "Cat", Emoji: "🐱",
		Body: "#f3b8d4", Belly: "#fff3fa",
		Happy: "#f7d6a6", Hungry: "#e88a70", Tired: "#bdb0e6", Sleep: "#b6a8e2",
		base: []string{
			"   /\\___/\\   ",
			"  /  . .  \\  ",
			" ( ( ◕ ◕ ) ) ",
			" =\\   ▾   /= ",
			"   ) \\_/ (   ",
			"  /`-----'\\  ",
			" / /     \\ \\ ",
			"( (_)   (_) )",
			"  \\_______/  ",
		},
		moods: map[petMood]map[int]string{
			moodBlink:  {2: " ( ( - - ) ) "},
			moodHappy:  {2: " ( ( ^ ^ ) ) ", 4: "   )\\___/(   "},
			moodHungry: {2: " ( ( O O ) ) ", 4: "   ) ~~~ (   "},
			moodTired:  {2: " ( ( u u ) ) "},
			moodSleep:  {2: " ( ( - - ) ) ", 4: "   ) zzz (   "},
			moodDead:   {2: " ( ( x x ) ) ", 4: "   ) ... (   "},
		},
	},
	// ── Biscuit — a dog: floppy ears, brows, snout, tongue, paws ────────────
	{
		Name: "Biscuit", Species: "Dog", Emoji: "🐶",
		Body: "#e0b487", Belly: "#fff1dc",
		Happy: "#a8df7e", Hungry: "#e87858", Tired: "#aebfe0", Sleep: "#a8c0d8",
		base: []string{
			" ___     ___ ",
			"/   \\   /   \\",
			"|    \\_/    |",
			"|  ●     ●  |",
			"|     ▾     |",
			" \\  (___)  / ",
			"  \\       /  ",
			"  /|     |\\  ",
			" (_U_____U_) ",
		},
		moods: map[petMood]map[int]string{
			moodBlink:  {3: "|  ―     ―  |"},
			moodHappy:  {3: "|  ^     ^  |", 5: " \\  \\‿‿/  /  "},
			moodHungry: {3: "|  O     O  |", 5: " \\  ~~~~~  / "},
			moodTired:  {3: "|  u     u  |"},
			moodSleep:  {3: "|  -     -  |", 5: " \\  (zzz)  / "},
			moodDead:   {3: "|  x     x  |", 5: " \\  (___)  / "},
		},
	},
	// ── Piko — a rabbit: tall ears, whiskers, buck teeth, big feet ──────────
	{
		Name: "Piko", Species: "Rabbit", Emoji: "🐰",
		Body: "#f7d3ec", Belly: "#fffafd",
		Happy: "#cdf3cd", Hungry: "#f0a0a0", Tired: "#c5bdee", Sleep: "#c0b8e8",
		base: []string{
			"  /\\   /\\  ",
			" |  | |  | ",
			" |  |_|  | ",
			"  \\     /  ",
			" ( ◕ . ◕ ) ",
			"=(   ▾   )=",
			" ( \\___/ ) ",
			" _/|   |\\_ ",
			"(__)   (__)",
		},
		moods: map[petMood]map[int]string{
			moodBlink:  {4: " ( - . - ) "},
			moodHappy:  {4: " ( ^ . ^ ) ", 6: " ( \\^_^/ ) "},
			moodHungry: {4: " ( O . O ) ", 6: " ( ~~~~~ ) "},
			moodTired:  {4: " ( u . u ) "},
			moodSleep:  {4: " ( - . - ) ", 6: " (  zzz  ) "},
			moodDead:   {4: " ( x . x ) ", 6: " (  ...  ) "},
		},
	},
	// ── Boo — a ghost: domed head, floating arms, triple-wave tail ──────────
	{
		Name: "Boo", Species: "Ghost", Emoji: "👻",
		Body: "#e2dbf6", Belly: "#ffffff",
		Happy: "#ffffff", Hungry: "#e87878", Tired: "#bcb4e6", Sleep: "#b0a8e0",
		base: []string{
			"   .---.   ",
			"  /     \\  ",
			" / ◕   ◕ \\ ",
			" |   ▾   | ",
			" |  \\_/  | ",
			" |       | ",
			"(  )   (  )",
			" \\/ \\_/ \\/ ",
			"  ~v~v~v~  ",
		},
		moods: map[petMood]map[int]string{
			moodBlink:  {2: " / -   - \\ "},
			moodHappy:  {2: " / ^   ^ \\ ", 4: " |  \\_/  | "},
			moodHungry: {2: " / O   O \\ ", 4: " |  ~~~  | "},
			moodTired:  {2: " / u   u \\ "},
			moodSleep:  {2: " / -   - \\ ", 4: " |  zzz  | "},
			moodDead:   {2: " / x   x \\ ", 4: " |  ...  | "},
		},
	},
	// ── Kuma — a bear: round ears, broad body, muzzle, big paws ─────────────
	{
		Name: "Kuma", Species: "Bear", Emoji: "🐻",
		Body: "#d2ae86", Belly: "#f6e2c0",
		Happy: "#f3c878", Hungry: "#e07050", Tired: "#bcb0da", Sleep: "#b0a8d0",
		base: []string{
			" (o)     (o) ",
			"  \\_______/  ",
			" /         \\ ",
			"|  ●     ●  |",
			"|     ▾     |",
			" \\  (___)  / ",
			"  \\       /  ",
			"  /|     |\\  ",
			" (_U_____U_) ",
		},
		moods: map[petMood]map[int]string{
			moodBlink:  {3: "|  ―     ―  |"},
			moodHappy:  {3: "|  ^     ^  |", 5: " \\  \\___/  / "},
			moodHungry: {3: "|  O     O  |", 5: " \\  ~~~~~  / "},
			moodTired:  {3: "|  u     u  |"},
			moodSleep:  {3: "|  -     -  |", 5: " \\  (zzz)  / "},
			moodDead:   {3: "|  x     x  |", 5: " \\  (___)  / "},
		},
	},
	// ── Hoshi — a hamster: tiny ears, huge cheeks, paws, a seed, feet ───────
	{
		Name: "Hoshi", Species: "Hamster", Emoji: "🐹",
		Body: "#f6d2a8", Belly: "#fff2dd",
		Happy: "#f8e87a", Hungry: "#e87860", Tired: "#c2bae6", Sleep: "#b8b0e0",
		base: []string{
			"   ^^ ^^   ",
			"  /─────\\  ",
			" ( ◕   ◕ ) ",
			"((   ▾   ))",
			" (  \\_/  ) ",
			" ( °     ) ",
			"  \\__|__/  ",
			"   U   U   ",
		},
		moods: map[petMood]map[int]string{
			moodBlink:  {2: " ( -   - ) "},
			moodHappy:  {2: " ( ^   ^ ) ", 4: " (  \\v/  ) "},
			moodHungry: {2: " ( O   O ) ", 4: " ( ~~~~~ ) "},
			moodTired:  {2: " ( u   u ) "},
			moodSleep:  {2: " ( -   - ) ", 4: " (  zzz  ) "},
			moodDead:   {2: " ( x   x ) ", 4: " (  ...  ) "},
		},
	},
}
