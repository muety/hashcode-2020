package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Problem struct {
	Days      int
	Libraries map[int]*Library
}

type Solver struct {
	Problem         *Problem
	SeedStrategy    string
	PickStrategy    string
	SortEvery       int
	days            int
	lastSorted      int
	libraries       map[int]*Library
	blacklist       map[int]bool
	solution        *Solution
	sortedLibraries []*Library
}

type Books []*Book

type Book struct {
	Id    int
	Score int
}

type Library struct {
	Id         int
	SignupDays int
	ShipRate   int
	Books      Books
}

type Signup struct {
	Day     int
	Library *Library
	Books   *Books
}

func (s *Signup) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", s.Library.Id, len(*s.Books)))
	for _, b := range *s.Books {
		sb.WriteString(fmt.Sprintf("%d ", b.Id))
	}
	return strings.TrimSpace(sb.String()) + "\n"
}

type Solution []*Signup

func (s *Solution) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(*s)))
	for _, signup := range *s {
		sb.WriteString(signup.String())
	}
	return sb.String()
}

func (s *Solution) Dump(file string) {
	if err := ioutil.WriteFile(file, []byte(s.String()), os.ModePerm); err != nil {
		panic(err)
	}
}

func (s *Solver) SolutionScore() (int, int, float64) {
	var visitedBooks map[int]bool

	var maxScore int
	visitedBooks = make(map[int]bool)
	for _, l := range s.Problem.Libraries {
		for _, b := range l.Books {
			if _, ok := visitedBooks[b.Id]; !ok {
				maxScore += b.Score
				visitedBooks[b.Id] = true
			}
		}
	}

	var totalScore int
	visitedBooks = make(map[int]bool)
	for _, signup := range *s.solution {
		for _, b := range *signup.Books {
			if _, ok := visitedBooks[b.Id]; !ok {
				totalScore += b.Score
				visitedBooks[b.Id] = true
			}
		}
	}

	var accuracy = (float64(totalScore) / float64(maxScore)) * 100

	return totalScore, maxScore, accuracy
}

func (s *Solver) SolutionStats() string {
	if s.solution == nil {
		panic("not solved yet")
	}

	totalScore, maxScore, accuracy := s.SolutionScore()
	return fmt.Sprintf("Score: %d / %d (%.2f %%), Remaining Days: %d", totalScore, maxScore, accuracy, s.days)
}

func (s *Solver) Run(n int) Solution {
	var bestSolution Solution
	var bestScore float64

	for i := 0; i < n; i++ {
		solution := s.Solve()
		_, _, score := s.SolutionScore()
		if score > bestScore {
			bestSolution = solution
			bestScore = score
		}
		fmt.Printf("Ran iteration %d / %d. Best score is %.2f %%\n", i+1, n, bestScore)
	}

	return bestSolution
}

func (s *Solver) Solve() Solution {
	s.init()

	solution := make(Solution, 0)

	var firstPick func() *Library
	var pick func() *Library

	if s.SeedStrategy == "best" {
		firstPick = s.pickBest
	} else if s.SeedStrategy == "random" {
		firstPick = s.pickRandom
	}

	if s.PickStrategy == "best" {
		pick = s.pickBest
	} else if s.PickStrategy == "random" {
		pick = s.pickRandom
	}

	for lib := firstPick(); s.days > 0 && lib != nil; lib = pick() {
		books := lib.Take(s.blacklist, s.days)
		if books == nil {
			continue
		}

		solution = append(solution, &Signup{Library: lib, Day: s.Problem.Days - s.days, Books: books})

		s.updateDays(lib)
		s.updateBlacklist(books)
	}

	s.solution = &solution
	return solution
}

func (s *Solver) init() {
	s.days = s.Problem.Days
	s.blacklist = make(map[int]bool)
	s.libraries = make(map[int]*Library)
	for k, v := range s.Problem.Libraries {
		s.libraries[k] = v
	}
}

func (s *Solver) updateBlacklist(books *Books) {
	for _, b := range *books {
		s.blacklist[b.Id] = true
	}
}

func (s *Solver) updateDays(lib *Library) {
	s.days = s.days - lib.SignupDays
}

func (s *Solver) resort() {
	s.sortedLibraries = s.sortLibraries()
	s.lastSorted = 0
}

func (s *Solver) pickBest() *Library {
	if len(s.libraries) <= s.lastSorted {
		return nil
	}
	if s.sortedLibraries == nil || s.lastSorted >= s.SortEvery-1 {
		s.resort()
	}
	s.lastSorted++
	delete(s.libraries, s.sortedLibraries[s.lastSorted].Id)
	return s.sortedLibraries[s.lastSorted]
}

func (s *Solver) pickRandom() *Library {
	if len(s.libraries) == 0 {
		return nil
	}

	var i int
	keys := make([]int, len(s.libraries))
	for k := range s.libraries {
		keys[i] = k
		i++
	}

	k := int(rand.Int31n(int32(len(s.libraries))))
	lib := s.libraries[keys[k]]
	delete(s.libraries, lib.Id)
	return lib
}

func (s *Solver) sortLibraries() []*Library {
	sortedLibraries := make([]*Library, len(s.libraries))

	var i int
	for _, v := range s.libraries {
		sortedLibraries[i] = v
		i++
	}

	sort.Slice(sortedLibraries, func(i, j int) bool {
		return sortedLibraries[i].Score(s.blacklist, s.days) > sortedLibraries[j].Score(s.blacklist, s.days)
	})

	return sortedLibraries
}

func (l *Library) SortBooks() {
	sort.Slice(l.Books, func(i, j int) bool {
		return l.Books[i].Score > l.Books[j].Score
	})
}

func (l *Library) Take(blacklist map[int]bool, remainingDays int) *Books {
	nMax := (remainingDays - l.SignupDays) * l.ShipRate
	if nMax < 0 {
		return nil
	}

	booksFiltered := make(Books, 0, len(l.Books))
	for _, b := range l.Books {
		if _, ok := blacklist[b.Id]; !ok {
			booksFiltered = append(booksFiltered, b)
		}
	}
	booksSubset := booksFiltered[0:intMin(len(booksFiltered), nMax)]

	return &booksSubset
}

func (l *Library) Score(blacklist map[int]bool, remainingDays int) float64 {
	return scoreSelection(l.Take(blacklist, remainingDays), l, remainingDays)
}

func scoreSelection(bookSelection *Books, library *Library, remainingDays int) float64 {
	var bookScores int

	if bookSelection == nil {
		return 0
	}

	for _, b := range *bookSelection {
		bookScores += b.Score
	}

	averageDailyScore := (float64(bookScores) / float64(len(*bookSelection))) * float64(library.ShipRate)
	return averageDailyScore / float64(library.SignupDays)
}

func ReadProblem(file string) *Problem {
	const maxCapacity = 512 * 1024

	nDays := 0
	nLibs := 0
	books := make(map[int]*Book)
	libraries := make(map[int]*Library)

	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := make([]byte, maxCapacity)
	scanner := bufio.NewScanner(f)
	scanner.Buffer(buf, maxCapacity)

	for i := 0; ; i++ {
		scanner.Scan()
		line := scanner.Text()
		if line == "" {
			break
		}
		parts := strings.Split(line, " ")

		if i == 0 {
			// Meta data
			nDays, _ = strconv.Atoi(parts[2])
		} else if i == 1 {
			// Book scores
			for j := 0; j < len(parts); j++ {
				s, _ := strconv.Atoi(parts[j])
				books[j] = &Book{Id: j, Score: s}
			}
		} else if i > 1 && i%2 == 0 {
			// Library meta data
			n, _ := strconv.Atoi(parts[0])
			s, _ := strconv.Atoi(parts[1])
			r, _ := strconv.Atoi(parts[2])
			libraries[nLibs] = &Library{Id: nLibs, SignupDays: s, ShipRate: r, Books: make([]*Book, n)}
		} else if i > 0 && i%2 == 1 {
			// Library books
			for j := 0; j < len(parts); j++ {
				bid, _ := strconv.Atoi(parts[j])
				if b, ok := books[bid]; ok {
					libraries[nLibs].Books[j] = b
				} else {
					panic(fmt.Sprintf("couldn't find book %d", bid))
				}
			}
			libraries[nLibs].SortBooks()
			nLibs++
		}
	}

	return &Problem{Days: nDays, Libraries: libraries}
}

func intMin(x, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	inFile := os.Args[1]
	outFile := "submission-" + strings.Split(inFile, "/")[1]

	problem := ReadProblem(inFile)
	solver := Solver{Problem: problem, SeedStrategy: "best", PickStrategy: "best", SortEvery: 10}
	solution := solver.Run(1)

	fmt.Println(solver.SolutionStats())

	solution.Dump(outFile)
}
