package main

import (
	"fmt"
	"runtime"
	"time"
)

type Book struct {
	NumPages int
}

type Cart struct {
	Books       []Book
	Cap         int
	ReadyToLoad bool
}

func main() {
	runtime.GOMAXPROCS(1)
	start := time.Now()

	pileOfBooks := MakePileOfBooks(500)
	booksMoved := 0
	booksBurned := 0
	cart := Cart{
		Books:       []Book{},
		ReadyToLoad: true,
		Cap:         15,
	}

	for len(pileOfBooks) > 0 {
		fmt.Printf("%10s %5d Books remain, %5d books moved, %5d books burned\n", "loading", len(pileOfBooks), booksMoved, booksBurned)
		pileOfBooks = LoadCart(pileOfBooks, &cart)
		fmt.Printf("%10s %5d Books remain, %5d books moved, %5d books burned\n", "moving", len(pileOfBooks), booksMoved, booksBurned)
		booksMoved += len(cart.Books)
		MoveCart(&cart)
		fmt.Printf("%10s %5d Books remain, %5d books moved, %5d books burned\n", "burning", len(pileOfBooks), booksMoved, booksBurned)
		booksBurned += len(cart.Books)
		BurnBooks(&cart)
		fmt.Printf("%10s %5d Books remain, %5d books moved, %5d books burned\n", "moving", len(pileOfBooks), booksMoved, booksBurned)
		MoveCart(&cart)
		// fmt.Println(len(pileOfBooks), "Books remain")
	}

	fmt.Println("Done:", time.Since(start))
}

func MakePileOfBooks(n int) []Book {
	books := make([]Book, n)
	for i := 0; i < n; i++ {
		books[i].NumPages = 1
	}
	return books
}

func BurnBooks(c *Cart) {
	for _, b := range c.Books {
		time.Sleep(time.Duration(b.NumPages) * time.Millisecond)
		// fmt.Println("burned book")
	}
	c.Books = []Book{}
	return
}

func LoadCart(pile []Book, c *Cart) []Book {
	// wait for cart to be ready
	for ready := c.ReadyToLoad; !ready; {
		time.Sleep(1 * time.Millisecond)
	}
	// wait for books to be loaded
	for _, b := range pile {
		time.Sleep(time.Duration(b.NumPages) * time.Millisecond)
		c.Books = append(c.Books, b)
		if len(c.Books) >= c.Cap {
			break
		}
	}
	//remove books from pile
	loaded := len(c.Books)
	newPile := MakePileOfBooks(len(pile) - loaded)
	return newPile
}

func MoveCart(c *Cart) {
	time.Sleep(10 * time.Millisecond)
	c.ReadyToLoad = !c.ReadyToLoad
}
