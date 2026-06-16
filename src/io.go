package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Log_file    string `json:"log file"`
	Random_seed int64  `json:"random seed"`
	Metrics     struct {
		Initial_state       bool `json:"initial state"`
		Num_misplaced_tiles bool `json:"num misplaced tiles"`
		Max_solution_length bool `json:"max solution length"`
		Status              bool `json:"status"`
		Execution_time      bool `json:"execution time"`
		Solution_length     bool `json:"solution length"`
		Nodes_explored      bool `json:"nodes explored"`
		Frontier_size       bool `json:"frontier size"`
		Nodes_evaluated     bool `json:"nodes evaluated"`
		Solution_path       bool `json:"solution path"`
	} `json:"metrics"`
	Default_inputs struct {
		Heuristics []int `json:"heuristics"`
		Time_limit int   `json:"time limit"`
	} `json:"default inputs"`
	Inputs []struct {
		Size          int   `json:"size"`
		Swaps         int   `json:"swaps"`
		Misplaced     int   `json:"misplaced"`
		Heuristics    []int `json:"heuristics"`
		Time_limit    int   `json:"time limit"`
		Use_prev_move bool  `json:"use prev move"`
	} `json:"inputs"`
}

func ConfigExists() bool {
	if _, err := os.Stat(config_file); err == nil {
		return true
	} else {
		return false
	}
}

func readConfig() Config {
	// Open our jsonFile
	jsonFile, err := os.Open(config_file)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err2 := io.ReadAll(jsonFile)
	if err2 != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var config Config
	if err3 := json.Unmarshal(byteValue, &config); err3 != nil {
		fmt.Printf("Invalid config file. Delete %v and run program to generate a new one\n", config_file)
		os.Exit(1)
	}

	for _, input := range config.Inputs {
		if (input.Misplaced != 0) && (input.Swaps != 0) {
			panic("cannot specify both swaps and misplaced in config.inputs")
		}
	}

	return config
}

func createConfig() {
	var contents = []string{
		"{",
		"\t\"log file\": \"log\",",
		"\t\"random seed\": 0,",
		"\t\"metrics\": {",
		"\t\t\"initial state\": true,",
		"\t\t\"num misplaced tiles\": true,",
		"\t\t\"max solution length\": true,",
		"\t\t\"status\": true,",
		"\t\t\"execution time\": true,",
		"\t\t\"solution length\": true,",
		"\t\t\"nodes explored\": true,",
		"\t\t\"frontier size\": true,",
		"\t\t\"nodes evaluated\": true,",
		"\t\t\"solution path\": false",
		"\t},",
		"\t\"default inputs\": {",
		"\t\t\"heuristics\": [2],",
		"\t\t\"time limit\": 60",
		"\t},",
		"\t\"inputs\": [",
		"\t\t{",
		"\t\t\t\"size\": 2,",
		"\t\t\t\"misplaced\": 3",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 9,",
		"\t\t\t\"swaps\": 40",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 3,",
		"\t\t\t\"misplaced\": 9,",
		"\t\t\t\"heuristics\": [1, 2, 3, 4]",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 3,",
		"\t\t\t\"swaps\": 20",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 4,",
		"\t\t\t\"swaps\": 20,",
		"\t\t\t\"heuristics\": [1, 2, 3, 4]",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 5,",
		"\t\t\t\"misplaced\": 10,",
		"\t\t\t\"time limit\": 100",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 4,",
		"\t\t\t\"misplaced\": 16,",
		"\t\t\t\"time limit\": 0",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 4,",
		"\t\t\t\"misplaced\": 25,",
		"\t\t\t\"time limit\": 0",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 4,",
		"\t\t\t\"misplaced\": 25,",
		"\t\t\t\"time limit\": 0,",
		"\t\t\t\"use prev move\": true",
		"\t\t},",
		"\t\t{",
		"\t\t\t\"size\": 5,",
		"\t\t\t\"misplaced\": 15,",
		"\t\t\t\"time limit\": 60",
		"\t\t}",
		"\t]",
		"}",
	}
	f, err := os.Create(config_file)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, line := range contents {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Config created at %v\n", config_file)
}

func openLogFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		panic(err)
	}

	return f
}

func logFileSpacer() string {
	return "\n+----------------------------------+\n\n"
}
