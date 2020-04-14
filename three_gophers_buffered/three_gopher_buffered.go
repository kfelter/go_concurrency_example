package main

import (
	"fmt"
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
	start := time.Now()

	pileOfBooks := MakePileOfBooks(500)
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

	pileLoaderChan := make(chan *Cart, 2)
	moverChan := make(chan *Cart, 2)
	incineratorLoaderChan := make(chan *Cart, 2)

	go func() {
		for c := range pileLoaderChan {
			fmt.Printf("%5d Books remain\n", len(pileOfBooks))
			pileOfBooks = LoadCart(pileOfBooks, c)
			moverChan <- c
		}
	}()

	go func() {
		for c := range moverChan {
			// fmt.Println("moving cart with", len(c.Books), "books")
			MoveCart(c)
			if c.ReadyToLoad {
				pileLoaderChan <- c
			} else {
				incineratorLoaderChan <- c
			}
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
			moverChan <- c
		}
	}()

	pileLoaderChan <- &cart
	pileLoaderChan <- &Cart{
		Books:       []Book{},
		Cap:         15,
		ReadyToLoad: true,
	}

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
