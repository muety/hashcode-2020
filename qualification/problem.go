package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Problem struct {
	Days      int
	Libraries map[int]*Library
}

type Signup struct {
	Day     int
	Library *Library
	Books   *Books
}

type Solution []*Signup

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
			libraries[nLibs].Finalize()
			nLibs++
		}
	}

	return &Problem{Days: nDays, Libraries: libraries}
}

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

func (s *Signup) String() string {
	var sb strings.Builder
	if len(*s.Books) > 0 {
		sb.WriteString(fmt.Sprintf("%d %d\n", s.Library.Id, len(*s.Books)))
	}
	for _, b := range *s.Books {
		sb.WriteString(fmt.Sprintf("%d ", b.Id))
	}
	return strings.TrimSpace(sb.String()) + "\n"
}
