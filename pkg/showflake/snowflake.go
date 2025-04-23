package showflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	timestampBits = 41
	clusterIDBits = 6  // max 64
	nodeIDBits    = 4  // max 16
	sequenceBits  = 12 //  4096

	maxClusterID = -1 ^ (-1 << clusterIDBits) // 63
	maxNodeID    = -1 ^ (-1 << nodeIDBits)    // 15
	maxSequence  = -1 ^ (-1 << sequenceBits)  // 4095

	timeShift    = clusterIDBits + nodeIDBits + sequenceBits
	clusterShift = nodeIDBits + sequenceBits
	nodeShift    = sequenceBits

	defaultEpoch int64 = 1735678800000 // 2025-01-01
)

type Snowflake struct {
	mu sync.Mutex

	epoch         int64
	lastTimestamp int64
	sequence      int64
	clusterID     int64
	nodeID        int64
}

type Option func(*Snowflake)

func WithEpoch(epoch time.Time) Option {
	return func(s *Snowflake) {
		s.epoch = epoch.UnixMilli()
	}
}

func NewSnowflake(clusterID, nodeID int64, opts ...Option) (*Snowflake, error) {
	if clusterID < 0 || clusterID > maxClusterID {
		return nil, fmt.Errorf("cluster ID must be between 0 and %d", maxClusterID)
	}

	if nodeID < 0 || nodeID > maxNodeID {
		return nil, fmt.Errorf("node ID must be between 0 and %d", maxNodeID)
	}

	sf := &Snowflake{
		mu:        sync.Mutex{},
		nodeID:    nodeID,
		clusterID: clusterID,
		epoch:     defaultEpoch,
	}

	for _, opt := range opts {
		opt(sf)
	}

	return sf, nil
}

func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastTimestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = now

	id := ((now - s.epoch) << timeShift) |
		(s.clusterID << clusterShift) |
		(s.nodeID << nodeShift) |
		s.sequence

	return id
}
