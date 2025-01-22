package golib

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func TestConcurrentManager(t *testing.T) {
	cm := ConcurrentManager(5)
	unordered := []int{}
	ordered := []int{}
	for i := range 10 {
		func(i int) {
			cm.Run(func(runId int) error {
				time.Sleep(time.Duration(rand.Intn(100)*10) * time.Millisecond)
				cm.Synced(func() error {
					unordered = append(unordered, i)
					return nil
				})
				cm.Ordered(runId, func() error {
					ordered = append(ordered, i)
					return nil
				})
				return nil
			})
		}(i)
	}
	cm.Wait()
	Pln("exec order: %v ordered: %v", unordered, ordered)
	if !assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, ordered) {
		return
	}
}

func TestConcurrentManager2(t *testing.T) {
	cm := ConcurrentManager(5)
	unordered := []int{}
	ordered := []int{}
	for i := range 10 {
		func(i int) {
			cm.Run(func(runId int) error {
				time.Sleep(time.Duration(rand.Intn(100)*10) * time.Millisecond)
				cm.Synced(func() error {
					unordered = append(unordered, i)
					return nil
				})
				if i == 3 || i == 5 {
					return fmt.Errorf("i %d caused an error", i)
				}
				cm.Ordered(runId, func() error {
					ordered = append(ordered, i)
					return nil
				})
				return nil
			})
		}(i)
	}
	err := cm.Wait()
	if !assert.Error(t, err) {
		return
	}
	Pln("error: %s", err.Error())
}
