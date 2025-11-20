package main

import "sync"

const (
	rows = 6
	cols = 7
)

type Cell int

const (
	Empty Cell = iota
	P1
	P2
)

type State string

const (
	Playing State = "playing"
	WinP1   State = "win_p1"
	WinP2   State = "win_p2"
	Draw    State = "draw"
)

type Game struct {
	mu            sync.Mutex
	Board         [rows][cols]Cell
	CurrentPlayer Cell
	State         State
	Message       string
}

func NewGame() *Game {
	return &Game{
		CurrentPlayer: P1,
		State:         Playing,
	}
}

func (g *Game) Play(col int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State != Playing {
		g.Message = "Partie terminée. Cliquez sur « Reset » pour recommencer."
		return false
	}
	if col < 0 || col >= cols {
		g.Message = "Colonne invalide."
		return false
	}

	row := -1
	for r := rows - 1; r >= 0; r-- {
		if g.Board[r][col] == Empty {
			row = r
			break
		}
	}
	if row == -1 {
		g.Message = "Cette colonne est pleine."
		return false
	}

	g.Board[row][col] = g.CurrentPlayer

	if g.isWinningMove(row, col) {
		if g.CurrentPlayer == P1 {
			g.State = WinP1
			g.Message = "🎉 Victoire du Joueur 1 !"
		} else {
			g.State = WinP2
			g.Message = "🎉 Victoire du Joueur 2 !"
		}
		return true
	}
	if g.isDraw() {
		g.State = Draw
		g.Message = "Match nul 🤝"
		return true
	}

	if g.CurrentPlayer == P1 {
		g.CurrentPlayer = P2
	} else {
		g.CurrentPlayer = P1
	}
	g.Message = ""
	return true
}

func (g *Game) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Board = [rows][cols]Cell{}
	g.CurrentPlayer = P1
	g.State = Playing
	g.Message = ""
}

func (g *Game) isDraw() bool {
	for c := 0; c < cols; c++ {
		if g.Board[0][c] == Empty {
			return false
		}
	}
	return true
}

func (g *Game) isWinningMove(r, c int) bool {
	player := g.Board[r][c]
	if player == Empty {
		return false
	}

	dirs := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, d := range dirs {
		count := 1
		count += g.countDirection(r, c, d[0], d[1], player)
		count += g.countDirection(r, c, -d[0], -d[1], player)
		if count >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) countDirection(r, c, dr, dc int, player Cell) int {
	cnt := 0
	rr, cc := r+dr, c+dc
	for rr >= 0 && rr < rows && cc >= 0 && cc < cols && g.Board[rr][cc] == player {
		cnt++
		rr += dr
		cc += dc
	}
	return cnt
}
