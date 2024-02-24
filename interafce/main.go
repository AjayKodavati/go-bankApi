package main

import (
	"fmt"
	"log"
	"strconv"
)

type Stringer interface {
	String() string
}

type Book struct {
	Title  string
	Author string
}

// This book type satisfies the interface because it has s String() string method, it is not really important what is this Book type or what it does
// The only thing that matters is that is has a method called String() which returns a string value
func (b *Book) String() string {
	return fmt.Sprintf("Book: %s - %s", b.Title, b.Author)
}

type Count int

func (c Count) String() string {
	return strconv.Itoa(int(c))
}

// This function which takes any object that satisfies the fmt.Stringer interface as parameter
func WriteLog(s fmt.Stringer) {
	log.Print(s.String())
}

func main() {
	book := &Book{
		"Alice in wonderland", "lewis carrol",
	}
	WriteLog(book)

	count := Count(3)
	WriteLog(count)
}