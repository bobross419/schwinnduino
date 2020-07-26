package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Console struct {
	button1          bool
	button2          bool // toggle time/distance
	revolutions      int
	distance         int
	revolution_timer time.Time
	display_timer    time.Time
	elapsed_timer    time.Time
}

func main() {
	fmt.Println("vim-go")
	ch := make(chan string)
	// Read keyboard input
	go func(ch chan string) {
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	c := Console{false, false, 0, 0, time.Now(), time.Now(), time.Now()}

	for {
		select {
		// Spoof inputs
		case stdin, _ := <-ch:
			if stdin == "q" { // Power
				os.Exit(0)
			}
			if stdin == "a" { // Output1 Toggle
				c.button1 = !c.button1
			}
			if stdin == "b" { // Output2 Toggle
				c.button2 = !c.button2
			}
			if stdin == "r" { // Reset Button
				c.revolutions = 0
				c.distance = 0
				c.revolution_timer = time.Now()
				c.display_timer = time.Now()
				c.elapsed_timer = time.Now()
			}

		default:
			// spoof the wheel spinning
			revolution_duration := time.Since(c.revolution_timer)
			if revolution_duration.Seconds() > 1.0 {
				c.revolutions++
				c.revolution_timer = time.Now()
			}

			// Only update the display every .5 seconds
			time_duration := time.Since(c.display_timer)
			if time_duration.Milliseconds() > 500 {
				c.distance = calculate_distance(c.revolutions) // don't need to calculate this unless we're updating the display
				output1 := strconv.Itoa(c.revolutions)
				var output2 string
				if c.button2 {
					output2 = strconv.Itoa(c.distance)
				} else {
					elapsed := time.Since(c.elapsed_timer)
					output2 = elapsed.String()
				}
				update_display(output1, output2)
				c.display_timer = time.Now()
			}
		}
	}
}

func calculate_distance(rev int) int {
	return rev * 16
}

func update_display(out1 string, out2 string) {
	fmt.Printf("%s %s\n", out1, out2)
}
