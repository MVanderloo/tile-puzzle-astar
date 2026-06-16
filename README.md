# Tile Puzzle AI

### CS 420 Project 1

### Michael Vanderloo

This project was done in golang. You will need to install golang and update GOPATH to include this directory.

## Install Golang

https://go.dev/doc/install

## Run program

This program is built and run using the go command line tool.
You can use

```code
go run .
```

to run the program without generating files

or

```code
go build .
```

to create an executable.

## Configs

The input to the experiment is config.json located in the project folder. If the config is not present, the program will generate one and terminate. The config generated should give a decent overview of the capabilities of the program

The config.json must be valid json. Extra values in the config will be ignored and missing values will be have a defualt value.

There are a few behaviors that the program will do with different config values. Missing values will be assigned the golang zero-value for it's datatype. In short numbers are 0, strings are empty, and slices are nil. Most of this behavior is accounted for in the program but I didn't test super thouroughly

Random seed = 0 will generate a random seed using the system time

Time limit \<= 0 will have no time limit

Both the time limit and heuristics have a default value that can be omitted in inputs to use, or overridden.

Heuristics is a list of values 1-4

Metrics are a list of metrics that may or may not be useful in analyzing the algorithm performance. Set one to false or delete it to disclude it from the log file. By default all metrics are enabled except solution path, as it prints every state along the solution path and makes the file harder to read, but can be used to verify a solution.

## Logging

Everything is printed to the logfile specified in config.json. The program will append a .txt extension to the provided filename. It will overwrite a file so be careful with this. If no file is provided, it will write to a file with current time like this: 15-04-05.txt

The naming of the Puzzles is of the format 1-3. The first number denotes a given input, and the second differentiates that puzzle with different heuristics. This is so it's easier to compare the same puzzle with different heuristics.

## Implementation Details

The puzzle is implemented as a 2d array. I wish I had done a 1d array because it would be less pointers but the overhead is minimal anyways.

One optimization is that the puzzle stores the location of 0, which is used to shortcut certain functions like finding successor states.

An optimization of A\* is that the node holds a puzzle and the previous move as an integer. In generating the successor states, the algorithm doesn't consider the successor that undoes the previous move, as this node is guarenteed to be evaluated.
You can see the results of this optimization by setting
"use prev move": true
per input in the config. The generated config shows 1 example of this on the same puzzle. For the seed that I used, this cut the execution time from 1.5s to 1s, which is pretty significant. Though the number of nodes evaluated is the exact same, it removes a linear search of the closed list and open list for every node evaluated.

The solution verifier just checks that the solution path ends in a solved state, and verifies that every node is a successor of the previous. I have not encountered a solution that was not verified to be incorrect by the algorithm.

The euclidian heuristic uses a lookup table, as repeat values are extremely common. It uses a map[int, int] -> float. All others are calculated as it's cheap enough.

## Performance

Go is a very performant language. Everything was designed to be as lightweight as possible and I think I achieved that. If you check the file archive_log.txt, it managed to evaluate a size 4 puzzle with all nodes out of place. The solution length is 38 and there were 254977 nodes kept in memory. This did take about an hour to run.

Of the heuristics, h2 was the best. Though h1 is cheaper to compute, h2 does a better job at evaluating which move to do next. h4 is similar performance to h2 but doesn't do as good a job as h2 as there are more nodes explored to get to the solution. The heuristics in order of speed (on 1 particular puzzle) is: 2, 4, 1, 3

## Next Steps

- implement faster lookup of open list nodes using a priority queue
- implement faster lookup of closed list nodes
- implement puzzle as 1d array
