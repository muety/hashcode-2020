package main

import "sort"

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
	BookScore  int
}

func (l *Library) Finalize() {
	// Evaluate total book score
	for _, b := range l.Books {
		l.BookScore += b.Score
	}

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

func (l *Library) Score(blacklist map[int]bool, remainingDays int) float64 {
	var bookScores int

	bookSelection := l.Take(blacklist, remainingDays)
	if bookSelection == nil {
		return 0
	}

	for _, b := range *bookSelection {
		bookScores += b.Score
	}

	averageDailyScore := (float64(bookScores) / float64(len(*bookSelection))) * float64(l.ShipRate)
	return averageDailyScore / float64(l.SignupDays)
}
