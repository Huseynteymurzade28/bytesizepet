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
		Species: "Cat",
		Emoji:   "🐱",
		Colors:  CharColors{"226", "82", "196", "147", "208"},
		Art: CharArt{
			Normal: " /\\ /\\  \n(=• ω •=)\n ( ___ ) \n  /   \\  \n (_) (_)",
			Happy:  " /\\ /\\  \n(=^ ω ^=)\n ( ♡♡♡ )\n  /   \\  \n  ♪   ♪ ",
			Hungry: " /\\ /\\  \n(=° A °=)\n (~~~~~) \n  /   \\  \n (_) (_)",
			Sleep:  " /\\ /\\  \n(=- w -=)\n ( ___ ) \n  /   \\  \n (_) (_)",
			Tired:  " /\\ /\\  \n(=~ . ~=)\n ( ___ ) \n  /   \\  \n (_) (_)",
			Dead:   " /\\ /\\  \n(=x . x=)\n (     ) \n  /   \\  \n (_) (_)",
		},
	},
	{
		Name:    "Biscuit",
		Species: "Dog",
		Emoji:   "🐶",
		Colors:  CharColors{"178", "82", "196", "147", "208"},
		Art: CharArt{
			Normal: " n     n \n(ᵔ .  ᵔ)\n (     ) \n  /   \\  \n  W   W ",
			Happy:  " n     n \n(ᵔ ω  ᵔ)\n ( ♡♡♡ )\n  /   \\  \n  ♪   ♪ ",
			Hungry: " n     n \n(° A  °)\n (~~~~~) \n  /   \\  \n  W   W ",
			Sleep:  " n     n \n(- .  -)\n (     ) \n  /   \\  \n  W   W ",
			Tired:  " n     n \n(~ .  ~)\n (     ) \n  /   \\  \n  W   W ",
			Dead:   " n     n \n(x .  x)\n (     ) \n  /   \\  \n  W   W ",
		},
	},
	{
		Name:    "Piko",
		Species: "Rabbit",
		Emoji:   "🐰",
		Colors:  CharColors{"225", "82", "196", "147", "208"},
		Art: CharArt{
			Normal: " (\\ /)  \n (• ω •)\n (     ) \n  |   |  \n  (___)  ",
			Happy:  " (\\ /)  \n (^ ω ^)\n ( ♡♡♡ )\n  |   |  \n  ♪___♪  ",
			Hungry: " (\\ /)  \n (° A °)\n (~~~~~) \n  |   |  \n  (___)  ",
			Sleep:  " (\\ /)  \n (- w -)\n (     ) \n  |   |  \n  (___)  ",
			Tired:  " (\\ /)  \n (~ . ~)\n (     ) \n  |   |  \n  (___)  ",
			Dead:   " (\\ /)  \n (x . x)\n (     ) \n  |   |  \n  (___)  ",
		},
	},
	{
		Name:    "Boo",
		Species: "Ghost",
		Emoji:   "👻",
		Colors:  CharColors{"189", "255", "196", "147", "208"},
		Art: CharArt{
			Normal: "  .~~~.  \n (o   o) \n  ( u )  \n /~~~~~\\ \n  \\___/  ",
			Happy:  "  .~~~.  \n (^ ^ )  \n  ( ♡ )  \n /~~~~~\\ \n  \\^_^/  ",
			Hungry: "  .~~~.  \n (o   o) \n  (~~~)  \n /~~~~~\\ \n  \\___/  ",
			Sleep:  "  .~~~.  \n (- - )  \n  ( u )  \n /~~~~~\\ \n  \\___/  ",
			Tired:  "  .~~~.  \n (~ ~ )  \n  ( u )  \n /~~~~~\\ \n  \\___/  ",
			Dead:   "  .~~~.  \n (x   x) \n  (   )  \n /~~~~~\\ \n  \\___/  ",
		},
	},
	{
		Name:    "Kuma",
		Species: "Bear",
		Emoji:   "🐻",
		Colors:  CharColors{"180", "82", "196", "147", "208"},
		Art: CharArt{
			Normal: "(U)   (U)\n (• ω •) \n (     ) \n  |   |  \n (_____) ",
			Happy:  "(U)   (U)\n (^ ω ^) \n ( ♡♡♡ ) \n  |   |  \n  ♪   ♪  ",
			Hungry: "(U)   (U)\n (° A °) \n (~~~~~) \n  |   |  \n (_____) ",
			Sleep:  "(U)   (U)\n (- w -) \n (     ) \n  |   |  \n (_____) ",
			Tired:  "(U)   (U)\n (~ . ~) \n (     ) \n  |   |  \n (_____) ",
			Dead:   "(U)   (U)\n (x . x) \n (     ) \n  |   |  \n (_____) ",
		},
	},
	{
		Name:    "Hoshi",
		Species: "Hamster",
		Emoji:   "🐹",
		Colors:  CharColors{"220", "226", "196", "147", "208"},
		Art: CharArt{
			Normal: "  .---.  \n (˘ . ˘) \n (  u  ) \n /|   |\\ \n (_) (_) ",
			Happy:  "  .---.  \n (˘ ω ˘) \n (  ♡  ) \n /|   |\\ \n  ♪   ♪  ",
			Hungry: "  .---.  \n (° A °) \n (~   ~) \n /|   |\\ \n (_) (_) ",
			Sleep:  "  .---.  \n (- . -) \n (  u  ) \n /|   |\\ \n (_) (_) ",
			Tired:  "  .---.  \n (~ . ~) \n (  u  ) \n /|   |\\ \n (_) (_) ",
			Dead:   "  .---.  \n (x . x) \n (     ) \n /|   |\\ \n (_) (_) ",
		},
	},
}
