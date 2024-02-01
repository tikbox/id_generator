package id_generator

import (
	"fmt"
	"testing"
	"time"
)

func TestIdGenerator(t *testing.T) {
	idGenerator := NewIdGenerator()

	// Generate unique 6-digit Ids
	//idGenerator.GenerateRandIds(899999)
	idGenerator.GenerateRandIds(10000)

	// Save Ids to a file
	err := idGenerator.SaveIdsToFile()
	if err != nil {
		fmt.Println("Error saving Ids to file:", err)
		return
	}

	// Load corresponding Ids into the in-memory map
	err = idGenerator.LoadIdsToMap(getStartOfHourSeconds())
	if err != nil {
		fmt.Println("Error loading Ids to map:", err)
		return
	}

	for i := 0; i < 100; i++ {
		// Simulate user registration and get an Id
		seconds := time.Now().Unix()
		id := idGenerator.GetId(seconds)
		fmt.Println("Registered Id:", id)

		// Mark the Id as used
		idGenerator.MarkIdAsUsed(id)

		time.Sleep(time.Second)
		//time.Sleep(time.Millisecond * 200)
	}

	// Sync data to the file
	err = idGenerator.SyncIdsToFile()
	if err != nil {
		fmt.Println("Error syncing Ids to file:", err)
		return
	}

	// Load corresponding Ids into the in-memory map
	err = idGenerator.LoadIdsToMap(getNextHourSeconds())
	if err != nil {
		fmt.Println("Error loading Ids to map:", err)
		return
	}
}
