package artnet

import (
	"fmt"
	"net"
	"sync"

	"github.com/jsimonetti/go-artnet"
	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

var (
	globalArtnet     *Artnet
	globalArtnetOnce sync.Once
)

// GetArtnet returns a global singleton instance.
// First call initializes it.
func GetArtnet(artsubnet, nodeName string) *Artnet {
	fmt.Printf("GetArtnet called with %s %s\n", artsubnet, nodeName)
	globalArtnetOnce.Do(func() {
		globalArtnet = NewArtnet(artsubnet, nodeName)
		globalArtnet.Start()
	})
	return globalArtnet
}

type Artnet struct {
	_, cidrnet *net.IPNet
	broadcast  *net.IP

	node *artnet.Node

	dmxData map[int][512]byte // universe -> dmx data
	lock    sync.RWMutex
}

func NewArtnet(artsubnet string, nodeName string) *Artnet {
	a := &Artnet{}
	_, a.cidrnet, _ = net.ParseCIDR(artsubnet)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("error getting ips: %s\n", err)
	}

	var ip net.IP

	for _, addr := range addrs {
		ip = addr.(*net.IPNet).IP
		fmt.Printf("Found IP: %s\n", ip.String())
		if a.cidrnet.Contains(ip) {
			break
		}
	}

	fmt.Printf("Using IP: %s\n", ip.String())

	// Calculate broadcast address for the subnet
	broadcast := make(net.IP, len(a.cidrnet.IP))
	copy(broadcast, a.cidrnet.IP)
	for i := range broadcast {
		broadcast[i] |= ^a.cidrnet.Mask[i]
	}
	fmt.Printf("Using broadcast address: %s\n", broadcast.String())

	a.broadcast = &broadcast
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:6454", broadcast.String()))
	opti := artnet.NodeBroadcastAddress(*addr)

	log := artnet.NewDefaultLogger()
	a.node = artnet.NewNode(nodeName, code.StVisual, ip, log, opti)

	a.dmxData = make(map[int][512]byte)

	// register callback for any packet
	a.node.RegisterCallback(code.OpDMX, func(p packet.ArtNetPacket) {
		dmx := p.(*packet.ArtDMXPacket)
		fmt.Printf("DMX universe %d len %d\n", dmx.SubUni, dmx.Length)
		a.storeDMXData(int(dmx.SubUni), dmx.Data)
	})

	return a
}

// Thread-safe read
func (a *Artnet) GetDMXData(universe int) [512]byte {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// return copy
	return a.dmxData[universe]
}

// Thread-safe write
func (a *Artnet) storeDMXData(universe int, data [512]byte) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.dmxData[universe] = data
}

func (a *Artnet) Start() {
	if err := a.node.Start(); err != nil {
		fmt.Printf("error starting artnet node: %s\n", err)
		return
	}
}

func (a *Artnet) Stop() {
	a.node.Stop()
}
