package metric

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterMetric(t *testing.T) {
	m := CounterMetric{Map: sync.Map{}}
	count := 1000
	expexted := 0
	for i := 0; i < count; i++ {
		expexted += i
	}
	expexted += 1

	wg := sync.WaitGroup{}
	wg.Add(count * 2)
	for i := 0; i < count; i++ {
		i := i
		go func() {
			defer wg.Done()
			m.Increment("test", int64(i))
		}()
		go func() {
			defer wg.Done()
			m.Increment("test1", int64(i))
		}()
	}
	wg.Wait()
	v := m.Increment("test", 1)
	v1 := m.Increment("test1", 1)

	assert.Equal(t, int64(expexted), v)
	assert.Equal(t, int64(expexted), v1)

}
