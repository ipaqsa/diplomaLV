package service

import (
	"fmt"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"ingress/pkg"
	"sync/atomic"
)

func (s *PoolT) initPool() {
	for _, address := range pkg.Config.Brokers {
		back := createBackend(address)
		s.backends = append(s.backends, back)
	}
}
func createBackend(address string) *Backend {
	return &Backend{
		Address: address,
	}
}

func (s *PoolT) nextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

func (s *PoolT) nextPeer() *Backend {
	next := s.nextIndex()
	l := len(s.backends) + next
	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].isAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

func (b *Backend) isAlive() bool {
	if b == nil || b.Address == "" {
		return false
	}
	fmt.Printf("Alive - %s\n", b.Address)
	opt := rpc.CreateOptions(false, 0, nil, nil)
	pack := packUtils.CreatePack("health", "alive?")
	send, err := rpc.Send(b.Address, "ServerBroker.Health", pack, opt)
	if err != nil {
		return false
	}
	if send.Body.Data != "alive" {
		return false
	}
	return true
}

func (s *PoolT) Prox(prox *packUtils.Package) (*packUtils.Package, error) {
	opt := rpc.CreateOptions(false, 0, nil, nil)
	var back *Backend
	for back == nil {
		back = s.nextPeer()
	}
	ans, err := rpc.Send(back.Address, "ServerBroker.PutMail", prox, opt)
	if err != nil {
		return nil, err
	}
	return ans, nil
}
