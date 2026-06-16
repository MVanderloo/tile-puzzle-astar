package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const config_file = "config.json"

var config Config
var logger log.Logger
var logfile *os.File

func init() {

	if !ConfigExists() { // if there is no config file, create one and return
		fmt.Printf("Couldn't find config file (%v)\n", config_file)
		fmt.Printf("Creating %v\n", config_file)
		createConfig()
		os.Exit(0)
	}

	config = readConfig()

	if config.Random_seed == 0 { // set global random seed
		config.Random_seed = int64(time.Now().UnixMilli())
	}
}

func main() {
	// open logfile and print header
	if config.Log_file == "" {
		logfile = openLogFile(time.Now().Format("15-04-05") + ".txt")
	} else {
		logfile = openLogFile(config.Log_file + ".txt")
	}

	defer logfile.Close()
	logger.SetOutput(logfile)
	logger.Println(time.Now().Format("15:04:05 02/01/06"))
	logger.Printf("Random Seed: %v", config.Random_seed)
	logger.Print(logFileSpacer())

	for i, input := range config.Inputs {
		// get heuristics for this input
		var heuristics []int
		if input.Heuristics == nil {
			heuristics = config.Default_inputs.Heuristics
		} else {
			heuristics = input.Heuristics
		}

		// get time limit for this input
		var time_limit int
		if input.Time_limit <= 0 {
			time_limit = config.Default_inputs.Time_limit
		} else {
			time_limit = input.Time_limit
		}

		// for each heuristic
		for _, heuristic_num := range heuristics {
			logger.Printf("Puzzle: %v-%v", i+1, heuristic_num)

			var p Puzzle // find what type of input was specified
			var swaps int
			if input.Swaps != 0 {
				p = newPuzzleSwapped(input.Size, input.Swaps)
				swaps = input.Swaps
			} else if input.Misplaced != 0 {
				p, swaps = newPuzzleMisplaced(input.Size, input.Misplaced)
			} else {
				p = newPuzzleSolved(input.Size)
				swaps = 0
			}
			logger.Printf("Size: %v\n", p.len())

			if config.Metrics.Initial_state {
				logger.Printf("Initial: \n%v", p.toStr())
			}

			logger.Printf("Initial Misplaced Tiles: %v / %v\n", h1(p), p.size()-1)
			logger.Printf("Swaps Used to Generate: %v\n", swaps)
			if input.Use_prev_move {
				logger.Printf("Using prev node in successor generation\n")
			}

			switch heuristic_num {
			case 1:
				logger.Print("Heuristic: Number of Misplaced (1)")
			case 2:
				logger.Print("Heuristic: Manhattan Distance (2)")
			case 3:
				logger.Print("Heuristic: Maxsort Swaps (3)")
			case 4:
				logger.Print("Heuristic: Euclidian Distance (4)")
			default:
				logger.Printf("Heuristic: Unknown (%v)", heuristic_num)
			}
			logger.Print("\n")

			solve(p, heuristic_num, time_limit, !input.Use_prev_move)
			logger.Print(logFileSpacer())
		}
	}
}
