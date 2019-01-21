package subprocessor

import (
	"fmt"
	"math"
	"testing"
)

func TestCoordinatorStart(t *testing.T) {

	c := NewCoordinator()

	go c.Run()

	c.Destroy() //Here under normal conditions the async destroy does not get called as the thread exits on main

}

func TestGetCoordinatorInputSimple(t *testing.T) {

	SieveOfEratosthenes := func(c chan bool, value int) {
		fmt.Println("Running SieveOfEratosthenes")
		f := make([]bool, value)
		for i := 2; i <= int(math.Sqrt(float64(value))); i++ {
			if f[i] == false {
				for j := i * i; j < value; j += i {
					f[j] = true
				}
			}
		}
		for i := 2; i < value; i++ {
			if f[i] == false {
				fmt.Printf("%v", i)
			}
		}
		c <- true
	}
	c := NewCoordinator()

	go c.Run()

	cb := make(chan bool)

	c.ResourceChannel <- Resource{
		Process: func() {
			SieveOfEratosthenes(cb, 42)
		},
	}

	<-cb
}

//
//func TestGetCoordinatorInputExtended(t *testing.T) {
//
//	c := GetCoordinator()
//
//	urls := []string{
//		"https://www.google.com",
//		"https://www.yahoo.com",
//		"https://www.msn.com",
//		"https://www.yandex.com",
//	}
//	var wg sync.WaitGroup
//	for _, url := range urls {
//		c.ResourceChannel <- Resource{
//			Process: func() {
//				wg.Add(1)
//				resp, err := http.Get(url)
//				if err != nil {
//					t.Fail()
//				}
//				fmt.Printf("Got data from %s\n", url)
//				color.Yellow(resp.Status)
//				defer resp.Body.Close()
//				wg.Done()
//			},
//		}
//	}
//	wg.Wait()
//	c.Stop()
//}
