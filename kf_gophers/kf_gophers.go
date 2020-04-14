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

const (
	NUM_BOOKS          = 500
	NUM_CARTS          = 10
	NUM_LOADERS        = 10
	NUM_MOVERS         = 10
	NUM_BURNERS        = 1
	LOAD_TIME_PER_PAGE = 1 * time.Millisecond
	MOVE_TIME          = 10 * time.Millisecond
	BURN_TIME_PER_PAGE = 1 * time.Millisecond
	MAX_PROCS          = 1
)

func main() {
	runtime.GOMAXPROCS(MAX_PROCS)
	start := time.Now()

	pileOfBooks := MakePileOfBooks(NUM_BOOKS)
	booksMoved := 0
	booksBurned := 0
	doneChan := make(chan bool)

	incineratorChan := make(chan Book)
	go func() {
		defer close(incineratorChan)
		for b := range incineratorChan {
			time.Sleep(time.Duration(b.NumPages) * BURN_TIME_PER_PAGE)
			booksBurned += 1
			if NUM_BOOKS == booksBurned {
				doneChan <- true
			}
		}
	}()

	type CartAndPile struct {
		Pile []Book
		Cart *Cart
	}

	pileParserChan := make(chan *Cart, NUM_CARTS)
	pileLoaderChan := make(chan CartAndPile, NUM_CARTS)
	moverChan := make(chan *Cart, NUM_CARTS)
	incineratorLoaderChan := make(chan *Cart, NUM_CARTS)

	go func() {
		for c := range pileParserChan {
			fmt.Printf("%10s %5d Books remain, %5d books moved, %5d books burned\n", "parsing", len(pileOfBooks), booksMoved, booksBurned)
			var pile []Book
			if len(pileOfBooks) < 1 {
				return
			}
			if len(pileOfBooks) > 15 {
				pile, pileOfBooks = pileOfBooks[:15], pileOfBooks[15:]
			} else {
				pile, pileOfBooks = pileOfBooks, []Book{}
			}
			pileLoaderChan <- CartAndPile{Pile: pile, Cart: c}
		}
	}()

	for i := 0; i < NUM_LOADERS; i++ {
		go func() {
			for cartAndPile := range pileLoaderChan {
				LoadCart(cartAndPile.Pile, cartAndPile.Cart)
				fmt.Printf("%10s %5d Books remain, %5d books moved, %5d books burned\n", "loading", len(pileOfBooks), booksMoved, booksBurned)
				moverChan <- cartAndPile.Cart
			}
		}()
	}

	for i := 0; i < NUM_MOVERS; i++ {
		go func() {
			for c := range moverChan {
				// fmt.Println("moving cart with", len(c.Books), "books")
				MoveCart(c)
				booksMoved += len(c.Books)
				fmt.Printf("%10s %5d Books remain, %5d books moved, %5d books burned\n", "moving", len(pileOfBooks), booksMoved, booksBurned)
				if c.ReadyToLoad {
					pileParserChan <- c
				} else {
					incineratorLoaderChan <- c
				}
			}
		}()
	}

	for i := 0; i < NUM_BURNERS; i++ {
		go func() {
			for c := range incineratorLoaderChan {
				// fmt.Println("burning")
				BurnBooks(c, incineratorChan)
				moverChan <- c
			}
		}()
	}

	for i := 0; i < NUM_CARTS; i++ {
		pileParserChan <- &Cart{
			Books:       []Book{},
			ReadyToLoad: true,
			Cap:         15,
		}
	}

	done := <-doneChan
	if done {
		fmt.Printf("%10s %5d Books remain, %5d books moved, %5d books burned\n", "done", len(pileOfBooks), booksMoved, booksBurned)
		fmt.Println("Done:", time.Since(start))
	}
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

func LoadCart(pile []Book, c *Cart) {

	if len(pile) == 0 {
		c.Books = []Book{}
		return
	}
	// wait for books to be loaded
	for _, b := range pile {
		time.Sleep(time.Duration(b.NumPages) * LOAD_TIME_PER_PAGE)
		c.Books = append(c.Books, b)
		if len(c.Books) >= c.Cap {
			break
		}
	}
	return
}

func MoveCart(c *Cart) {
	time.Sleep(MOVE_TIME)
	c.ReadyToLoad = !c.ReadyToLoad
}
