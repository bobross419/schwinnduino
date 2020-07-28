package main

/*
	This was meant to exercise different types of datasets: arrays, rings, lists
	Criteria is size, speed, and complexity of code
	Tests were attempted with container/ring as well. However, due to the limited
	functions it would add needless complexity to compare any but two elements
   	that are right next to each other (there is a word for that but it eludes me)
*/

import (
	"container/list" // slightly larger and magnitudes slower on cycling than the array
	//"container/ring"*
	"fmt"
	"time"
	"unsafe"
)

/*
 */

func main() {
	a := [5]int{}
	l := list.New()

	fmt.Printf("a: %T, %d\n", a, unsafe.Sizeof(a))
	fmt.Printf("l: %T, %d\n", l, unsafe.Sizeof(*l))

	for i := 0; i < len(a); i++ {
		a[i] = i
		l.PushBack(i)
	}

	fmt.Println("===========FILLED===========")
	fmt.Printf("a: %T, %d\n", a, unsafe.Sizeof(a))
	fmt.Printf("l: %T, %d\n", l, unsafe.Sizeof(*l))

	fmt.Println("===========TIMERS===========")

	cycles := 10000

	// TIME ARRAY
	at := time.Now()
	for i := 0; i < cycles; i++ {
		for j := len(a) - 2; j >= 0; j-- {
			a[j+1] = a[j]
		}
		a[0] = i
	}
	fmt.Printf("a cycle: %v\n", time.Since(at))

	// TIME ARRAY FUNCTION
	// Just invoking the function makes this noticeably slower
	at = time.Now()
	for i := 0; i < cycles; i++ {
		cycle_array(&a, i)
	}
	fmt.Printf("a cycle(): %v\n", time.Since(at))

	// TIME LIST
	lt := time.Now()
	for i := 0; i < cycles; i++ {
		l.Remove(l.Back())
		l.PushFront(i)
	}
	fmt.Printf("l cycle: %v\n", time.Since(lt))

	fmt.Printf("====%d CYCLES====\n", cycles)
	fmt.Printf("a: %T, %d\n", a, unsafe.Sizeof(a))
	fmt.Printf("l: %T, %d\n", l, unsafe.Sizeof(*l))
}

func cycle_array(a *[5]int, newdata int) {
	for j := len(a) - 2; j >= 0; j-- {
		a[j+1] = a[j]
	}
	a[0] = newdata
}
