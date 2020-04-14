package main

import (
	"fmt"
	"sync"
	"time"
	"runtime"
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

	pileOfBooks := MakePileOfBooks(250)
	cart := Cart{
		Books:       []Book{},
		ReadyToLoad: true,
		Cap:         15,
	}

	incineratorChan := make(chan Book)
	go func() {
		defer close(incineratorChan)
		for b := range incineratorChan {
			time.Sleep(time.Duration(b.NumPages) * time.Millisecond)
		}
	}()

	pileLoaderChan := make(chan *Cart)
	moverFullChan := make(chan *Cart)
	moverEmptyChan := make(chan *Cart)
	incineratorLoaderChan := make(chan *Cart)

	go func() {
		for c := range pileLoaderChan {
			fmt.Printf("%5d Books remain\n", len(pileOfBooks))
			pileOfBooks = LoadCart(pileOfBooks, c)
			moverFullChan <- c
		}
	}()

	go func() {
		for c := range moverFullChan {
			// fmt.Println("moving cart with", len(c.Books), "books")
			MoveCart(c)
			incineratorLoaderChan <- c
		}
	}()

	go func() {
		for c := range moverEmptyChan {
			// fmt.Println("moving cart with", len(c.Books), "books")
			MoveCart(c)
			pileLoaderChan <- c
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for c := range incineratorLoaderChan {
			// fmt.Println("burning")
			BurnBooks(c, incineratorChan)
			if len(pileOfBooks) < 1 {
				wg.Done()
				return
			}
			moverEmptyChan <- c
		}
	}()

	pileLoaderChan <- &cart

	pileOfBooks2 := MakePileOfBooks(250)
	cart2 := Cart{
		Books:       []Book{},
		ReadyToLoad: true,
		Cap:         15,
	}

	pileLoaderChan2 := make(chan *Cart)
	moverFullChan2 := make(chan *Cart)
	moverEmptyChan2 := make(chan *Cart)
	incineratorLoaderChan2 := make(chan *Cart)

	go func() {
		for c := range pileLoaderChan2 {
			fmt.Printf("%5d Books remain\n", len(pileOfBooks2))
			pileOfBooks2 = LoadCart(pileOfBooks2, c)
			moverFullChan2 <- c
		}
	}()

	go func() {
		for c := range moverFullChan2 {
			// fmt.Println("moving cart with", len(c.Books), "books")
			MoveCart(c)
			incineratorLoaderChan2 <- c
		}
	}()

	go func() {
		for c := range moverEmptyChan2 {
			// fmt.Println("moving cart with", len(c.Books), "books")
			MoveCart(c)
			pileLoaderChan2 <- c
		}
	}()

	go func() {
		for c := range incineratorLoaderChan2 {
			// fmt.Println("burning")
			BurnBooks(c, incineratorChan)
			moverEmptyChan2 <- c
		}
	}()

	pileLoaderChan2 <- &cart2

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
