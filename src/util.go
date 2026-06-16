package main

import (
	"math"
)

/**
* This file contains some utilities used by the program including...
* the implementation of RowCol struct used for indexing a 2d board without confusing axis.
* directions Up, Down, Left, Right and None to make sense in context of the puzzle
**/

type RowCol struct {
	row int
	col int
}

/**
 * Enum for cardinal directions
 */
type Move int

const (
	Down Move = iota - 2
	Left
	None
	Right
	Up
)

/**
 * Returns the opposite of the inputted move
 **/
func (m Move) opposite() Move {
	switch m {
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	case Up:
		return Down
	default:
		return None
	}
}

/**
 * Returns the index of the 2d array as if it was a single array
 **/
func (rc RowCol) toN(len int) int {
	return len*rc.row + rc.col
}

/**
 * input validation to ensure a puzzle only has unique digits from 0 to n-1
 * assumes isSquare has been called so does not check for a square number range
 **/
func containsAllIndices(arr []int) bool {
	var seen []bool = make([]bool, len(arr))
	for _, e := range arr {
		if e < 0 || e >= len(arr) {
			return false
		}

		if seen[e] {
			return false
		} else {
			seen[e] = true
		}
	}
	return true
}

/*
 * return if the number is square and what it's integer base is
 */
func isSquare(n int) (bool, int) {
	base, rem := math.Modf(math.Sqrt(float64(n)))
	return rem == 0, int(base)
}

/**
 * remove an element from a slice (doesn't retain ordering)
 */
func remove(s []int, i int) []int {
	if i < 0 || i >= len(s) {
		panic("index out of bounds for remove")
	}

	if len(s) == 0 || len(s) == 1 {
		return []int{}
	}

	if len(s) == 2 {
		if i == 0 {
			s[0] = s[1]
		}
		return s[:1]
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

var euclid_lookup = make(map[RowCol]float32, 9)

func euclidean_dist(x int, y int) float32 {
	if x < 0 {
		x = -x
	}

	if y < 0 {
		y = -y
	}

	if x > y {
		temp := x
		x = y
		y = temp
	}

	if val, ok := euclid_lookup[RowCol{x, y}]; ok {
		return val
	} else {
		var dist float32 = float32(math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2)))
		euclid_lookup[RowCol{x, y}] = dist
		return dist
	}
}
