package game

import (
	"fmt"
)

const DEBUG = false

type BoardState struct {
	Board        [][]int `json:"board"`
	Pieces       []Piece `json:"pieces"`
	TurnNumber   int     `json:"turn_number"`
	Score        [2]int  `json:"score"`
	HasTraversed int     `json:"has_traversed"`
	History      []Move  `json:"history"`
}

type PieceColor string
type PieceType string
type Direction string

const (
	Red  PieceColor = "Red"
	Blue PieceColor = "Blue"

	Explorer PieceType = "E"
	Vanguard PieceType = "V"
	Sentinel PieceType = "S"

	North Direction = "N"
	South Direction = "S"
	East  Direction = "E"
	West  Direction = "W"
)

type Piece struct {
	Color PieceColor `json:"color"`
	Type  PieceType  `json:"type"`
	Side  Direction  `json:"direction"`
	X     int        `json:"x"`
	Y     int        `json:"y"`
	Alive bool       `json:"alive"`
}

type Move struct {
	FromX int   `json:"from_x"`
	FromY int   `json:"from_y"`
	ToX   int   `json:"to_x"`
	ToY   int   `json:"to_y"`
	Piece Piece `json:"piece"`
	Takes bool  `json:"takes"`
}

func (p Piece) Speed() (speed int) {
	switch p.Type {
	case Explorer:
		return 5
	case Vanguard:
		return 4
	case Sentinel:
		return 3
	default:
		panic(fmt.Errorf("[Speed] (%v) unknown piece type (expected 'E', 'V' or 'S', got '%s')", p, p.Type))
	}
}

func (p Piece) Strength() (strength int) {
	switch p.Type {
	case Explorer:
		return 1
	case Vanguard:
		return 2
	case Sentinel:
		return 3
	default:
		panic(fmt.Errorf("[Strength] (%v) unknown piece type (expected 'E', 'V' or 'S', got '%s')", p, p.Type))
	}
}

func (p Piece) String() string {
	return fmt.Sprintf("%s%s", p.Type, p.Side)
}

func (p *Piece) DebugMoveTo(x, y int) {
	if !DEBUG {
		panic("[DEBUG MOVE] attempted a debug action in a non-debug scenario")
	}
	from, _ := CoordsToNotation(p.X, p.Y)
	to, _ := CoordsToNotation(x, y)
	fmt.Printf("[DEBUG MOVE]: %s %s (%s) -> (%s)\n", p.Color, p.String(), from, to)
	p.X = x
	p.Y = y
}

func (m Move) String() string {
	empty := Move{}
	if m == empty {
		return "####"
	}
	notation_position, err := CoordsToNotation(m.ToX, m.ToY)
	if err != nil {
		return "####"
	}
	if m.Takes {
		return fmt.Sprintf("%sx%s", m.Piece.String(), notation_position)
	} else {
		return fmt.Sprintf("%s%s", m.Piece.String(), notation_position)
	}
}

func InitializeBoardState(resp RespStruct) BoardState {
	return BoardState{
		Board: resp.Board,
		Pieces: []Piece{
			{Color: Red, Type: Sentinel, Side: West, X: 2, Y: 0, Alive: true}, // "c1"
			{Color: Red, Type: Vanguard, Side: West, X: 3, Y: 0, Alive: true}, // "d1"
			{Color: Red, Type: Explorer, Side: West, X: 4, Y: 0, Alive: true}, // "e1"
			{Color: Red, Type: Explorer, Side: East, X: 5, Y: 0, Alive: true}, // "f1"
			{Color: Red, Type: Vanguard, Side: East, X: 6, Y: 0, Alive: true}, // "g1"
			{Color: Red, Type: Sentinel, Side: East, X: 7, Y: 0, Alive: true}, // "h1"

			{Color: Blue, Type: Sentinel, Side: West, X: 2, Y: 9, Alive: true}, //"c10"
			{Color: Blue, Type: Vanguard, Side: West, X: 3, Y: 9, Alive: true}, //"d10"
			{Color: Blue, Type: Explorer, Side: West, X: 4, Y: 9, Alive: true}, //"e10"
			{Color: Blue, Type: Explorer, Side: East, X: 5, Y: 9, Alive: true}, //"f10"
			{Color: Blue, Type: Vanguard, Side: East, X: 6, Y: 9, Alive: true}, //"g10"
			{Color: Blue, Type: Sentinel, Side: East, X: 7, Y: 9, Alive: true}, //"h10"
		},
		TurnNumber:   0,
		Score:        [2]int{0, 0},
		History:      []Move{},
		HasTraversed: 0,
	}
}

func (b BoardState) At(x, y int) (height int, occupiedBy Piece, outside bool, err error) {
	if x < 0 || x > 9 {
		return height, occupiedBy, outside, fmt.Errorf("[BoardState.At] invalid x value (expected [0..9], got '%d')", x)
	}
	if y < 0 || y > 9 {
		if y == -1 {
			height = b.Board[0][x]
			return height, NO_PIECE, true, nil
		}
		if y == 10 {
			height = b.Board[9][x]
			return height, NO_PIECE, true, nil
		}
		return height, occupiedBy, outside, fmt.Errorf("[BoardState.At] invalid y value (expected [0..9], got '%d')", y)
	}
	height = b.Board[y][x]
	occupiedBy = NO_PIECE
	for _, p := range b.Pieces {
		if p.Alive && p.X == x && p.Y == y {
			occupiedBy = p
			break
		}
	}
	return height, occupiedBy, false, nil
}

func (b BoardState) GetPiecePointer(p Piece) (pointer *Piece, err error) {
	for i := range b.Pieces {
		if b.Pieces[i] == p {
			return &b.Pieces[i], nil
		}
	}
	return nil, fmt.Errorf("[GetPiecePointer] unknown piece %v", p)
}
