package clients

import (
	"sync"

	"github.com/bits-and-blooms/bitset"
	"github.com/jin06/mercury/pkg/mqtt"
)

type packetIDManager struct {
	mu      sync.Mutex
	usedIDs *bitset.BitSet
	maxID   mqtt.PacketID
}

func NewPacketIDManager() *packetIDManager {
	return &packetIDManager{
		usedIDs: bitset.New(65536), // 16-bit Packet ID range
		maxID:   65535,
	}
}

func (manager *packetIDManager) Get() mqtt.PacketID {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for id := uint(1); id <= uint(manager.maxID); id++ {
		if !manager.usedIDs.Test(id) {
			manager.usedIDs.Set(id)
			return mqtt.PacketID(id)
		}
	}
	return 0 // No available PacketID
}

func (manager *packetIDManager) Put(id mqtt.PacketID) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.usedIDs.Clear(uint(id))
}
