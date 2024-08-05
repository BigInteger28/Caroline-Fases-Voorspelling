package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Phase struct to hold the phase name and duration
type Phase struct {
	Name     string
	Durations []int
}

// Cycle struct to hold the start and end dates of each phase
type Cycle struct {
	Phase string
	Start time.Time
	End   time.Time
}

var (
	phases = []Phase{
		{"Menstruatie", []int{5, 6, 6, 4, 5, 6}},
		{"Piek", []int{6, 4, 5, 5, 6, 5}},
		{"Ovulatie", []int{7, 7, 7, 7, 7, 7}},
		{"Luteaal", []int{11, 11, 11, 11, 11, 11}},
	}
	cycles []Cycle
)

// Function to generate cycles
func generateCycles(startDate time.Time, phases []Phase, numDays int) []Cycle {
	var cycles []Cycle
	currentDate := startDate
	for currentDate.Before(startDate.AddDate(0, 0, numDays)) {
		for _, phase := range phases {
			duration := phase.Durations[len(cycles)%len(phase.Durations)]
			endDate := currentDate.AddDate(0, 0, duration-1)
			cycles = append(cycles, Cycle{Phase: phase.Name, Start: currentDate, End: endDate})
			currentDate = endDate.AddDate(0, 0, 1)
		}
	}
	return cycles
}

// Function to find the phase for a specific date
func findPhaseForDate(date time.Time, cycles []Cycle) string {
	for _, cycle := range cycles {
		if (date.After(cycle.Start) && date.Before(cycle.End)) || date.Equal(cycle.Start) || date.Equal(cycle.End) {
			return cycle.Phase
		}
	}
	return "Unknown"
}

// Function to read month and year input from user
func readMonthYearInput(prompt string) (time.Time, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	date, err := time.Parse("01 2006", input)
	return date, err
}

func main() {
	// Startdatum is vastgesteld op 22 maart 2024
	startDate := time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)

	// Read start and end month/year from user
	startMonthYear, err := readMonthYearInput("Enter start month and year (MM YYYY): ")
	if err != nil {
		fmt.Println("Invalid date format. Please use MM YYYY.")
		return
	}
	endMonthYear, err := readMonthYearInput("Enter end month and year (MM YYYY): ")
	if err != nil {
		fmt.Println("Invalid date format. Please use MM YYYY.")
		return
	}

	// Ensure that the start date is not before the fixed start date
	if startMonthYear.Before(startDate) {
		startMonthYear = startDate
	}

	// Calculate the number of days between start and end date
	numDays := int(endMonthYear.Sub(startMonthYear).Hours() / 24)

	// Generate cycles
	cycles = generateCycles(startDate, phases, numDays)

	// Print all cycles for the specified period
	fmt.Println("Calculated cycles:")
	for _, cycle := range cycles {
		if cycle.Start.After(endMonthYear) {
			break
		}
		if cycle.End.After(startMonthYear) && cycle.Start.Before(endMonthYear) {
			fmt.Printf("Phase: %s, Start: %s, End: %s\n", cycle.Phase, cycle.Start.Format("02-01-2006"), cycle.End.Format("02-01-2006"))
		}
	}
}
