package main

import (
	"fmt"
	"log"
	"progetto/restaurant/server/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.InitDatabase("../server/restaurant.db")
	defer database.CloseDatabase()

	fmt.Println("=== Populating Tables ===")
	fmt.Println()

	tables := []struct {
		seats int
		count int
	}{
		{2, 4}, // 4 tables of 2 seats
		{4, 4}, // 4 tables of 4 seats
		{6, 1}, // 1 table of 6 seats
	}

	totalTables := 0
	for _, tableType := range tables {
		totalTables += tableType.count
	}

	fmt.Printf("Inserting %d physical tables...\n\n", totalTables)

	tableNumber := 1
	for _, tableType := range tables {
		for i := 0; i < tableType.count; i++ {
			err := database.InsertTable(tableType.seats)
			if err != nil {
				log.Printf("Error inserting table %d: %v", tableNumber, err)
			} else {
				fmt.Printf("Table #%d: %d seats\n", tableNumber, tableType.seats)
			}
			tableNumber++
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("          SUMMARY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  Table composition:")
	for _, tableType := range tables {
		fmt.Printf("     • %d tables × %d seats\n", tableType.count, tableType.seats)
	}

	// Calcola capacità totale
	totalSeats := 0
	for _, tableType := range tables {
		totalSeats += tableType.count * tableType.seats
	}
	fmt.Printf("\nTotal capacity: %d seats\n", totalSeats)
	fmt.Printf("Total tables: %d\n", totalTables)
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("	Restaurant ready for reservations!")
	fmt.Println("   Each table is available for all time slots (12:00-22:00)")
	fmt.Println("   Reservations last 2 hours per booking")
}
