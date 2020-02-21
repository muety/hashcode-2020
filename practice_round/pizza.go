package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput() (int, []int) {
	var M int
	var N int
	var types []int

	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; ; i++ {
		scanner.Scan()
		line := scanner.Text()
		if line == "" {
			break
		}

		if i == 0 {
			numbers := strings.Split(line, " ")
			M, _ = strconv.Atoi(numbers[0])
			N, _ = strconv.Atoi(numbers[1])

			types = make([]int, N)
		} else {
			numbers := strings.Split(line, " ")
			for j := 0; j < len(numbers); j++ {
				nSlices, _ := strconv.Atoi(numbers[j])
				types[j] = nSlices
			}
		}
	}

	return M, types
}

func printOutput(data []bool) {
	var sb strings.Builder
	var c int

	for i := 0; i < len(data); i++ {
		if data[i] {
			sb.WriteString(strconv.Itoa(i) + " ")
			c++
		}
	}

	fmt.Println(c)
	fmt.Println(strings.TrimSpace(sb.String()))
}

func solveGreedy(M int, types []int) ([]bool, int) {
	solution := make([]bool, len(types))

	var sum int
	for i := len(types) - 1; i >= 0; i-- {
		if sum+types[i] <= M {
			sum += types[i]
			solution[i] = true
		}
	}

	return solution, sum
}

func main() {
	M, types := readInput()
	solution, sum := solveGreedy(M, types)

	score := (float64(sum) / float64(M)) * 100
	fmt.Fprintf(os.Stderr, "Total Slices: %d, Score: %.2f %%\n", sum, score)

	printOutput(solution)
}
