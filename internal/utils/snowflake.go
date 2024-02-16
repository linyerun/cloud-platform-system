package utils

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

const (
	workerBits     uint8  = 10                      // 工作机器ID，可以部署在1024个节点上
	workerMax      uint16 = -1 ^ (-1 << workerBits) // 节点ID的最大值
	seqBits        uint8  = 12
	seqMax         uint16 = -1 ^ (-1 << seqBits)
	timestampShift        = workerBits + seqBits
	workerIdShift         = seqBits
)

var (
	once       = new(sync.Once)
	snowFlaker *snowflake
)

// 雪花ID的组成(有最高位到最低位):
//		1bit: 符号位，始终为0
//	   41bit: 时间戳, 精确到毫秒
//	   10bit: 工作机器ID(Max: 2^10 - 1)
//     12bit: 序列化(Max: 2^12 - 1)

// 每个机器一个序列ID，每个机器每毫秒最多生成(2^12 - 1)个雪花ID

type snowflake struct {
	mu sync.Mutex

	// 作用: timestamp相等就就增加seqId, 不一样就把seqId只为0
	timestamp int64
	seqId     uint16

	// 最后的timestamp要减去下面这个值
	epoch int64

	// 机器编号
	workerId uint16
}

func newSnowFlake(workerId uint16) (*snowflake, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}
	return &snowflake{timestamp: 0, seqId: 0, epoch: time.Now().UnixMilli(), workerId: workerId}, nil
}

func (s *snowflake) getId() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	if now != s.timestamp {
		s.seqId = 0 // 不相等就置为0即可
	} else {
		s.seqId = (s.seqId + 1) & seqMax
		if s.seqId == 0 {
			// 如果当前工作节点在1毫秒内生成的ID已经超过上限 需要等待1毫秒再继续生成
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	}
	s.timestamp = now
	ret := ((s.timestamp-s.epoch)&(2<<41-1))<<timestampShift | int64(s.workerId<<workerIdShift) | int64(s.seqId)
	return ret
}

func GetSnowFlakeIdAndBase64() string {
	buf := make([]byte, 8)
	binary.PutVarint(buf, GetSnowFlakeId())
	return base64.URLEncoding.EncodeToString(buf)
}

func GetSnowFlakeId() int64 {
	once.Do(func() {
		flake, err := newSnowFlake(0)
		if err != nil {
			panic(err)
		}
		snowFlaker = flake
	})
	return snowFlaker.getId()
}
