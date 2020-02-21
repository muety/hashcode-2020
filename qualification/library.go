package main

import (
	"sort"
)

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
	Score      float64
}

func (l *Library) Finalize() {
	// Evaluate total book score
	var bookScore int
	for _, b := range l.Books {
		bookScore += b.Score
	}
	l.Score = float64(bookScore) / float64(l.SignupDays)

	// Sort books by score
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
