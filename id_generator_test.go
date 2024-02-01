package id_generator

import (
	"bytes"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	IdFile      = "id_list.txt"
	idMapLength = 3600
)

// IdGenerator class
type IdGenerator struct {
	ids       []int
	idMap     map[int64]int
	usedIdMap map[int]bool
	filename  string
	mu        sync.Mutex
}

// NewIdGenerator creates a new instance of IdGenerator
func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		usedIdMap: make(map[int]bool),
		filename:  IdFile,
	}
}

// GenerateRandIds generates randomly shuffled 6-digit numeric Ids
func (g *IdGenerator) GenerateRandIds(count int) {
	g.ids = make([]int, count)

	for i := 0; i < count; i++ {
		g.ids[i] = i + 100000 // Generate 6-digit numeric Id
	}

	rand.Seed(time.Now().UnixNano())

	// Fisher-Yates shuffle algorithm
	for i := count - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		g.ids[i], g.ids[j] = g.ids[j], g.ids[i]
	}
}

// SaveIdsToFile saves Ids to a file
func (g *IdGenerator) SaveIdsToFile() error {
	file, err := os.Create(g.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	idsStr := make([]string, len(g.ids))
	for i, id := range g.ids {
		idsStr[i] = strconv.Itoa(id)
	}

	_, err = file.WriteString(strings.Join(idsStr, "\n"))
	if err != nil {
		return err
	}

	return nil
}

// LoadIdsToMap loads corresponding Ids into the in-memory map
func (g *IdGenerator) LoadIdsToMap(startSeconds int64) error {
	g.idMap = make(map[int64]int)

	content, err := os.ReadFile(g.filename)
	if err != nil {
		return err
	}

	idsStr := bytes.Split(content, []byte("\n"))

	for i := 0; i < idMapLength; i++ {
		key := startSeconds + int64(i)
		g.idMap[key], _ = strconv.Atoi(string(idsStr[i]))
	}

	return nil
}

// GetId retrieves the Id corresponding to the specified seconds
func (g *IdGenerator) GetId(seconds int64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	if id, ok := g.idMap[seconds]; ok {
		if !g.usedIdMap[id] {
			return id
		}
	}

	return 0
}

// MarkIdAsUsed marks the Id as used
func (g *IdGenerator) MarkIdAsUsed(id int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if id == 0 {
		return
	}

	g.usedIdMap[id] = true
}

// SyncIdsToFile synchronizes Id data to the file
func (g *IdGenerator) SyncIdsToFile() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	file, err := os.Create(g.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	idsStr := make([]string, 0)
	for _, id := range g.ids {
		if !g.usedIdMap[id] {
			idsStr = append(idsStr, strconv.Itoa(id))
		}
	}

	_, err = file.WriteString(strings.Join(idsStr, "\n"))
	if err != nil {
		return err
	}

	g.usedIdMap = make(map[int]bool)

	return nil
}

func getStartOfHourSeconds() int64 {
	now := time.Now()
	startOfHour := now.Truncate(time.Hour).Unix()
	return startOfHour
}

func getNextHourSeconds() int64 {
	nextHour := time.Now().Add(time.Hour)
	nextHour = time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), nextHour.Hour(), 0, 0, 0, nextHour.Location())
	return nextHour.Unix()
}
