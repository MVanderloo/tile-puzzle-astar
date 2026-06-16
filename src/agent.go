package main

import (
	"math"
	"time"
)

// type State interface {
// 	equals(State) bool
// 	getSuccessors() []State
// 	heuristic() float64
// 	print()
// 	isFinal() bool
// }

type Node struct {
	state Puzzle
	g     int8    // cost
	h     float32 // heuristic
	prev  *Node   // the predecessor
}

func (n Node) getF() float32 {
	return float32(n.g) + n.h
}

type Status string

const (
	Solved     = "solved"
	Timeout    = "timeout"
	Unsolvable = "unsolvable"
)

/**
 * Return a slice of all nodes with successor states to the current
 *
 * Note: All info other than the state is invalid as these states
 * are likely repeats and have already been calculated
 **/
func (n Node) getSuccessorStates(ignore_prev_moves bool) []Puzzle {
	var succStates = []Puzzle{}

	succStates = append(succStates, n.state.getSuccessors(ignore_prev_moves)...)

	return succStates
}

type Heuristic func(p Puzzle) float32

var heuristics = [...]Heuristic{h1, h2, h3, h4}

func h1(p Puzzle) float32 {
	var cost float32 = 0
	for i := 0; i < p.size(); i++ {
		if e := p.getN(i); (e != 0) && (e != i) {
			cost++
		}
	}

	return cost
}

func h2(p Puzzle) float32 {
	var cost float32 = 0
	for r := 0; r < p.len(); r++ {
		for c := 0; c < p.len(); c++ {
			goalPos := p.getGoalPos(p.get(RowCol{row: r, col: c}))
			cost += float32(math.Abs(float64(r - goalPos.row)))
			cost += float32(math.Abs(float64(c - goalPos.col)))
		}
	}
	return cost
}

func h3(p Puzzle) float32 {
	var cost float32 = 0
	var arr []int = make([]int, p.size())
	var idx int = 0
	for _, row := range p.arr {
		for _, e := range row {
			arr[idx] = e
			idx++
		}
	}

	for i := len(arr) - 1; i > 0; i-- {
		max := 0
		for j := 0; j <= i; j++ {
			if arr[max] < arr[j] {
				max = j
			}
		}
		if i != max {
			temp := arr[i]
			arr[i] = arr[max]
			arr[max] = temp
			cost++
		}
	}

	return cost
}

func h4(p Puzzle) float32 {
	var cost float32 = 0
	for r := 0; r < p.len(); r++ {
		for c := 0; c < p.len(); c++ {
			goalPos := p.getGoalPos(p.get(RowCol{row: r, col: c}))
			cost += euclidean_dist(goalPos.row-r, goalPos.col-c)
		}
	}
	return float32(cost)
}

func indexOf(list []*Node, state Puzzle) int {
	for i, e := range list {
		if e.state.equals(state) {
			return i
		}
	}
	return -1
}

func minNodeIdx(list []*Node) int {
	var minF float32 = list[0].getF()
	var idx int = 0
	for i, e := range list {
		test := e.getF()
		if test < minF {
			minF = test
			idx = i
		}
	}

	return idx
}

func popLowest(list []*Node) (*Node, []*Node) {
	var idx int = minNodeIdx(list)

	var minNode *Node = list[idx]

	copy(list[idx:], list[idx+1:])
	list[len(list)-1] = nil
	list = list[:len(list)-1]

	return minNode, list
}

func (n Node) isFinal() bool {
	return n.state.isSolved()
}

/**
 * return solution path, open list size, closed list size, avg branching factor
 * calc runtime outside of func
 **/
func a_star(initial Puzzle, h Heuristic, time_limit int, ignore_prev_moves bool) (status Status, path []Puzzle, openSize int, closedSize int) {
	start := time.Now()
	var openList = []*Node{{
		state: initial,
		g:     0,
		h:     h(initial),
	}} // frontier starts with the initial state

	var closedList = []*Node{} // explored is empty
	var cur *Node
	for len(openList) > 0 { // while there are nodes to explore
		// timeout condition
		if time_limit > 0 { // 0 or negative time limit is ignored
			if time.Since(start).Seconds() >= float64(time_limit) {
				return Timeout, make([]Puzzle, 0), len(openList), len(closedList)
			}
		}

		cur, openList = popLowest(openList)

		if cur.isFinal() { // found solution
			var steps int = 0
			var path []Puzzle = make([]Puzzle, 0)
			for node := cur; node != nil; node = node.prev {
				path = append(path, node.state.copy())
				steps++
			}

			return Solved, path, len(openList), len(closedList)

		} else { // still exploring
			closedList = append(closedList, cur)
			for _, state := range cur.getSuccessorStates(ignore_prev_moves) {
				nIdx := indexOf(openList, state)

				if nIdx != -1 { // state is in open list
					if cur.g+1 < openList[nIdx].g { // update with better path
						openList[nIdx].g = cur.g + 1
						openList[nIdx].prev = cur
					}
				} else if indexOf(closedList, state) != -1 { // state is in closed list
					continue
				} else { // state has not been seen yet
					openList = append(openList, &Node{
						state: state,
						g:     cur.g + 1,
						h:     h(state),
						prev:  cur,
					})
				}
			}
		}
	}
	return Unsolvable, make([]Puzzle, 0), len(openList), len(closedList)
}

func solve(initial Puzzle, heuristic_num int, time_limit int, ignore_prev_moves bool) {
	start := time.Now()
	status, path, openLen, closedLen := a_star(initial, heuristics[heuristic_num-1], time_limit, ignore_prev_moves)
	duration := time.Since(start)

	if config.Metrics.Status {
		if status == Solved {
			logger.Printf("Status: %v, Valid: %v\n", status, verifySolution(path))
		} else {
			logger.Printf("Status: %v\n", status)
		}
	}

	if config.Metrics.Execution_time {
		logger.Printf("Execution time: %.3fs\n", duration.Seconds())
	}

	if status == Solved && config.Metrics.Solution_length {
		logger.Printf("Solution Length: %v\n", len(path)-1)
	}

	if config.Metrics.Nodes_explored {
		logger.Printf("Nodes Explored: %v\n", openLen+closedLen)
	}

	if config.Metrics.Frontier_size {
		logger.Printf("Frontier Size: %v\n", openLen)
	}

	if config.Metrics.Nodes_evaluated {
		logger.Printf("Nodes Evaluated: %v\n", closedLen)
	}

	if status == Solved && config.Metrics.Solution_path {
		logger.Printf("Solution Path:\n")
		for i := len(path) - 1; i >= 0; i-- {
			logger.Print(path[i].toStr())
			logger.Println()
		}
	}
}

func verifySolution(path []Puzzle) bool {
	if len(path) == 0 {
		return false
	}

	if !path[0].isSolved() {
		return false
	}

	var cur, prev Puzzle
	for i := range path {
		if i >= len(path)-1 {
			break
		}
		cur = path[i]
		prev = path[i+1]
		if !cur.isSuccessorTo(prev) {
			return false
		}
	}

	return true
}
