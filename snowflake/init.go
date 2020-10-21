package snowflake

/*
github.com/twitter/snowflake in golang

id =>  timestamp retain center worker sequence
          40      4       5      5      10
*/

import (
	"fmt"
	"sync"
	"time"
)

var defSnowFlake *snowFlake

func init() {
	defSnowFlake, _ = New(0, 0)
}

// Gen 生成ID
func Gen() (uint64, error) {
	return defSnowFlake.Gen()
}

const (
	nano = 1000 * 1000

	timestampBits = 40                         // timestamp
	maxtimestamp  = -1 ^ (-1 << timestampBits) // timestamp mask
	retainedBits  = 4
	maxRetain     = -1 ^ (-1 << retainedBits)
	centerBits    = 5
	maxCenter     = -1 ^ (-1 << centerBits) // center mask
	workerBits    = 5
	maxWorker     = -1 ^ (-1 << workerBits)   // worker mask
	sequenceBits  = 10                        // sequence
	maxSequence   = -1 ^ (-1 << sequenceBits) // sequence mask
)

var (
	since  int64                 = time.Date(2017, 5, 1, 0, 0, 0, 0, time.Local).UnixNano() / nano
	poolMu sync.RWMutex          = sync.RWMutex{}
	pool   map[uint64]*snowFlake = make(map[uint64]*snowFlake)
)

type snowFlake struct {
	lastTimestamp uint64
	retain        uint32
	center        uint32
	worker        uint32
	sequence      uint32
	lock          sync.Mutex
}

// New 创建一个SnowFlake
func New(centerID uint32, workerID uint32) (*snowFlake, error) {
	if centerID > maxCenter {
		return nil, fmt.Errorf("CenterID %v is invalid", centerID)
	} else if workerID > maxWorker {
		return nil, fmt.Errorf("WorkerID %v is invalid", workerID)
	}
	return &snowFlake{
		worker: workerID,
		center: centerID,
	}, nil
}

// Gen 生成一个ID
func (sf *snowFlake) Gen() (uint64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	ts := timestamp()
	if ts == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & maxSequence
		if sf.sequence == 0 {
			ts = tilNextMillis(ts)
		}
	} else {
		sf.sequence = 0
	}

	if ts < sf.lastTimestamp {
		return 0, fmt.Errorf("Invalid timestamp: %v - precedes %v", ts, sf)
	}
	sf.lastTimestamp = ts
	return sf.uint64(), nil
}

func (sf *snowFlake) uint64() uint64 {
	return (sf.lastTimestamp << (retainedBits + centerBits + workerBits + sequenceBits)) |
		(uint64(sf.retain) << (centerBits + workerBits + sequenceBits)) |
		(uint64(sf.center) << (workerBits + sequenceBits)) |
		(uint64(sf.worker) << sequenceBits) |
		uint64(sf.sequence)
}
func timestamp() uint64 {
	return uint64(time.Now().UnixNano()/nano - since)
}

func tilNextMillis(ts uint64) uint64 {
	i := timestamp()
	for i <= ts {
		i = timestamp()
	}
	return i
}
