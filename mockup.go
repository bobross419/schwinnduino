package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

type Console struct {
	button1          bool // toggle revolutions/rpms
	button2          bool // toggle time/distance
	revolutions      int
	revolution_timer time.Time // Timer used for spoofing revolutions
	display_timer    time.Time // Timer used to set display delay (500ms)
	elapsed_timer    time.Time // Timer used to track workout duration
	rpmdata          [5]rpmdatapoint
	//rpm              int // this is based on rpmdata[]
	//distance         int // this is calced so no need to hang onto it
}

type rpmdatapoint struct {
	revolutions int
	timestamp   time.Duration
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

	c := Console{false, false, 0, time.Now(), time.Now(), time.Now(), [5]rpmdatapoint{}}

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
				//c.distance = 0
				c.revolution_timer = time.Now()
				c.display_timer = time.Now()
				c.elapsed_timer = time.Now()
				c.rpmdata = [5]rpmdatapoint{}
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
				cycle_datapoints(&c)

				var (
					output1 string
					output2 string
				)

				// Button 1 Toggle
				if c.button1 {
					output1 = output_revolutions(&c)
				} else {
					output1 = output_rpm(&c)
				}

				// Button 2 Toggle
				if c.button2 {
					output2 = output_distance(&c)
				} else {
					output2 = output_timer(&c)
				}
				update_display(output1, output2)
				c.display_timer = time.Now()
			}
		}
	}
}

// Calc Functions
func calc_calories() {
	//still need to do some research on this, but should be a function of speed or rpms
}
func calc_distance(rev int) int {
	// brute force calc for now
	// this isn't very "real time" though
	// the results are stepped instead of linear
	// e.g. 16, 32, 48
	return rev * 16
}
func calc_rpm(a *[5]rpmdatapoint) int {
	first, last := a[0], a[len(a)-1]
	rev_delta := first.revolutions - last.revolutions
	time_delta := first.timestamp.Seconds() - last.timestamp.Seconds()
	rpm := float64(rev_delta) / time_delta
	//fmt.Printf("%d %f", rev_delta, time_delta)
	r := int(math.Ceil(rpm * 60))
	return r
}
func calc_speed() {
	//speed can be based on calculated RPM
}

// Output Functions
func output_distance(cPtr *Console) string {
	distance := calc_distance(cPtr.revolutions) // calculate distance
	return fmt.Sprintf("Distance: %dm", distance)
}
func output_revolutions(cPtr *Console) string {
	return fmt.Sprintf("Revolutions: %d", cPtr.revolutions)
}
func output_rpm(cPtr *Console) string {
	return fmt.Sprintf("RPM: %d", calc_rpm(&cPtr.rpmdata))
}
func output_timer(cPtr *Console) string {
	elapsed := time.Since(cPtr.elapsed_timer)
	return fmt.Sprintf("Time: %s", elapsed.String())
}

// Other Functions
func cycle_datapoints(cPtr *Console) {
	newdata := rpmdatapoint{cPtr.revolutions, time.Since(cPtr.elapsed_timer)}
	a := &cPtr.rpmdata
	for i := len(a) - 2; i >= 0; i-- {
		a[i+1] = a[i]
	}
	a[0] = newdata
}
func update_display(out1 string, out2 string) {
	fmt.Printf("%s %s\n", out1, out2)
}
