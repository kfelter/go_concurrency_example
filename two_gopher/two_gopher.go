package main

import (
	"fmt"
	"runtime"
	"sync"
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

	wg := sync.WaitGroup{}

	incineratorChan := make(chan Book)
	go func() {
		for b := range incineratorChan {
			time.Sleep(time.Duration(b.NumPages) * time.Millisecond)
		}
	}()

	pileOfBooks1 := MakePileOfBooks(250)
	cart1 := Cart{
		Books:       []Book{},
		ReadyToLoad: true,
		Cap:         15,
	}

	wg.Add(1)
	go func() {
		for len(pileOfBooks1) > 0 {
			fmt.Printf("%5d Books remain\n", len(pileOfBooks1))
			pileOfBooks1 = LoadCart(pileOfBooks1, &cart1)
			// fmt.Println("moving cart to incinerator with", cart1.Cap, "books")
			MoveCart(&cart1)
			// fmt.Println("burning books")
			BurnBooks(&cart1, incineratorChan)
			// fmt.Println("moving cart to book pile")
			MoveCart(&cart1)
			// fmt.Println(len(pileOfBooks1), "Books remain")
		}
		wg.Done()
	}()

	pileOfBooks2 := MakePileOfBooks(250)
	cart2 := Cart{
		Books:       []Book{},
		ReadyToLoad: true,
		Cap:         15,
	}

	wg.Add(1)
	go func() {
		for len(pileOfBooks2) > 0 {
			fmt.Printf("%5d Books remain\n", len(pileOfBooks2))
			pileOfBooks2 = LoadCart(pileOfBooks2, &cart2)
			// fmt.Println("moving cart to incinerator with", cart2.Cap, "books")
			MoveCart(&cart2)
			// fmt.Println("burning")
			BurnBooks(&cart2, incineratorChan)
			// fmt.Println("moving cart to book pile")
			MoveCart(&cart2)
			// fmt.Println(len(pileOfBooks2), "Books remain")
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Done:", time.Since(start))
}

func MakePileOfBooks(n int) []Book {
	books := make([]Book, n)
	for i := 0; i < n; i++ {
		books[i].NumPages = 1
	}
	return books
}

func BurnBooks(c *Cart, incinerator chan Book) {
	for _, b := range c.Books {
		incinerator <- b
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
