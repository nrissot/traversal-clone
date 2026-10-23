package game

import (
	"fmt"
	"slices"
	"strings"
)

type cell struct {
	X       int
	Y       int
	H       int
	P       Piece
	Outside bool
}

var NO_PIECE = Piece{}

func (b *BoardState) ApplyLegalMove(m Move) (err error) {
	movingPiecePointer, err := b.GetPiecePointer(m.Piece)
	if err != nil {
		return fmt.Errorf("[ApplyLegalMove]:%v", err)
	}
	if m.Takes {
		_, killedPiece, _, err := b.At(m.ToX, m.ToY)
		if err != nil {
			return fmt.Errorf("[ApplyLegalMove]:%v", err)
		}
		killedPiecePointer, err := b.GetPiecePointer(killedPiece)
		if err != nil {
			return fmt.Errorf("[ApplyLegalMove]:%v", err)
		}
		killedPiecePointer.Alive = false
	}
	movingPiecePointer.X = m.ToX
	movingPiecePointer.Y = m.ToY
	// finisher
	if m.ToY == -1 {
		movingPiecePointer.Alive = false
		b.Score[b.TurnNumber%2] += movingPiecePointer.Strength()
		b.HasTraversed++
		// win detection
		if movingPiecePointer.Type == Sentinel || b.HasTraversed >= 2 {
			if b.Score[0] > b.Score[1] {
				return fmt.Errorf(RedText(fmt.Sprintf("[END OF GAME] %s won! (%d-%d)", movingPiecePointer.Color, b.Score[0], b.Score[1])))
			} else {
				return fmt.Errorf(BlueText(fmt.Sprintf("[END OF GAME] %s won! (%d-%d)", movingPiecePointer.Color, b.Score[0], b.Score[1])))
			}
		}
	}
	b.History = append(b.History, m)
	b.TurnNumber++
	return nil
}

func (b *BoardState) ApplyMove(m Move) (err error) {
	legal, err := m.IsMoveLegal(*b)
	if err != nil {
		return fmt.Errorf("[ApplyMove] cannot compute move legality: %v", err)
	}
	if !legal {
		return fmt.Errorf("[ApplyMove] (Illegal Move) %v", m)
	}
	return b.ApplyLegalMove(m)
}

func (b BoardState) LegalMoveFromString(s string) (m Move, err error) {
	var targetPieceType PieceType
	var targetPieceSide Direction
	var targetColor PieceColor

	if b.TurnNumber%2 == 0 {
		targetColor = Red
	} else {
		targetColor = Blue
	}

	switch string(s[0]) {
	case "S", "V", "E":
		targetPieceType = PieceType(string(s[0]))
	default:
		return m, fmt.Errorf("[LegalMoveFromString] %s invalid piece type identifier (expected 'S', 'V' or 'E', got '%s')", s, string(s[0]))
	}
	switch string(s[1]) {
	case "E", "W":
		targetPieceSide = Direction(string(s[1]))
	default:
		return m, fmt.Errorf("[LegalMoveFromString] %s invalid piece side identifier (expected 'W' or 'E', got '%s')", s, string(s[1]))
	}

	var p Piece = NO_PIECE
	flag := true
	for _, candidate := range b.Pieces {
		if candidate.Alive && candidate.Color == targetColor && candidate.Type == targetPieceType && candidate.Side == targetPieceSide {
			p = candidate
			flag = false
			break
		}
	}
	if flag {
		alivePiece := []string{}
		for _, toprintpiece := range b.Pieces {
			if toprintpiece.Alive {
				alivePiece = append(alivePiece, toprintpiece.String())
			}
		}
		return m, fmt.Errorf("[LegalMoveFromString] %s unable to find an alive piece matching the provided identifier (expected %v ,got %s%s)", s, b.Pieces, targetPieceType, targetPieceSide)
	}

	coordsPart, _ := strings.CutPrefix(string(s[2:]), "x")
	x, y, err := NotationToCoords(string(coordsPart))
	if err != nil {
		return m, fmt.Errorf("[LegalMoveFromString] %s unable to process coordinates: %v", s, err)
	}

	var takes bool = false
	if strings.Contains(s, "x") {
		takes = true
	}

	m = Move{FromX: p.X, FromY: p.Y, ToX: x, ToY: y, Piece: p, Takes: takes}

	ok, err := m.IsMoveLegal(b)
	if err != nil {
		return m, fmt.Errorf("[LegalMoveFromString] unable to check move legality %v: %v", m, err)
	}
	if !ok {
		moves, _ := b.GetLegalMoves(m.Piece)
		return m, fmt.Errorf("[LegalMoveFromString] Illegal Move (expected %v, got '%v')", moves, m)
	}
	return m, nil
}

func (m Move) IsMoveLegal(b BoardState) (is_move_legal bool, err error) {
	legal_moves, err := b.GetLegalMoves(m.Piece)
	if err != nil {
		return is_move_legal, fmt.Errorf("[IsMoveLegal] %v cannot compute a list of legal moves: %v", m, err)
	}
	return slices.Contains(legal_moves, m), err
}

func (b BoardState) GetAllLegalMoves() (moves []Move, err error) {
	targetColor := Blue
	if b.TurnNumber%2 == 0 {
		targetColor = Red
	}
	for _, p := range b.Pieces {
		if p.Color == targetColor {
			piece_moves, err := b.GetLegalMoves(p)
			if err != nil {
				return moves, fmt.Errorf("[GetAllLegalMoves] error while computing legal moves for %v :%v", p, err)
			}
			moves = append(moves, piece_moves...)
		}
	}
	return moves, nil
}

func (b BoardState) GetLegalMoves(p Piece) (moves []Move, err error) {
	if !slices.Contains(b.Pieces, p) {
		return moves, fmt.Errorf("[GetLegalMoves] piece %v doesnt exist in the boardState %v", p, b)
	}
	speed := p.Speed()

	//    v cant start outside
	h, _, _, err := b.At(p.X, p.Y)

	initial_cell := cell{p.X, p.Y, h, p, false}

	var cost_table map[cell]int = make(map[cell]int)
	var kills map[cell]int = make(map[cell]int)
	var frontier []cell = []cell{}
	var new_frontier []cell = []cell{}

	frontier = append(frontier, initial_cell)
	cost_table[initial_cell] = 0

	if p.Type == Explorer {
		jumps := b.ExplorerJump(p.X, p.Y)
		for _, jump := range jumps {
			if jump.Outside && ((p.Color == Red && jump.Y == -1) || (p.Color == Blue && jump.Y == 10)) {
				continue
			}
			if jump.P == NO_PIECE {
				cost_table[jump] = speed
			} else {
				can_attack, err := p.CanAttack(b, jump.P, 5)
				if err != nil {
					return moves, fmt.Errorf("[GetLegalMoves] invalid attack computation: %v", err)
				} else if can_attack {
					kills[jump] = speed
				}
			}
		}
	}

	// Backtracking BFS
	for {
		for _, candidate := range frontier {
			neighborhood := b.Around(candidate.X, candidate.Y)
			for _, neighbor := range neighborhood {
				new_cost := cost_table[candidate] + candidate._distanceToNeighbor(neighbor)
				if neighbor.Outside && ((p.Color == Red && neighbor.Y == -1) || (p.Color == Blue && neighbor.Y == 10)) {
					continue
				}
				if neighbor.P == NO_PIECE {
					old_cost, ok := cost_table[neighbor]
					if ok && (new_cost < old_cost) || !ok && new_cost <= speed {
						cost_table[neighbor] = new_cost
						if !neighbor.Outside {
							new_frontier = append(new_frontier, neighbor)
						}
					}
				} else {
					can_attack, err := p.CanAttack(b, neighbor.P, new_cost)
					if err != nil {
						return moves, fmt.Errorf("[GetLegalMoves] invalid attack computation: %v", err)
					}
					old_cost, ok := kills[neighbor]
					if can_attack && ((ok && new_cost < old_cost) || (!ok && new_cost <= speed)) {
						kills[neighbor] = new_cost
					}

				}
			}
		}
		if len(new_frontier) != 0 {
			frontier = new_frontier
			new_frontier = []cell{}
		} else {
			break
		}
	}

	for candidate, cost := range cost_table {
		//               v useless check
		if (cost > 0) && (cost <= speed) {
			if !candidate.Outside {
				moves = append(moves, Move{
					FromX: p.X,
					FromY: p.Y,
					ToX:   candidate.X,
					ToY:   candidate.Y,
					Piece: p,
					Takes: false,
				})
			} else {
				moves = append(moves, Move{
					FromX: p.X,
					FromY: p.Y,
					ToX:   candidate.X,
					ToY:   -1,
					Piece: p,
					Takes: false,
				})
			}
		}
	}
	for candidate, cost := range kills {
		if (cost > 0) && (cost <= speed) {
			moves = append(moves, Move{
				FromX: p.X,
				FromY: p.Y,
				ToX:   candidate.X,
				ToY:   candidate.Y,
				Piece: p,
				Takes: true,
			})
		}
	}

	return moves, err
}

func (b BoardState) Around(x, y int) (neighborhood []cell) {
	for _, neighbor := range [4][2]int{{x, y + 1}, {x, y - 1}, {x + 1, y}, {x - 1, y}} {
		nx := neighbor[0]
		ny := neighbor[1]
		h, p, outside, err := b.At(nx, ny)
		if err == nil {
			neighborhood = append(neighborhood, cell{nx, ny, h, p, outside})
		}
	}
	return neighborhood
}

func (b BoardState) ExplorerJump(x, y int) (neighborhood []cell) {
	for _, neighbor := range [8][2]int{{x, y + 2}, {x, y - 2}, {x + 2, y}, {x - 2, y}, {x + 1, y + 1}, {x - 1, y - 1}, {x + 1, y - 1}, {x - 1, y + 1}} {
		nx := neighbor[0]
		ny := neighbor[1]
		h, p, o, err := b.At(nx, ny)
		if err == nil {
			neighborhood = append(neighborhood, cell{nx, ny, h, p, o})
		}
	}
	return neighborhood
}

func (p_atk Piece) CanAttack(b BoardState, p_def Piece, distance int) (can_attack bool, err error) {
	if !slices.Contains(b.Pieces, p_def) {
		return can_attack, fmt.Errorf("[CanAttack] piece %v doesnt exist in the boardState %v", p_def, b)
	}
	if !p_atk.Alive {
		return can_attack, fmt.Errorf("[CanAttack] (ZOMBIE) attacking piece %v is dead, cannot attack", p_atk)
	}
	if !p_def.Alive {
		return can_attack, fmt.Errorf("[CanAttack] (ZOMBIE) defending piece %v is dead, cannot defend", p_def)
	}
	if p_atk.Color == p_def.Color {
		return false, nil
	}
	atk := p_atk.Strength()
	def := p_def.Strength()
	neighborhood := b.Around(p_def.X, p_def.Y)
	for _, neighbor := range neighborhood {
		if neighbor.P != NO_PIECE && neighbor.P != p_atk {
			if neighbor.P.Color == p_atk.Color {
				// if there is an adjacent ally of p_atk
				atk += neighbor.P.Strength()
			} else {
				// if there is an adjacent ally of p_def
				def += 1
			}
		}
	}
	if p_atk.Type == Vanguard && distance == 4 {
		atk += 1
	}
	h_atk, _, o, err := b.At(p_atk.X, p_atk.Y)
	if err != nil {
		return can_attack, fmt.Errorf("[CanAttack] %v positions: %v", p_atk, err)
	}
	if o {
		return false, nil
	}
	h_def, _, o, err := b.At(p_def.X, p_def.Y)
	if err != nil {
		return can_attack, fmt.Errorf("[CanAttack] %v positions: %v", p_def, err)
	}
	if o {
		return false, nil
	}
	if h_atk > h_def {
		atk += 1
	} else if h_def > h_atk {
		def += 1
	}
	return atk >= def, nil
}

func (c1 cell) _distanceToNeighbor(c2 cell) (distance int) {
	return 1 + max(c1.H, c2.H) - min(c1.H, c2.H)
}
