package gopublic

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"time"
)

var (
	randomStringGenerator *RandomStringGenerator = NewRandomStringGenerator(time.Now().UnixNano())
)

// RandonString 生成随机字符串
func RandonString(lenth int) string {
	return randomStringGenerator.NewString(lenth)
}

func RandonStringWithPrefix(prefix string, length int) string {
	return prefix + "-" + RandonString(length)
}

func NewRandomStringGenerator(seed int64) *RandomStringGenerator {
	return &RandomStringGenerator{randSource: rand.New(rand.NewSource(seed))}
}

type RandomStringGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

func (gen *RandomStringGenerator) NewString(lenth int) string {
	gen.Lock()
	defer gen.Unlock()
	data := make([]byte, lenth)
	_, _ = gen.randSource.Read(data)
	return hex.EncodeToString(data)
}
