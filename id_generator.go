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
	DefaultCycleDuration = time.Hour
	DefaultUnitDuration  = time.Second
	DefaultFilename      = "id_list.txt"
	DefaultIdCount       = 1_000_000
	DefaultMinId         = 100_000
	DefaultMaxId         = 1_000_000
)

// IdGenerator class represents a generator for generating unique IDs.
type IdGenerator struct {
	ids           []int         // Slice of generated IDs
	idMap         map[int64]int // Map storing corresponding IDs
	usedIdMap     map[int]bool  // Map tracking used ID
	filename      string        // Filename for ID storage
	cycleDuration time.Duration // Duration of each cycle
	unitDuration  time.Duration // Duration of each unit
	idMapLength   int           // Length of idMap
	idCount       int           // Number of IDs to generate
	minId         int           // Minimum ID value
	maxId         int           // Maximum ID value
	mu            sync.Mutex    // Mutex for synchronization
}

// NewIdGenerator creates a new instance of IdGenerator with optional cycle duration, unit duration, and filename
func NewIdGenerator(options ...func(*IdGenerator)) *IdGenerator {
	g := &IdGenerator{
		usedIdMap:     make(map[int]bool),
		cycleDuration: DefaultCycleDuration,
		unitDuration:  DefaultUnitDuration,
		filename:      DefaultFilename,
		idCount:       DefaultIdCount,
		minId:         DefaultMinId,
		maxId:         DefaultMaxId,
	}

	for _, option := range options {
		option(g)
	}

	g.idMapLength = int(g.cycleDuration / g.unitDuration)

	return g
}

// WithCycleDuration sets the cycle duration for the IdGenerator
func WithCycleDuration(duration time.Duration) func(*IdGenerator) {
	return func(g *IdGenerator) {
		g.cycleDuration = duration
	}
}

// WithUnitDuration sets the unit duration for the IdGenerator
func WithUnitDuration(duration time.Duration) func(*IdGenerator) {
	return func(g *IdGenerator) {
		g.unitDuration = duration
	}
}

// WithFilename sets the filename for the IdGenerator
func WithFilename(filename string) func(*IdGenerator) {
	return func(g *IdGenerator) {
		g.filename = filename
	}
}

func WithIdRange(minId, maxId int) func(*IdGenerator) {
	return func(g *IdGenerator) {
		g.minId = minId
		g.maxId = maxId
	}
}

func WithIdCount(count int) func(*IdGenerator) {
	return func(g *IdGenerator) {
		g.idCount = count
	}
}

// Initialize loads Ids from file or generates new Ids and saves them to a file
func (g *IdGenerator) Initialize() error {
	if fileInfo, err := os.Stat(g.filename); os.IsNotExist(err) || fileInfo.Size() == 0 {
		// IdFile does not exist, generate new Ids and save them to the file
		g.GenerateRandIds()
		return g.SaveIdsToFile()
	}

	// IdFile exists, load Ids from the file
	return g.LoadIds(g.getCycleStartTime(0))
}

// LoadIds loads corresponding Ids into the in-memory map
func (g *IdGenerator) LoadIds(startTime time.Time) error {
	g.ids = make([]int, 0)
	g.idMap = make(map[int64]int)

	content, err := os.ReadFile(g.filename)
	if err != nil {
		return err
	}

	idsStr := bytes.Split(content, []byte("\n"))

	startTimeNum := startTime.UnixNano() / int64(g.unitDuration)
	for i := 0; i < g.idMapLength; i++ {
		key := startTimeNum + int64(i)
		id, _ := strconv.Atoi(string(idsStr[i]))
		g.idMap[key] = id
		g.ids = append(g.ids, id)
	}

	return nil
}

// GenerateRandIds generates randomly shuffled numeric Ids within the specified range
func (g *IdGenerator) GenerateRandIds() {
	count := g.idCount

	if g.minId > g.maxId {
		g.minId, g.maxId = g.maxId, g.minId
	}

	g.ids = make([]int, count)

	for i := g.minId; i < g.minId+count; i++ {
		g.ids[i-g.minId] = i
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

// GetId retrieves the Id corresponding to the specified seconds
func (g *IdGenerator) GetId(key int64) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	if id, ok := g.idMap[key]; ok {
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

func (g *IdGenerator) getCycleStartTime(offset int) time.Time {
	now := time.Now()
	cycleOffset := time.Duration(offset) * g.cycleDuration
	startOfCycleTime := now.Truncate(g.cycleDuration).Add(cycleOffset)
	return startOfCycleTime
}
