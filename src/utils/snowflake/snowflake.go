package snowflake

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake ID 生成器
// Twitter 的雪花算法实现
// 64位ID: 1位符号位(0) + 41位时间戳 + 10位机器ID + 12位序列号

const (
	epoch          = int64(1609459200000)              // 起始时间戳 (2021-01-01 00:00:00 UTC)
	timestampBits  = uint(41)                          // 时间戳占用位数
	machineIDBits  = uint(10)                          // 机器ID占用位数
	sequenceBits   = uint(12)                          // 序列号占用位数
	maxMachineID   = int64(-1 ^ (-1 << machineIDBits)) // 最大机器ID (1023)
	maxSequence    = int64(-1 ^ (-1 << sequenceBits))  // 最大序列号 (4095)
	machineIDShift = sequenceBits                      // 机器ID左移位数
	timestampShift = sequenceBits + machineIDBits      // 时间戳左移位数
)

// 雪花ID生成器
type Generator struct {
	mu        sync.Mutex
	timestamp int64 // 上次生成ID的时间戳
	machineID int64 // 机器ID
	sequence  int64 // 序列号
}

var (
	//  默认的全局生成器
	DefaultGenerator *Generator
	initOnce         sync.Once
)

// 创建新的雪花ID生成器
func NewGenerator(machineID int64) (*Generator, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("机器ID必须在 0 到 %d 之间", maxMachineID)
	}

	return &Generator{
		timestamp: 0,
		machineID: machineID,
		sequence:  0,
	}, nil
}

// 初始化默认生成器
func Init(machineID int64) error {
	var err error
	initOnce.Do(func() {
		DefaultGenerator, err = NewGenerator(machineID)
	})
	return err
}

// 生成下一个ID
func (g *Generator) NextID() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	// 如果当前时间小于上次生成ID的时间，说明时钟回退了
	if now < g.timestamp {
		return 0, fmt.Errorf("时钟回退，拒绝生成ID")
	}

	// 如果在同一毫秒内
	if now == g.timestamp {
		// 序列号递增
		g.sequence = (g.sequence + 1) & maxSequence

		// 序列号溢出，等待下一毫秒
		if g.sequence == 0 {
			for now <= g.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 不同毫秒，序列号重置为0
		g.sequence = 0
	}

	g.timestamp = now

	// 组装ID
	id := ((now - epoch) << timestampShift) |
		(g.machineID << machineIDShift) |
		g.sequence

	return id, nil
}

// 使用默认生成器生成ID
func Generate() (int64, error) {
	if DefaultGenerator == nil {
		return 0, fmt.Errorf("雪花ID生成器未初始化，请先调用 Init()")
	}
	return DefaultGenerator.NextID()
}

// 生成ID，如果失败则 panic
func MustGenerate() int64 {
	id, err := Generate()
	if err != nil {
		panic(err)
	}
	return id
}

// 解析雪花ID，返回时间戳、机器ID和序列号
func ParseID(id int64) (timestamp int64, machineID int64, sequence int64) {
	timestamp = (id >> timestampShift) + epoch
	machineID = (id >> machineIDShift) & maxMachineID
	sequence = id & maxSequence
	return
}

// 从雪花ID中提取时间戳
func GetTimestamp(id int64) time.Time {
	timestamp, _, _ := ParseID(id)
	return time.UnixMilli(timestamp)
}
