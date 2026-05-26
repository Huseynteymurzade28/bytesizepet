package main

import "github.com/charmbracelet/lipgloss"

// CharArt holds all emotional states plus an alternate idle frame (NormalB)
// for the blink/idle animation. Every field must have the same line count.
type CharArt struct {
	Normal, NormalB                          string
	Happy, Hungry, Sleep, Tired, Dead string
}

type CharColors struct {
	Normal, Happy, Hungry, Sleep, Tired lipgloss.Color
}

type CharDef struct {
	Name     string
	Species  string
	Emoji    string
	Colors   CharColors
	Art      CharArt
	HabitatA string // habitat frame A  (9 chars wide × 2 lines, separated by \n)
	HabitatB string // habitat frame B
}

// AllChars is the selectable roster shown on the character-select screen.
var AllChars = []CharDef{
	{
		Name:    "Mochi",
		Species: "Cat",
		Emoji:   "🐱",
		Colors:  CharColors{"#e8b4d0", "#f5d0a0", "#e88a70", "#b4a8e0", "#b0a0c0"},
		Art: CharArt{
			Normal:  " /\\ /\\  \n(=• ω •=)\n ( ___ ) \n  /   \\  \n (_) (_)",
			NormalB: " /\\ /\\  \n(=^ ω ^=)\n ( ___ ) \n  /   \\  \n (_) (_)",
			Happy:   " /\\ /\\  \n(=^ ω ^=)\n ( ♡♡♡ )\n  /   \\  \n  ♪   ♪ ",
			Hungry:  " /\\ /\\  \n(=° A °=)\n (~~~~~) \n  /   \\  \n (_) (_)",
			Sleep:   " /\\ /\\  \n(=- w -=)\n ( ___ ) \n  /   \\  \n (_) (_)",
			Tired:   " /\\ /\\  \n(=~ . ~=)\n ( ___ ) \n  /   \\  \n (_) (_)",
			Dead:    " /\\ /\\  \n(=x . x=)\n (     ) \n  /   \\  \n (_) (_)",
		},
		HabitatA: "~·~·~·~·~\n  (@ )   ",
		HabitatB: "·~·~·~·~·\n  ( @)   ",
	},
	{
		Name:    "Biscuit",
		Species: "Dog",
		Emoji:   "🐶",
		Colors:  CharColors{"#d4a87a", "#a0d870", "#e87858", "#a8c0d8", "#b0a890"},
		Art: CharArt{
			Normal:  " n     n \n(ᵔ .  ᵔ)\n (     ) \n  /   \\  \n  W   W ",
			NormalB: " n     n \n(^ .  ^)\n (     ) \n  /   \\  \n  W   W ",
			Happy:   " n     n \n(ᵔ ω  ᵔ)\n ( ♡♡♡ )\n  /   \\  \n  ♪   ♪ ",
			Hungry:  " n     n \n(° A  °)\n (~~~~~) \n  /   \\  \n  W   W ",
			Sleep:   " n     n \n(- .  -)\n (     ) \n  /   \\  \n  W   W ",
			Tired:   " n     n \n(~ .  ~)\n (     ) \n  /   \\  \n  W   W ",
			Dead:    " n     n \n(x .  x)\n (     ) \n  /   \\  \n  W   W ",
		},
		HabitatA: ".v.v.v.v.\n  * .  * ",
		HabitatB: "v.v.v.v.v\n  . *  . ",
	},
	{
		Name:    "Piko",
		Species: "Rabbit",
		Emoji:   "🐰",
		Colors:  CharColors{"#f5d0e8", "#c8f0c8", "#f0a0a0", "#c0b8e8", "#c0b0c0"},
		Art: CharArt{
			Normal:  " (\\ /)  \n (• ω •)\n (     ) \n  |   |  \n  (___)  ",
			NormalB: " (\\ /)  \n (^ ω ^)\n (     ) \n  |   |  \n  (___)  ",
			Happy:   " (\\ /)  \n (^ ω ^)\n ( ♡♡♡ )\n  |   |  \n  ♪___♪  ",
			Hungry:  " (\\ /)  \n (° A °)\n (~~~~~) \n  |   |  \n  (___)  ",
			Sleep:   " (\\ /)  \n (- w -)\n (     ) \n  |   |  \n  (___)  ",
			Tired:   " (\\ /)  \n (~ . ~)\n (     ) \n  |   |  \n  (___)  ",
			Dead:    " (\\ /)  \n (x . x)\n (     ) \n  |   |  \n  (___)  ",
		},
		HabitatA: "^.^.^.^.^\n ,*, . , ",
		HabitatB: ".^.^.^.^.\n . ,*, . ",
	},
	{
		Name:    "Boo",
		Species: "Ghost",
		Emoji:   "👻",
		Colors:  CharColors{"#d8d0f0", "#ffffff", "#e87878", "#b0a8e0", "#a8a0c0"},
		Art: CharArt{
			Normal:  "  .~~~.  \n (o   o) \n  ( u )  \n /~~~~~\\ \n  \\___/  ",
			NormalB: "  .~~~.  \n (O   O) \n  ( u )  \n /~~~~~\\ \n  \\___/  ",
			Happy:   "  .~~~.  \n (^ ^ )  \n  ( ♡ )  \n /~~~~~\\ \n  \\^_^/  ",
			Hungry:  "  .~~~.  \n (o   o) \n  (~~~)  \n /~~~~~\\ \n  \\___/  ",
			Sleep:   "  .~~~.  \n (- - )  \n  ( u )  \n /~~~~~\\ \n  \\___/  ",
			Tired:   "  .~~~.  \n (~ ~ )  \n  ( u )  \n /~~~~~\\ \n  \\___/  ",
			Dead:    "  .~~~.  \n (x   x) \n  (   )  \n /~~~~~\\ \n  \\___/  ",
		},
		HabitatA: "* . · * .\n  · * .  ",
		HabitatB: ". · * . ·\n  * . ·  ",
	},
	{
		Name:    "Kuma",
		Species: "Bear",
		Emoji:   "🐻",
		Colors:  CharColors{"#c4a07a", "#f0c070", "#e07050", "#b0a8d0", "#a89880"},
		Art: CharArt{
			Normal:  "(U)   (U)\n (• ω •) \n (     ) \n  |   |  \n (_____) ",
			NormalB: "(U)   (U)\n (^ ω ^) \n (     ) \n  |   |  \n (_____) ",
			Happy:   "(U)   (U)\n (^ ω ^) \n ( ♡♡♡ ) \n  |   |  \n  ♪   ♪  ",
			Hungry:  "(U)   (U)\n (° A °) \n (~~~~~) \n  |   |  \n (_____) ",
			Sleep:   "(U)   (U)\n (- w -) \n (     ) \n  |   |  \n (_____) ",
			Tired:   "(U)   (U)\n (~ . ~) \n (     ) \n  |   |  \n (_____) ",
			Dead:    "(U)   (U)\n (x . x) \n (     ) \n  |   |  \n (_____) ",
		},
		HabitatA: "^.^.^.^.^\n  ,*, *  ",
		HabitatB: ".^.^.^.^.\n  * ,*,  ",
	},
	{
		Name:    "Hoshi",
		Species: "Hamster",
		Emoji:   "🐹",
		Colors:  CharColors{"#f0c8a0", "#f8e870", "#e87860", "#b8b0e0", "#c0b0a0"},
		Art: CharArt{
			Normal:  "  .---.  \n (˘ . ˘) \n (  u  ) \n /|   |\\ \n (_) (_) ",
			NormalB: "  .---.  \n (^ . ^) \n (  u  ) \n /|   |\\ \n (_) (_) ",
			Happy:   "  .---.  \n (˘ ω ˘) \n (  ♡  ) \n /|   |\\ \n  ♪   ♪  ",
			Hungry:  "  .---.  \n (° A °) \n (~   ~) \n /|   |\\ \n (_) (_) ",
			Sleep:   "  .---.  \n (- . -) \n (  u  ) \n /|   |\\ \n (_) (_) ",
			Tired:   "  .---.  \n (~ . ~) \n (  u  ) \n /|   |\\ \n (_) (_) ",
			Dead:    "  .---.  \n (x . x) \n (     ) \n /|   |\\ \n (_) (_) ",
		},
		HabitatA: ".,.,.,.,.\n  * . ,  ",
		HabitatB: ",.,.,.,.,\n  , . *  ",
	},
}
