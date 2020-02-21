package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())

	if len(os.Args) < 5 {
		fmt.Printf("Not enough arguments given\nExample: go run *.go data/a_example.txt 10 random best\n")
		os.Exit(1)
	}
}

func main() {
	inFile := os.Args[1]
	outFile := "submission-" + strings.Split(inFile, "/")[1]

	nIterations, _ := strconv.Atoi(os.Args[2])
	seedStrategy := os.Args[3] // "best" or "random"
	pickStrategy := os.Args[4] // "best" or "random"

	problem := ReadProblem(inFile)
	solver := Solver{Problem: problem, SeedStrategy: seedStrategy, PickStrategy: pickStrategy}
	solution := solver.Run(intMax(1, nIterations))

	fmt.Println(solver.SolutionStats())

	solution.Dump(outFile)
}
