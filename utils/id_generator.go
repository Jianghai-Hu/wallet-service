package utils

import (
	"context"
	"sync"
	"time"

	"github.com/golang/glog"

	"jianghai-hu/wallet-service/internal/common"
)

// Constants for ID generation
const (
	epoch int64 = 1733011200000 // 2024-12-01 00:00:00 UTC
	// Adjusted bits allocation to maximize the range
	// timeBits      uint  = 41
	machineIDBits uint  = 10
	sequenceBits  uint  = 12
	maxMachineID  int32 = (1 << machineIDBits) - 1
	maxSequence   int32 = (1 << sequenceBits) - 1
)

var globalIDGenerator *IDGenerator

// IDGenerator structure for generating unique IDs
type IDGenerator struct {
	mu            sync.Mutex
	lastTimestamp int64
	sequence      int32
	machineID     int32
}

// InitIDGenerator initializes the global ID generator
func InitIDGenerator(ctx context.Context, machineID int32) {
	if machineID < 0 || machineID > maxMachineID {
		glog.FatalContext(ctx, "machineID out of range")
	}

	globalIDGenerator = &IDGenerator{
		lastTimestamp: 0,
		sequence:      0,
		machineID:     machineID,
	}
}

// GetIDGenerator returns the global ID generator
func GetIDGenerator(ctx context.Context) *IDGenerator {
	if globalIDGenerator == nil {
		glog.FatalContext(ctx, "globalIDGenerator is nil")
	}

	return globalIDGenerator
}

// Generate generates a unique ID
func (g *IDGenerator) Generate() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := currentTimestamp()

	// Handle clock moving backwards
	if timestamp < g.lastTimestamp {
		return 0, NewMyError(common.Constant_ERROR_SERVICE_INTERNAL, "clock moved backwards, refusing to generate id")
	}

	// Sequence increment or wait for next timestamp
	if timestamp == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			timestamp = waitUntilNextMillie(g.lastTimestamp)
		}
	} else {
		g.sequence = 0
	}

	g.lastTimestamp = timestamp

	// Assemble the ID
	id := ((timestamp - epoch) << (machineIDBits + sequenceBits)) |
		(int64(g.machineID) << sequenceBits) |
		int64(g.sequence)

	// Ensure ID does not exceed MaxInt64
	if id < 0 {
		return 0, NewMyError(common.Constant_ERROR_SERVICE_INTERNAL, "generated id exceeds int64 range")
	}

	return id, nil
}

// currentTimestamp returns the current timestamp in microseconds
func currentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// waitUntilNextMicro waits until the next microsecond
func waitUntilNextMillie(lastTimestamp int64) int64 {
	timestamp := currentTimestamp()
	for timestamp <= lastTimestamp {
		timestamp = currentTimestamp()
	}

	return timestamp
}
