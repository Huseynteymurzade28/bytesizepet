package main

import "github.com/charmbracelet/lipgloss"

// CharArt holds all emotional states for a character. Every field must have
// the same number of lines so the UI stays stable.
type CharArt struct {
	Normal, Happy, Hungry, Sleep, Tired, Dead string
}

type CharColors struct {
	Normal, Happy, Hungry, Sleep, Tired lipgloss.Color
}

type CharDef struct {
	Name    string
	Species string
	Emoji   string
	Colors  CharColors
	Art     CharArt
}

// AllChars is the selectable roster shown on the character-select screen.
var AllChars = []CharDef{
	{
		Name:    "Mochi",
		Species: "Kedi",
		Emoji:   "🐱",
		Colors:  CharColors{"226", "82", "196", "147", "208"},
		Art: CharArt{
			Normal: " /\\ /\\\n(=^.^=)\n (\")(\")",
			Happy:  " /\\ /\\\n(=^ω^=)\n (\")(\")",
			Hungry: " /\\ /\\\n(=>.<=)\n (\")(\")",
			Sleep:  " /\\ /\\\n(=-.-=)\n (\")(\")",
			Tired:  " /\\ /\\\n(=~.~=)\n (\")(\")",
			Dead:   " /\\ /\\\n(=x.x=)\n (\")(\")",
		},
	},
	{
		Name:    "Biscuit",
		Species: "Köpek",
		Emoji:   "🐶",
		Colors:  CharColors{"178", "82", "196", "147", "208"},
		Art: CharArt{
			Normal: " n   n\n(ᵔ.ᵔ)\n /W\\ ",
			Happy:  " n   n\n(ᵔωᵔ)\n /W\\ ",
			Hungry: " n   n\n(>.< )\n /W\\ ",
			Sleep:  " n   n\n(-.-) \n /W\\ ",
			Tired:  " n   n\n(~.~ )\n /W\\ ",
			Dead:   " n   n\n(x.x )\n /W\\ ",
		},
	},
	{
		Name:    "Piko",
		Species: "Tavşan",
		Emoji:   "🐰",
		Colors:  CharColors{"225", "82", "196", "147", "208"},
		Art: CharArt{
			Normal: "(\\  /)\n(^.^)\n(___)  ",
			Happy:  "(\\  /)\n(^ω^)\n(___)  ",
			Hungry: "(\\  /)\n(>.< )\n(___)  ",
			Sleep:  "(\\  /)\n(-.- )\n(___)  ",
			Tired:  "(\\  /)\n(~.~ )\n(___)  ",
			Dead:   "(\\  /)\n(x.x )\n(___)  ",
		},
	},
	{
		Name:    "Boo",
		Species: "Hayalet",
		Emoji:   "👻",
		Colors:  CharColors{"189", "82", "196", "147", "208"},
		Art: CharArt{
			Normal: "  .~. \n (o.o)\n/~~~~~\\",
			Happy:  "  .~. \n (^ω^)\n/~~~~~\\",
			Hungry: "  .~. \n (>.< )\n/~~~~~\\",
			Sleep:  "  .~. \n (-.- )\n/~~~~~\\",
			Tired:  "  .~. \n (~.~ )\n/~~~~~\\",
			Dead:   "  .~. \n (x.x )\n/~~~~~\\",
		},
	},
}
