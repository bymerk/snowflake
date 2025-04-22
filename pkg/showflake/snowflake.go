package showflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch         = 1609459200000 // 1 Jan 2021 (UTC)
	timestampBits = 41
	nodeIDBits    = 10
	sequenceBits  = 12

	maxNodeID   = -1 ^ (-1 << nodeIDBits)
	maxSequence = -1 ^ (-1 << sequenceBits)

	nodeIDShift    = sequenceBits
	timestampShift = sequenceBits + nodeIDBits
)

type Snowflake struct {
	mu       sync.Mutex
	epoch    int64
	lastTime int64
	sequence int64
	nodeID   int64
}

type Option func(*Snowflake)

func WithEpoch(epoch time.Time) Option {
	return func(s *Snowflake) {
		s.epoch = epoch.UnixMilli()
	}
}

func NewSnowflake(nodeID int64, opts ...Option) (*Snowflake, error) {
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, fmt.Errorf("node ID must be between 0 and %d", maxNodeID)
	}

	sf := &Snowflake{
		epoch:  epoch,
		nodeID: nodeID,
		mu:     sync.Mutex{},
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

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	id := ((now - epoch) << timestampShift) |
		(s.nodeID << nodeIDShift) |
		s.sequence

	return id
}
