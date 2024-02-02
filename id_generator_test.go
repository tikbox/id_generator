package id_generator

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestUsageFlow(t *testing.T) {
	// Create an IdGenerator instance
	generator := NewIdGenerator(WithIdCount(10000))

	// Initialize the IdGenerator
	err := generator.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize IdGenerator: %v", err)
	}

	err = generator.LoadIds(generator.getCycleStartTime(0))
	if err != nil {
		t.Fatalf("Failed to LoadIds: %v", err)
	}

	// Note: Here is simulating the process of obtaining IDs from the generator.
	// It takes a timestamp as input and retrieves the corresponding ID. The ID is then marked as used.
	// The sleep duration between ID requests is included for demonstration purposes.
	for i := 0; i < 10; i++ {
		// Get an ID from the generator
		key := time.Now().UnixNano() / int64(generator.unitDuration)
		id := generator.GetId(key)
		t.Log("id=", id)

		// Mark the ID as used
		generator.MarkIdAsUsed(id)

		time.Sleep(generator.unitDuration * 2)
	}

	// Note: In a normal scenario, the synchronization of IDs to the file would typically be performed at the end of each cycle,
	//rather than immediately after marking an ID as used. The code here is for demonstration purposes only.
	err = generator.SyncIdsToFile()
	if err != nil {
		t.Fatalf("Failed to sync IDs to file: %v", err)
	}
	err = generator.LoadIds(generator.getCycleStartTime(1))
	if err != nil {
		t.Fatalf("Failed to LoadIds: %v", err)
	}
}

func TestInitialize(t *testing.T) {
	// Create a temporary file to simulate the existence of an ID file
	file, err := os.CreateTemp("", "id_list.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Create an IdGenerator instance
	generator := NewIdGenerator(WithFilename(file.Name()))

	// Initialize the IdGenerator
	err = generator.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize IdGenerator: %v", err)
	}

	// Ensure the correct number of generated IDs
	if len(generator.ids) != DefaultIdCount {
		t.Errorf("Expected %d IDs, but got %d", DefaultIdCount, len(generator.ids))
	}

	// Ensure the ID file is created and contains the correct number of IDs
	fileContent, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Failed to read ID file: %v", err)
	}

	lines := bytes.Count(fileContent, []byte("\n"))
	if lines != DefaultIdCount {
		t.Errorf("Expected %d lines in ID file, but got %d", DefaultIdCount, lines)
	}
}

func TestLoadIds(t *testing.T) {
	// Create a temporary file to simulate the existence of an ID file
	file, err := os.CreateTemp("", "id_list.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write simulated ID data to the file
	ids := []string{"1", "2", "3", "4", "5"}
	fileContent := []byte(strings.Join(ids, "\n"))
	err = os.WriteFile(file.Name(), fileContent, 0644)
	if err != nil {
		t.Fatalf("Failed to write ID data to file: %v", err)
	}

	// Create an IdGenerator instance
	generator := NewIdGenerator(WithFilename(file.Name()))

	// Load ID data into the idMap from the file
	startTime := time.Now()
	err = generator.LoadIds(startTime)
	if err != nil {
		t.Fatalf("Failed to load ID data to map: %v", err)
	}

	// Verify that idMap contains the correct keys and corresponding IDs
	for i, id := range ids {
		key := (startTime.UnixNano()/int64(time.Millisecond) + int64(i)) * int64(DefaultUnitDuration)
		expectedID, _ := strconv.Atoi(id)
		if generator.idMap[key] != expectedID {
			t.Errorf("Expected ID %d for key %d, but got %d", expectedID, key, generator.idMap[key])
		}
	}
}

func TestGenerateRandIds(t *testing.T) {
	// Create an IdGenerator instance
	generator := NewIdGenerator()

	// Generate random IDs
	generator.GenerateRandIds()

	// Ensure the correct number of generated IDs
	if len(generator.ids) != DefaultIdCount {
		t.Errorf("Expected %d generated IDs, but got %d", DefaultIdCount, len(generator.ids))
	}

	// Ensure that all generated IDs are within the specified range
	for _, id := range generator.ids {
		if id < generator.minId || id > generator.maxId {
			t.Errorf("Generated ID %d is out of range [%d, %d]", id, generator.minId, generator.maxId)
		}
	}
}

func TestSaveIdsToFile(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "id_list.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Create an IdGenerator instance
	generator := NewIdGenerator(WithFilename(file.Name()))

	// Generate random IDs
	generator.GenerateRandIds()

	// Save the generated IDs to a file
	err = generator.SaveIdsToFile()
	if err != nil {
		t.Fatalf("Failed to save IDs to file: %v", err)
	}

	// Read the file content
	fileContent, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Failed to read ID file: %v", err)
	}

	// Verify that the file content matches the generated IDs
	idStrings := make([]string, len(generator.ids))
	for i, id := range generator.ids {
		idStrings[i] = strconv.Itoa(id)
	}

	expectedContent := strings.Join(idStrings, "\n")
	actualContent := string(fileContent)

	if actualContent != expectedContent {
		t.Errorf("Expected file content:\n%s\n\nBut got:\n%s", expectedContent, actualContent)
	}
}

func TestGetId(t *testing.T) {
	// Create an IdGenerator instance
	generator := NewIdGenerator()

	// Generate random IDs
	generator.GenerateRandIds()

	// Get an ID from the generator
	id := generator.GetId(time.Now().UnixMicro() / int64(time.Millisecond))

	// Ensure that the ID is within the specified range
	if id < generator.minId || id > generator.maxId {
		t.Errorf("Got an ID %d that is out of range [%d, %d]", id, generator.minId, generator.maxId)
	}

	// Ensure that the ID is marked as used
	if !generator.usedIdMap[id] {
		t.Errorf("ID %d is not marked as used", id)
	}
}

func TestMarkIdAsUsed(t *testing.T) {
	// Create an IdGenerator instance
	generator := NewIdGenerator()

	// Generate random IDs
	generator.GenerateRandIds()

	// Get an ID from the generator
	id := generator.GetId(time.Now().UnixMicro() / int64(time.Millisecond))

	// Mark the ID as used
	generator.MarkIdAsUsed(id)

	// Ensure that the ID is marked as used
	if !generator.usedIdMap[id] {
		t.Errorf("ID %d is not marked as used", id)
	}
}

func TestSyncIdsToFile(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "id_list.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Create an IdGenerator instance
	generator := NewIdGenerator(WithFilename(file.Name()))

	// Generate random IDs
	generator.GenerateRandIds()

	// Sync the IDs to the file
	err = generator.SyncIdsToFile()
	if err != nil {
		t.Fatalf("Failed to sync IDs to file: %v", err)
	}

	// Read the file content
	fileContent, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Failed to read ID file: %v", err)
	}

	// Verify that the file content matches the generated IDs
	idStrings := make([]string, len(generator.ids))
	for i, id := range generator.ids {
		idStrings[i] = strconv.Itoa(id)
	}

	expectedContent := strings.Join(idStrings, "\n")
	actualContent := string(fileContent)

	if actualContent != expectedContent {
		t.Errorf("Expected file content:\n%s\n\nBut got:\n%s", expectedContent, actualContent)
	}
}

func TestGetCycleStartTime(t *testing.T) {
	// Create an IdGenerator instance
	generator := NewIdGenerator()

	// Get the cycle start time
	startTime := generator.getCycleStartTime(0)

	// Ensure that the start time is not zero
	if startTime.IsZero() {
		t.Errorf("Cycle start time is zero")
	}

	// Ensure that the start time is within a reasonable range
	currentTime := time.Now()
	minStartTime := currentTime.Add(-time.Hour)
	maxStartTime := currentTime.Add(time.Hour)

	if startTime.Before(minStartTime) || startTime.After(maxStartTime) {
		t.Errorf("Cycle start time %v is not within a reasonable range", startTime)
	}
}
