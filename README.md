# tamagotchi

A tiny virtual pet that lives in your terminal — now with a living diorama:
a day/night sky, **real-world weather**, drifting particles, and a pet that
**persists between runs** and ages while you're away.

Built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea) /
[Lipgloss](https://github.com/charmbracelet/lipgloss).

```
   /\    /\
  .--------.
  |  ^  ^  |
  |   v    |
  | \____/ |
  '--------'
   (_)  (_)
```

## What makes it pretty

- **Living diorama** — the pet sits in a composited scene with a sky, ground,
  and animated particles, all drawn on a custom cell canvas.
- **Day / night cycle** — the palette, sky glyphs and border shift with your
  real system clock (dawn, day, dusk, night).
- **Real weather** — pulls the current conditions for your location from
  [wttr.in](https://wttr.in) (no API key needed) and rains, snows or shines
  inside the diorama to match. Override the city with `TAMAGOTCHI_CITY`.
- **Gradient stat bars** — bars flow through a colour gradient instead of a
  flat fill (via [go-colorful](https://github.com/lucasb-eyer/go-colorful)).
- **Persistent life** — your pet is saved to your OS config dir and keeps
  ageing while the game is closed. Come back and it remembers you.
- **Expressive sprites** — rounded, cartoon-style faces with per-mood
  expressions, idle blinks, sleep z's and a feeding fish.

## Getting started

```bash
go run .
```

Or build and run:

```bash
go build -o tamagotchi && ./tamagotchi
```

## How to play

Pick one of six pets on the character select screen, then keep them happy and healthy.

| Key | Action        |
|-----|---------------|
| `f` | Feed          |
| `p` | Play          |
| `s` | Sleep / wake  |
| `n` | Adopt a new pet (saves the current one) |
| `1` | Minigame: Number guess   |
| `2` | Minigame: Reaction time  |
| `3` | Minigame: Rock Paper Scissors |
| `4` | Minigame: Math quiz      |
| `q` | Quit          |

Stats decay over time — even while the game is closed. If hunger gets too high, happiness drops too low, or energy runs out, your pet will let you know. Don't neglect them for too long.

## Configuration

| Env var            | Effect                                            |
|--------------------|---------------------------------------------------|
| `TAMAGOTCHI_CITY`  | Force a weather location (otherwise auto by IP).  |

The save file lives at `<os-config-dir>/tamagotchi/save.json`.

## Pets

| Name    | Species |
|---------|---------|
| Mochi   | Cat     |
| Biscuit | Dog     |
| Piko    | Rabbit  |
| Boo     | Ghost   |
| Kuma    | Bear    |
| Hoshi   | Hamster |

## Requirements

Go 1.21+
