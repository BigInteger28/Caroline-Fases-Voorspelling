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

// Function to read date input from user
func readDateInput(prompt string) (time.Time, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	date, err := time.Parse("02 01 2006", input)
	return date, err
}

// Function to find the phase for a specific date
func findPhaseForDate(date time.Time, cycles []Cycle) (string, time.Time, time.Time) {
	for _, cycle := range cycles {
		if (date.After(cycle.Start) && date.Before(cycle.End)) || date.Equal(cycle.Start) || date.Equal(cycle.End) {
			return cycle.Phase, cycle.Start, cycle.End
		}
	}
	return "Unknown", time.Time{}, time.Time{}
}

func main() {
	// Startdatum is vastgesteld op 22 maart 2024
	fixedStartDate := time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)

	// Lees de start- en eindmaand/jaar van de gebruiker
	for {
		fmt.Println("Kies een optie:")
		fmt.Println("1. Toon kalender")
		fmt.Println("2. Voer een datum in")
		fmt.Println("3. Afsluiten")
		fmt.Print("Maak een keuze: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			startMonthYear, err := readMonthYearInput("Enter start month and year (MM YYYY): ")
			if err != nil {
				fmt.Println("Invalid date format. Please use MM YYYY.")
				continue
			}
			endMonthYear, err := readMonthYearInput("Enter end month and year (MM YYYY): ")
			if err != nil {
				fmt.Println("Invalid date format. Please use MM YYYY.")
				continue
			}
			endMonthYear = endMonthYear.AddDate(0, 1, 0)

			// Bereken het aantal dagen tussen de vaste startdatum en de door de gebruiker opgegeven einddatum
			numDaysFromFixedStart := int(endMonthYear.Sub(fixedStartDate).Hours() / 24)

			// Genereer alle cycli vanaf de vaste startdatum
			allCycles := generateCycles(fixedStartDate, phases, numDaysFromFixedStart)

			// Filter cycli binnen het door de gebruiker opgegeven bereik
			fmt.Println("Calculated cycles:")
			for _, cycle := range allCycles {
				if cycle.Start.After(endMonthYear) {
					break
				}
				if cycle.End.After(startMonthYear) && cycle.Start.Before(endMonthYear) {
					fmt.Printf("Phase: %s, Start: %s, End: %s\n", cycle.Phase, cycle.Start.Format("02-01-2006"), cycle.End.Format("02-01-2006"))
				}
			}
		case 2:
			targetDate, err := readDateInput("Enter date (DD MM YYYY): ")
			if err != nil {
				fmt.Println("Invalid date format. Please use DD MM YYYY.")
				continue
			}

			// Bereken het aantal dagen tussen de vaste startdatum en de door de gebruiker opgegeven datum
			numDaysFromFixedStart := int(targetDate.Sub(fixedStartDate).Hours() / 24)

			// Genereer alle cycli vanaf de vaste startdatum
			allCycles := generateCycles(fixedStartDate, phases, numDaysFromFixedStart)

			// Zoek de fase voor de opgegeven datum
			phase, start, end := findPhaseForDate(targetDate, allCycles)
			if phase != "Unknown" {
				fmt.Printf("Date: %s is in phase %s which starts on %s and ends on %s\n", targetDate.Format("02-01-2006"), phase, start.Format("02-01-2006"), end.Format("02-01-2006"))
			} else {
				fmt.Println("Phase not found for the given date.")
			}
		case 3:
			fmt.Println("Afsluiten...")
			return
		default:
			fmt.Println("Ongeldige optie. Probeer het opnieuw.")
		}
	}
}
