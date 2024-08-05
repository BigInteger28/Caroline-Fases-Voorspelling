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
	Name      string
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
	phaseIndex := 0
	durationIndex := 0

	for dayCount := 0; dayCount < numDays; {
		phase := phases[phaseIndex%len(phases)]
		duration := phase.Durations[durationIndex%len(phase.Durations)]
		endDate := currentDate.AddDate(0, 0, duration-1)
		cycles = append(cycles, Cycle{Phase: phase.Name, Start: currentDate, End: endDate})
		currentDate = endDate.AddDate(0, 0, 1)
		dayCount += duration
		phaseIndex++
		if phaseIndex%len(phases) == 0 {
			durationIndex++
		}
	}
	return cycles
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
	for {
		// Startdatum is vastgesteld op 22 maart 2024
		fixedStartDate := time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)
	
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
		endMonthYear = endMonthYear.AddDate(0, 1, 0)
	
		// Calculate the number of days between the fixed start date and the user-specified end date
		numDaysFromFixedStart := int(endMonthYear.Sub(fixedStartDate).Hours() / 24)
	
		// Generate all cycles from the fixed start date
		allCycles := generateCycles(fixedStartDate, phases, numDaysFromFixedStart)
	
		// Filter cycles within the user-specified range
		fmt.Println("Calculated cycles:")
		for _, cycle := range allCycles {
			if cycle.Start.After(endMonthYear) {
				break
			}
			if cycle.End.After(startMonthYear) && cycle.Start.Before(endMonthYear) {
				fmt.Printf("Phase: %s, Start: %s, End: %s\n", cycle.Phase, cycle.Start.Format("02-01-2006"), cycle.End.Format("02-01-2006"))
			}
		}
	}
}
