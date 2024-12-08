package utils

import (
	"context"
	"github.com/golang/glog"
	"jianghai-hu/wallet-service/internal/common"
	"sync"
	"time"
)

const (
	epoch int64 = 1733011200000 // 2024-12-01 00:00:00 UTC
	// timeBits      uint  = 17
	machineIDBits uint  = 10
	sequenceBits  uint  = 5
	maxMachineID  int32 = (1 << machineIDBits) - 1
	maxSequence   int32 = (1 << sequenceBits) - 1
	maxID         int32 = (1 << 31) - 1
)

var globalIDGenerator *idGenerator

type idGenerator struct {
	mu            sync.Mutex
	lastTimestamp int64
	sequence      int32
	machineID     int32
}

func InitIDGenerator(ctx context.Context, machineID int32) {
	if machineID < 0 || machineID > maxMachineID {
		glog.FatalContext(ctx, "machineID out of range")
	}
	globalIDGenerator = &idGenerator{
		lastTimestamp: 0,
		sequence:      0,
		machineID:     machineID,
	}
}

func GetIDGenerator(ctx context.Context) *idGenerator {
	if globalIDGenerator == nil {
		glog.FatalContext(ctx, "globalIDGenerator is nil")
	}
	return globalIDGenerator
}

func (g *idGenerator) Generate() (int32, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := currentTimestamp()

	if timestamp < g.lastTimestamp {
		return 0, NewMyError(common.Constant_ERROR_SERVICE_INTERNAL, "clock moved backwards, refusing to generate id")
	}

	if timestamp == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			timestamp = waitUntilNextMillis(g.lastTimestamp)
		}
	} else {
		g.sequence = 0
	}

	g.lastTimestamp = timestamp

	id := ((timestamp - epoch) << (machineIDBits + sequenceBits)) |
		(int64(g.machineID) << sequenceBits) |
		int64(g.sequence)

	if id > int64(maxID) {
		return 0, NewMyError(common.Constant_ERROR_SERVICE_INTERNAL, "generated id exceeds int32 range")
	}

	return int32(id), nil
}

func currentTimestamp() int64 {
	return time.Now().UnixMilli()
}

func waitUntilNextMillis(lastTimestamp int64) int64 {
	timestamp := currentTimestamp()
	for timestamp <= lastTimestamp {
		timestamp = currentTimestamp()
	}
	return timestamp
}
