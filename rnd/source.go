package rnd

import (
	crypt_rand "crypto/rand"
	"encoding/binary"
	"math/rand/v2"
	"sync"
)

type SafeSource struct {
	mu     sync.Mutex
	source rand.Source
}

func (s *SafeSource) Seed(seed uint64) {
	// Non-op
}

func (s *SafeSource) Uint64() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.source != nil {
		return s.source.Uint64()
	}
	panic("no random source available")
}

func NewChaCha8(seed *uint64) rand.Source {
	s := make([]byte, 32)
	if seed == nil {
		if _, err := crypt_rand.Read(s); err != nil {
			panic(err)
		}
	} else {
		seedBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(seedBytes, *seed)
		for i := 0; i < 4; i++ {
			copy(s[i*8:(i+1)*8], seedBytes)
		}
	}

	return rand.NewChaCha8([32]byte(s))
}

func NewSafeSource(source rand.Source) rand.Source {
	if source == nil {
		source = NewChaCha8(nil)
	}

	return &SafeSource{
		source: source,
	}
}

func ThreadSafeSource(src rand.Source) rand.Source {
	return &SafeSource{source: src}
}
