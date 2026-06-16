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
