package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	phaseDurations = map[string][]int{
		"Menstruatie": {4, 5, 6},
		"Piek":        {4, 5, 6},
		"Ovulatie":    {7},
		"Luteaal":     {11},
	}
	startDate  = time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)
	maxYears   = 1
	threadPool = 10
)

// Cycle struct to hold the start and end dates of each phase
type Cycle struct {
	Phase string
	Start time.Time
	End   time.Time
}

// Function to calculate all possible cycles
func calculateAllCycles(startDate time.Time, durations map[string][]int, maxYears int) [][]Cycle {
	var allCycles [][]Cycle
	combinations := generateCombinations(durations)
	for _, combination := range combinations {
		cycles := generateCyclesForCombination(startDate, combination, maxYears)
		allCycles = append(allCycles, cycles...)
	}
	return allCycles
}

// Function to generate all possible combinations of phase durations
func generateCombinations(durations map[string][]int) [][]int {
	var result [][]int
	var keys []string
	for k := range durations {
		keys = append(keys, k)
	}
	var recurse func([]int, int)
	recurse = func(current []int, depth int) {
		if depth == len(keys) {
			temp := make([]int, len(current))
			copy(temp, current)
			result = append(result, temp)
			return
		}
		for _, v := range durations[keys[depth]] {
			recurse(append(current, v), depth+1)
		}
	}
	recurse([]int{}, 0)
	return result
}

// Function to generate cycles for a given combination of phase durations
func generateCyclesForCombination(startDate time.Time, combination []int, maxYears int) [][]Cycle {
	var cycles [][]Cycle
	currentDate := startDate
	endDate := startDate.AddDate(maxYears, 0, 0)

	for currentDate.Before(endDate) {
		var cycle []Cycle
		phaseNames := []string{"Menstruatie", "Piek", "Ovulatie", "Luteaal"}
		for i, phase := range phaseNames {
			duration := combination[i]
			cycle = append(cycle, Cycle{
				Phase: phase,
				Start: currentDate,
				End:   currentDate.AddDate(0, 0, duration-1),
			})
			currentDate = currentDate.AddDate(0, 0, duration)
		}
		cycles = append(cycles, cycle)
	}
	return cycles
}

// Function to calculate phase probabilities for a given date
func calculateProbabilities(cycles [][]Cycle, targetDate time.Time, result chan map[string]float64, wg *sync.WaitGroup) {
	defer wg.Done()
	probabilities := make(map[string]float64)
	for _, cycle := range cycles {
		for _, phase := range cycle {
			if (targetDate.After(phase.Start) && targetDate.Before(phase.End)) || targetDate.Equal(phase.Start) || targetDate.Equal(phase.End) {
				probabilities[phase.Phase]++
			}
		}
	}

	total := 0.0
	for _, count := range probabilities {
		total += count
	}
	for phase := range probabilities {
		probabilities[phase] = (probabilities[phase] / total) * 100
	}
	result <- probabilities
}

// Function to find the days in a given month where a specific phase can occur
func findPhaseDays(cycles [][]Cycle, month time.Month, year int, targetPhase string) map[time.Time]float64 {
	daysInPhase := make(map[time.Time]float64)
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	for _, cycle := range cycles {
		for _, phase := range cycle {
			if phase.Phase == targetPhase {
				for d := phase.Start; !d.After(phase.End); d = d.AddDate(0, 0, 1) {
					if !d.Before(startOfMonth) && !d.After(endOfMonth) {
						daysInPhase[d]++
					}
				}
			}
		}
	}

	totalDays := float64(len(cycles))
	for day := range daysInPhase {
		daysInPhase[day] = (daysInPhase[day] / totalDays) * 100
	}
	return daysInPhase
}

// Function to find the best start day for a phase in each month of a given year
func findBestStartDays(cycles [][]Cycle, year int, targetPhase string) map[time.Month]time.Time {
	bestDays := make(map[time.Month]time.Time)
	for month := 1; month <= 12; month++ {
		daysInPhase := findPhaseDays(cycles, time.Month(month), year, targetPhase)
		var bestDay time.Time
		highestProbability := 0.0
		for day, probability := range daysInPhase {
			if probability > highestProbability {
				bestDay = day
				highestProbability = probability
			}
		}
		if !bestDay.IsZero() {
			bestDays[time.Month(month)] = bestDay
		}
	}
	return bestDays
}

// Function to read date input from user in the format dd mm yyyy
func readDateInput() (time.Time, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter date (dd mm yyyy): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	date, err := time.Parse("02 01 2006", input)
	return date, err
}

// Function to read month and year input from user in the format mm yyyy
func readMonthYearInput() (time.Month, int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter month and year (mm yyyy): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	date, err := time.Parse("01 2006", input)
	return date.Month(), date.Year(), err
}

// Function to read year input from user in the format yyyy
func readYearInput() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter year (yyyy): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	date, err := time.Parse("2006", input)
	return date.Year(), err
}

// Function to read phase input from user
func readPhaseInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter phase (Menstruatie, Piek, Ovulatie, Luteaal): ")
	phase, _ := reader.ReadString('\n')
	return strings.TrimSpace(phase)
}

func main() {
	// Ask for number of years to calculate
	fmt.Print("Enter the number of years to calculate: ")
	fmt.Scanln(&maxYears)

	// Calculate all possible cycles
	allCycles := calculateAllCycles(startDate, phaseDurations, maxYears)

	// Menu options
	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Phase probability for a specific date")
		fmt.Println("2. Phase days for a specific month")
		fmt.Println("3. Best start days for a phase in each month of a given year")
		fmt.Println("4. Show calculated cycles for a specific number of years")
		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			// Read target date from user
			targetDate, err := readDateInput()
			if err != nil {
				fmt.Println("Invalid date format. Please use dd mm yyyy.")
				continue
			}

			// Calculate probabilities using goroutines
			results := make(chan map[string]float64, threadPool)
			var wg sync.WaitGroup

			cycleChunks := len(allCycles) / threadPool
			for i := 0; i < threadPool; i++ {
				start := i * cycleChunks
				end := start + cycleChunks
				if i == threadPool-1 {
					end = len(allCycles)
				}

				wg.Add(1)
				go calculateProbabilities(allCycles[start:end], targetDate, results, &wg)
			}

			go func() {
				wg.Wait()
				close(results)
			}()

			// Aggregate results
			finalProbabilities := make(map[string]float64)
			for res := range results {
				for phase, prob := range res {
					finalProbabilities[phase] += prob
				}
			}

			// Normalize results
			total := 0.0
			for _, count := range finalProbabilities {
				total += count
			}
			for phase := range finalProbabilities {
				finalProbabilities[phase] = (finalProbabilities[phase] / total) * 100
			}

			// Print probabilities
			fmt.Printf("Probabilities for %s:\n", targetDate.Format("02 01 2006"))
			for phase, probability := range finalProbabilities {
				fmt.Printf("%s: %.2f%%\n", phase, probability)
			}

		case 2:
			// Read month and year from user
			month, year, err := readMonthYearInput()
			if err != nil {
				fmt.Println("Invalid date format. Please use mm yyyy.")
				continue
			}

			// Read phase from user
			targetPhase := readPhaseInput()

			// Find days in the given month where the target phase can occur
			daysInPhase := findPhaseDays(allCycles, month, year, targetPhase)

			// Find the day with the highest probability
			var bestDay time.Time
			highestProbability := 0.0
			for day, probability := range daysInPhase {
				if probability > highestProbability {
					bestDay = day
					highestProbability = probability
				}
			}

			// Print the best day and its probability
			if !bestDay.IsZero() {
				fmt.Printf("The best day for phase %s in %s %d is %s with a probability of %.2f%%\n", targetPhase, month, year, bestDay.Format("02 01 2006"), highestProbability)
			} else {
				fmt.Printf("No days found for phase %s in %s %d\n", targetPhase, month, year)
			}

		case 3:
			// Read year from user
			year, err := readYearInput()
			if err != nil {
				fmt.Println("Invalid year format. Please use yyyy.")
				continue
			}

			// Read phase from user
			targetPhase := readPhaseInput()

			// Find the best start days for the phase in each month of the given year
			bestDays := findBestStartDays(allCycles, year, targetPhase)

			// Sort the months for display
			months := make([]time.Month, 0, len(bestDays))
			for month := range bestDays {
				months = append(months, month)
			}
			sort.Slice(months, func(i, j int) bool {
				return months[i] < months[j]
			})

			// Print the best days
			fmt.Printf("Best start days for phase %s in %d:\n", targetPhase, year)
			for _, month := range months {
				day := bestDays[month]
				fmt.Printf("%s: %s\n", month.String(), day.Format("02 01 2006"))
			}

		case 4:
			// Read number of years from user
			fmt.Print("Enter the number of years to show: ")
			var numYears int
			fmt.Scanln(&numYears)

			// Calculate the end date based on the number of years
			endDate := startDate.AddDate(numYears, 0, 0)

			// Print calculated cycles for the specified number of years
			fmt.Println("Calculated cycles:")
			for _, cycle := range allCycles {
				for _, phase := range cycle {
					if phase.Start.Before(endDate) {
						fmt.Printf("Phase: %s, Start: %s, End: %s\n", phase.Phase, phase.Start.Format("02 01 2006"), phase.End.Format("02 01 2006"))
					}
				}
			}

		default:
			fmt.Println("Invalid option.")
		}
	}
}
