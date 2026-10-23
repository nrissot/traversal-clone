package game

import (
	"fmt"
	"slices"
	"strings"
)

// https://en.wikipedia.org/wiki/ANSI_escape_code

const RESET = "\x1b[0m"

// === UTILS ===

func GrayBG(code uint8, text string) string {
	text_code := 0
	if code <= 243 {
		text_code = 15
	}
	return fmt.Sprintf("\x1b[48:5:%d;38:5:%dm%s", code, text_code, text)
}

func RedBG(s string) string {
	return fmt.Sprintf("\x1b[101;1m%s", s)
}

func RedText(s string) string {
	return fmt.Sprintf("\x1b[31m%s%s", s, RESET)
}

func BlueBG(s string) string {
	return fmt.Sprintf("\x1b[104;1m%s", s)
}

func BlueText(s string) string {
	return fmt.Sprintf("\x1b[34m%s%s", s, RESET)
}

func GreenBG(s string) string {
	return fmt.Sprintf("\x1b[42;1m%s", s)
}

func WarningText(s string) string {
	return fmt.Sprintf("\x1b[33;1m%s%s", s, RESET)
}

// === DISPLAY ===

func (state BoardState) DisplayBoard(show_heights bool) {
	var s strings.Builder

	lastMove := ""
	if state.TurnNumber > 0 {
		if state.TurnNumber%2 == 1 {
			lastMove = fmt.Sprintf("%v", RedText(state.History[len(state.History)-1].String()))
		} else {
			lastMove = fmt.Sprintf("%v", BlueText(state.History[len(state.History)-1].String()))
		}
	}

	fmt.Fprintf(&s, "Turn n°%d (%d-%d) %s\n", state.TurnNumber, state.Score[0], state.Score[1], lastMove)

	localMins := []int{}
	localMaxs := []int{}
	for _, row := range state.Board {
		localMins = append(localMins, slices.Min(row))
		localMaxs = append(localMaxs, slices.Max(row))
	}
	min := slices.Min(localMins)
	max := slices.Max(localMaxs)

	fmt.Fprintln(&s, "    a  b  c  d  e  f  g  h  i  j")
	for y, row := range state.Board {
		if y > 0 {
			s.WriteString("\n")
		}
		fmt.Fprintf(&s, "%2d ", y+1)
		for x, cell := range row {
			c := "  "
			for _, p := range state.Pieces {
				if p.Alive && p.X == x && p.Y == y {
					if p.Color == Red {
						c = RedBG(p.String())
					} else {
						c = BlueBG(p.String())
					}
				}
			}
			code := 232 + uint8((float64(cell-min)/float64(max-min))*23)
			if show_heights {
				s.WriteString(GrayBG(code, fmt.Sprintf("%d", cell)))
			} else {
				s.WriteString(GrayBG(code, " "))
			}
			s.WriteString(c)
			s.WriteString(RESET)
		}
	}
	fmt.Println(s.String())
}

func (state BoardState) DisplayAllowedMoves(show_heights bool, moves []Move) {
	var s strings.Builder

	fmt.Fprintf(&s, "Displaying Moveset:\n%v\n", moves)

	localMins := []int{}
	localMaxs := []int{}
	for _, row := range state.Board {
		localMins = append(localMins, slices.Min(row))
		localMaxs = append(localMaxs, slices.Max(row))
	}
	min := slices.Min(localMins)
	max := slices.Max(localMaxs)

	fmt.Fprintln(&s, "    a  b  c  d  e  f  g  h  i  j")
	for y, row := range state.Board {
		if y > 0 {
			s.WriteString("\n")
		}
		fmt.Fprintf(&s, "%2d ", y+1)
		for x, cell := range row {
			code := 232 + uint8((float64(cell-min)/float64(max-min))*23)
			c := "  "
			for _, p := range state.Pieces {
				if p.Alive && p.X == x && p.Y == y {
					if p.Color == Red {
						c = RedBG(p.String())
					} else {
						c = BlueBG(p.String())
					}
				}
			}
			for _, m := range moves {
				if m.ToX == x && m.ToY == y {
					if m.Takes {
						c = GreenBG("x") + GrayBG(code, " ")
					} else {
						c = GreenBG("o") + GrayBG(code, " ")
					}
				}
			}
			if show_heights {
				s.WriteString(GrayBG(code, fmt.Sprintf("%d", cell)))
			} else {
				s.WriteString(GrayBG(code, " "))
			}
			s.WriteString(c)
			s.WriteString(RESET)
		}
	}
	fmt.Println(s.String())
}

func (resp RespStruct) DisplayBoard() {
	var s strings.Builder

	fmt.Fprintf(&s, "Board Type: %s (%d*%d)\n", resp.TerrainSelected, resp.Columns, resp.Rows)

	localMins := []int{}
	localMaxs := []int{}
	for _, row := range resp.Board {
		localMins = append(localMins, slices.Min(row))
		localMaxs = append(localMaxs, slices.Max(row))
	}
	min := slices.Min(localMins)
	max := slices.Max(localMaxs)

	for _, row := range resp.Board {
		s.WriteString("\n")
		for _, cell := range row {
			s.WriteString(GrayBG(232+uint8((float64(cell-min)/float64(max-min))*23), fmt.Sprintf("%2d %s", cell, RESET)))
		}
	}
	fmt.Println(s.String())
}
