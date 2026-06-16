package main

import (
	"fmt"
	"math/rand"
	"strings"
)

/**
 * Struct defines the state of a puzzle
 * tracking zero_loc allows moves to be found in constant time rather than n
 * tracking solved can be much simpler by checking if the last move puts both tiles in correct position
 **/
type Puzzle struct {
	arr       [][]int
	zero_loc  RowCol
	last_move Move
}

/**
 * Returns if the states of the two puzzles are equal. Does not care about metadata
 * like zero location or last move
 **/
func (p1 Puzzle) equals(p2 Puzzle) bool {
	if p1.size() != p2.size() {
		return false
	}

	for i := 0; i < p1.size(); i++ {
		if p1.getN(i) != p2.getN(i) {
			return false
		}
	}
	return true
}

func (p Puzzle) copy() Puzzle {
	arr_copy := make([][]int, p.len())
	for i := range p.arr {
		arr_copy[i] = make([]int, p.len())
		copy(arr_copy[i], p.arr[i])
	}

	return Puzzle{
		arr:       arr_copy,
		zero_loc:  p.zero_loc,
		last_move: p.last_move,
	}
}

/**
 * function called during creation of puzzle
 **/
func (p *Puzzle) set(rc RowCol, val int) {
	if val == 0 { // track zero_loc
		p.zero_loc = rc
	}
	p.arr[rc.row][rc.col] = val
}

func newPuzzle(arr []int) Puzzle {
	// check puzzle is valid
	isSquare, len := isSquare(len(arr))
	if !isSquare || !containsAllIndices(arr) {
		panic("Invalid input arr for puzzle")
	}

	// create puzzle array
	var p = Puzzle{
		arr:       make([][]int, len),
		zero_loc:  RowCol{0, 0},
		last_move: None,
	}

	// alloc rows
	for i := range p.arr {
		p.arr[i] = make([]int, len)
	}

	// init puzzle array
	for i, e := range arr {
		p.set(p.nToCoord(i), e)
	}

	return p
}

func newPuzzleSolved(len int) Puzzle {
	var arr = make([]int, len*len)
	for i := range arr {
		arr[i] = i
	}
	return newPuzzle(arr)
}

func newPuzzleSwapped(n int, swaps int) Puzzle {
	rand.Seed(config.Random_seed)
	var p Puzzle = newPuzzleSolved(n)

	// make n random moves on the board
	for i := 0; i < swaps; i++ {
		moves := p.getNewMoves()
		p.makeMove(moves[rand.Intn(len(moves))])
	}

	p.last_move = None

	return p
}

func newPuzzleMisplaced(n int, misplaced int) (Puzzle, int) {
	if misplaced > n*n-1 {
		misplaced = n*n - 1
	} else if misplaced < 0 {
		misplaced = 0
	}

	rand.Seed(config.Random_seed)
	var p Puzzle = newPuzzleSolved(n)

	// make n random moves on the board
	var swaps int = 0
	for int(h1(p)) < misplaced {
		moves := p.getNewMoves()
		p.makeMove(moves[rand.Intn(len(moves))])
		swaps++
	}

	p.last_move = None

	return p, swaps
}

func (p Puzzle) nToRow(n int) int {
	return n / p.len()
}

func (p Puzzle) nToCol(n int) int {
	return n % p.len()
}

func (p Puzzle) nToCoord(n int) RowCol {
	return RowCol{
		row: p.nToRow(n),
		col: p.nToCol(n),
	}
}

func (p Puzzle) len() int {
	return len(p.arr)
}

func (p Puzzle) size() int {
	return p.len() * p.len()
}

func (p Puzzle) get(rc RowCol) int {
	return p.arr[rc.row][rc.col]
}

func (p Puzzle) getN(n int) int {
	return p.get(p.nToCoord(n))
}

func (p Puzzle) getGoalPos(val int) RowCol {
	return p.nToCoord(val)
}

/**
 * Prints out the state of the puzzle
 * only works up to 2 digit numbers
 * looks weird at puzzle len 10, which is infeasible anyways
 */
func (p Puzzle) print() {
	fmt.Print(p.toStr())
}

func (p Puzzle) toStr() string {
	var sb strings.Builder
	var len int = p.len()
	for i := 0; i < len; i++ {
		sb.WriteString(strings.Repeat("+----", len) + "+\n")
		for j := 0; j < len; j++ {
			sb.WriteString("| ")

			e := p.get(RowCol{
				row: i,
				col: j,
			})

			if e == 0 {
				sb.WriteString("  ")
			} else {
				if e < 10 {
					sb.WriteString(" ")
				}
				sb.WriteString(fmt.Sprint(e))
			}

			sb.WriteString(" ")
		}
		sb.WriteString("|\n")
	}

	sb.WriteString(strings.Repeat("+----", len) + "+\n")
	return sb.String()
}

func (p Puzzle) isSolved() bool {
	for i := 0; i < p.size(); i++ {
		if i != p.getN(i) {
			return false
		}
	}

	return true
}

func (p Puzzle) getMoves() []Move {
	var moves []Move

	if p.zero_loc.row > 0 {
		moves = append(moves, Down)
	}

	if p.zero_loc.row < p.len()-1 {
		moves = append(moves, Up)
	}

	if p.zero_loc.col > 0 {
		moves = append(moves, Right)
	}

	if p.zero_loc.col < p.len()-1 {
		moves = append(moves, Left)
	}

	return moves
}

func (p Puzzle) getNewMoves() []Move {
	var moves []Move

	if p.zero_loc.row > 0 && p.last_move != Up {
		moves = append(moves, Down)
	}

	if p.zero_loc.row < p.len()-1 && p.last_move != Down {
		moves = append(moves, Up)
	}

	if p.zero_loc.col > 0 && p.last_move != Left {
		moves = append(moves, Right)
	}

	if p.zero_loc.col < p.len()-1 && p.last_move != Right {
		moves = append(moves, Left)
	}

	return moves
}

/*
 * Returns the successors to a puzzle
 * ignore_prev_moves doesn't include the opposite of the last move played,
 * as it's already been explored. false will include them.
 */
func (p Puzzle) getSuccessors(ignore_prev_moves bool) []Puzzle {
	var successors []Puzzle = []Puzzle{}

	var moves []Move
	if ignore_prev_moves {
		moves = p.getNewMoves()
	} else {
		moves = p.getMoves()
	}

	for _, move := range moves {
		successors = append(successors, p.tryMove(move))
	}

	return successors
}

func (p *Puzzle) swap(rc1 RowCol, rc2 RowCol) {
	var temp int = p.get(rc1)
	p.set(rc1, p.get(rc2))
	p.set(rc2, temp)
}

/**
 * makes a move on the puzzle
 **/
func (p *Puzzle) makeMove(m Move) {
	p.last_move = m
	switch m {
	case Up:
		p.swap(p.zero_loc, RowCol{
			row: p.zero_loc.row + 1,
			col: p.zero_loc.col,
		})

	case Down:
		p.swap(p.zero_loc, RowCol{
			row: p.zero_loc.row - 1,
			col: p.zero_loc.col,
		})

	case Left:
		p.swap(p.zero_loc, RowCol{
			row: p.zero_loc.row,
			col: p.zero_loc.col + 1,
		})

	case Right:
		p.swap(p.zero_loc, RowCol{
			row: p.zero_loc.row,
			col: p.zero_loc.col - 1,
		})
	}
}

/**
 * returns a copy of the puzzle with the move made
 **/
func (p Puzzle) tryMove(m Move) Puzzle {
	p = p.copy()
	p.makeMove(m)
	return p
}

/**
 * returns if p is a successor to prev, meaning it takes 1 move to go between p and prev
 **/
func (p Puzzle) isSuccessorTo(prev Puzzle) bool {
	for _, m := range prev.getMoves() {
		if prev.tryMove(m).equals(p) {
			return true
		}
	}
	return false
}
