# tamagotchi

A tiny virtual pet that lives in your terminal.

Built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

```
 /\ /\
(=• ω •=)
 ( ___ )
  /   \
 (_) (_)
```

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
| `1` | Minigame: Number guess   |
| `2` | Minigame: Reaction time  |
| `3` | Minigame: Rock Paper Scissors |
| `4` | Minigame: Math quiz      |
| `q` | Quit          |

Stats decay over time. If hunger gets too high, happiness drops too low, or energy runs out, your pet will let you know. Don't neglect them for too long.

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
