package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
)

type Console struct {
	button1          bool // toggle revolutions/rpms
	button2          bool // toggle time/distance
	revolutions      uint
	revolution_timer time.Time // Timer used for spoofing revolutions
	display_timer    time.Time // Timer used to set display delay (500ms)
	elapsed_timer    time.Time // Timer used to track workout duration
	rpmdata          [5]rpmdatapoint
	//rpm              int // this is based on rpmdata[]
	//distance         int // this is calced so no need to hang onto it
}

type rpmdatapoint struct {
	revolutions uint
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
				reset_console(&c)
			}
			if stdin == "m" { // Debugger, show memory
				display_memory()
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
func calc_distance(rev uint) uint {
	// brute force calc for now
	// this isn't very "real time" though
	// the results are stepped instead of linear
	// e.g. 16, 32, 48
	// being a function of rpmdata is probably better
	return rev * 16
}
func calc_rpm(a *[5]rpmdatapoint) int {
	// Some of the issues I currently have with calculating RPM is the way I'm polling the data
	// Because I'm fixed at 60 seconds and the screen refresh and stuff are all .5 or 1s inteverals
	// Things are getting a bit munged
	// I should be able to use the time delta to calculate the RPM in real time
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
func calc_watts() {
	// Watts = .0011RPM3 + .0026RPM2 + 0.5642*RPM  (this might not be accurate but probs a good place to start)
	// This formula doesnt' line up with the AD3 Load Calibration.  Load 3 = 50RPM = 147.1 Watts, Formula with 50RPM = 172.2W
	// Ref: https://www.reddit.com/r/crossfit/comments/7bj678/assault_bike_rpm_watts_to_calories_graph/

	// This has some data I might be able to use as well
	// Ref: https://docs.google.com/spreadsheets/d/1Leq-FH6-2smtYyYNMdphR8L_6yC48N8vYlibjd6FxAQ/edit#gid=1018308001

	// According to our research we found that it takes about 1200 watts to generate 1 calorie on the Echo Bike. 50 rpm = 149 watts x 8seconds = 1192 = 1 calorie.
	// Also an rpm -> watts -> time/cal chart that could be a handy reference
	// Ref: https://heatonminded.com/rogue-echo-bike-calorie-conversion-chart/

	// Energy (kCal or "calories") = Av Power (Watts) x Duration (hours) x 3.6.
	// Not sure how accurate this is, buried in one of the answers
	// Ref: https://www.fixya.com/support/t25636091-user_guide_manual_schwinn_airdyne_ad3

	// Chart from the AD3 manual that I can use to compare results
	// https://www.facebook.com/229236226353/photos/a.10153499277066354/10153499324746354/?type=3&theater

	// Elevation calibration from AD3 manual
	// https://www.facebook.com/229236226353/photos/a.10153499277066354/10153499325521354/?type=3&theater

	// Checking calibration from AD3 manaul
	// Load 3 = 50 RPM
	// https://www.facebook.com/229236226353/photos/a.10153499277066354/10153499325536354/?type=3&theater
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
func cycle_datapoints(cPtr *Console) { // TODO Move this out of a function
	/*
		I experimented with cycling rings and linked lists,
		but I was unsatisfied with the results.
		rings: Code is needlessly complex for no noticeable speed
				however, the size was slightly smaller
		lists: Easy to implement, but magnitudes slower than arrays
				and slightly larger
		arrays: For our limited purposes this is small, fast, and
				fairly easy to follow.  Wrapping this in a function
				reduces speed.
	*/
	newdata := rpmdatapoint{cPtr.revolutions, time.Since(cPtr.elapsed_timer)}
	a := &cPtr.rpmdata
	for i := len(a) - 2; i >= 0; i-- {
		a[i+1] = a[i]
	}
	a[0] = newdata
}

func reset_console(cPtr *Console) {
	cPtr.revolutions = 0
	//c.distance = 0
	cPtr.revolution_timer = time.Now()
	cPtr.display_timer = time.Now()
	cPtr.elapsed_timer = time.Now()
	cPtr.rpmdata = [5]rpmdatapoint{}
}

func update_display(out1 string, out2 string) {
	fmt.Printf("%s %s\n", out1, out2)
}

// Debugging functions
func display_memory() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
