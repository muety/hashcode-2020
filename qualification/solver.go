package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Solver struct {
	Problem         *Problem
	SeedStrategy    string
	PickStrategy    string
	days            int
	lastPicked      int
	libraries       map[int]*Library
	sortedLibraries []*Library
	blacklist       map[int]bool
	solution        *Solution
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

		if len(*books) > 0 {
			solution = append(solution, &Signup{Library: lib, Day: s.Problem.Days - s.days, Books: books})
		}

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
	s.lastPicked = 0
}

func (s *Solver) pickBest() *Library {
	if len(s.libraries) <= s.lastPicked {
		return nil
	}
	if s.sortedLibraries == nil {
		s.resort()
	} else {
		s.lastPicked++
	}
	delete(s.libraries, s.sortedLibraries[s.lastPicked].Id)
	return s.sortedLibraries[s.lastPicked]
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
		return sortedLibraries[i].Score > sortedLibraries[j].Score
	})

	return sortedLibraries
}
