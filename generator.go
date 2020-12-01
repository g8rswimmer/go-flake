package flake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ID is the flake id 46 bits of time, 6 bits of worker and 12 bits of sequence
type ID uint64

// Decimal is the ID decimal string
func (i ID) Decimal() string {
	return fmt.Sprintf("%d", i)
}

func (i ID) String() string {
	t := uint64(i) >> 18
	w := uint64(i) >> 12
	w &= 0x3F
	s := uint64(i) & 0xFFF

	return fmt.Sprintf("%012X-%02X-%03X", t, w, s)
}

// ErrSequenceRollover is an error when the sequence number rolls over and the epoch time has not incremented
var ErrSequenceRollover = errors.New("flake generator sequence rollover")

const (
	// InvalidID is an invalid id
	InvalidID ID = 0
)

// Generate will create an id from the worker and sequence requested
func Generate(worker, seq uint64) (ID, error) {
	switch {
	case seq == 0 || seq >= 4096:
		return InvalidID, fmt.Errorf("flake generator: sequence can not be zero or over 4095 [%d]", seq)
	case worker == 0 || worker >= 64:
		return InvalidID, fmt.Errorf("flake generator: worker can not be zero or over 63 [%d]", worker)
	default:
	}

	epoch := uint64(time.Now().UnixNano() / int64(time.Millisecond))

	return generate(epoch, worker, seq), nil
}

func generate(epoch, worker, seq uint64) ID {
	id := epoch & 0x3FFFFFFFFFFF
	id <<= 6
	id |= uint64(worker & 0x3F)
	id <<= 12
	id |= uint64(seq) & 0xFFF
	return ID(id)
}

// New will create a new generator for the worker
func New(worker uint64) (*Generator, error) {
	switch {
	case worker == 0 || worker >= 64:
		return nil, fmt.Errorf("flake generator: worker can not be zero or over 63 [%d]", worker)
	default:
	}
	return &Generator{
		worker: worker,
		mutex:  sync.Mutex{},
	}, nil
}

// Generator will create a new id from the worker
type Generator struct {
	epoch  uint64
	seq    uint64
	worker uint64
	mutex  sync.Mutex
}

// Generate will create a new id
func (g *Generator) Generate() (ID, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	epoch := uint64(time.Now().UnixNano() / int64(time.Millisecond))
	if epoch != g.epoch {
		g.epoch = epoch
		g.seq = 0
	}
	g.seq++
	if g.seq >= 4096 {
		return InvalidID, ErrSequenceRollover
	}

	return generate(g.epoch, g.worker, g.seq), nil
}

// Wait will sleep for a millisecond, if there was an error on generate, this will help make sure that the sequence rollover will be handled.
func (g *Generator) Wait() {
	time.Sleep(time.Millisecond)
}
