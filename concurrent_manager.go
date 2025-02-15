package golib

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
)

type cmT struct {
	workers  int
	sem      chan bool
	wg       sync.WaitGroup
	mu       sync.Mutex
	ordered  map[int][]OrderedFunc
	runners  int
	finished bool
	errs     []error
}

type RunFunc func(runId int) error
type OrderedFunc func() error
type SyncedFunc func() error

// ConcurrentManager is used to concurrently execute a func
// in n workers. If n is <= 0, the number of CPUs is used.
func ConcurrentManager(n int) *cmT {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	cm := cmT{
		sem:     make(chan bool, n),
		workers: n,
		ordered: map[int][]OrderedFunc{},
	}
	return &cm
}

// Workers returns the number of workers
func (cm *cmT) Workers() int {
	return cm.workers
}

// Synced executes f under a mutex. Use this to sync state which belongs to all
// workers.
func (cm *cmT) Synced(f SyncedFunc) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return f()
}

// Synced executes f under a mutex. Use this to sync state which belongs to all
// workers.
func (cm *cmT) Ordered(runId int, f OrderedFunc) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.ordered[runId] = append(cm.ordered[runId], f)
}

// Run executes f. Waits for a worker to be available. This should be called in a Go routine. Panics
// are caught and can be retrieved when using cm.Error(). If the manager is in error state,
// Run returns an error
func (cm *cmT) Run(f RunFunc) (runID int) {
	cm.sem <- true
	cm.wg.Add(1)
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.finished {
		panic("cm.Run cannot be called after cm.Wait")
	}
	cm.runners++
	cm.ordered[cm.runners] = []OrderedFunc{}
	runId := cm.runners
	go func() {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("panic: %v", r)
				}
			}
			if err != nil {
				cm.mu.Lock()
				cm.errs = append(cm.errs, err)
				cm.mu.Unlock()
			}
			<-cm.sem
			cm.wg.Done()
		}()
		err = f(runId)
	}()
	return cm.runners
}

// Errors returns the accumulated errors so far.
func (cm *cmT) Errors() []error {
	return cm.errs
}

// Wait waits until all workers have finished. It returns an error if any of the
// func "Run" calls returned an error or panicked. In case that no errors
// occurred, Wait also call all functions registered with "Ordered" in the order
// of "Run" execution. If an "Ordered" function returns an error, execution
// stops and Wait returns that error. "Run" cannot be called after wait has been
// called.
func (cm *cmT) Wait() error {
	cm.mu.Lock()
	cm.finished = true
	cm.mu.Unlock()
	cm.wg.Wait()
	if len(cm.errs) > 0 {
		errS := []string{}
		for _, err := range cm.errs {
			errS = append(errS, err.Error())
		}
		return errors.New(strings.Join(errS, ", "))
	}
	for runId := 1; runId <= cm.runners; runId++ {
		for _, f := range cm.ordered[runId] {
			err := f()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Reset resets the manager to the
func (cm *cmT) Reset() {
	cm.ordered = map[int][]OrderedFunc{}
	cm.finished = false
	cm.errs = []error{}
	cm.runners = 0
}
