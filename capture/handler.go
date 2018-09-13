package capture

import (
	"runtime"
	"sync"
	"time"

	"gitlab.x.lan/yunshan/droplet-libs/datatype"
	"gitlab.x.lan/yunshan/droplet-libs/dedup"
	"gitlab.x.lan/yunshan/droplet-libs/queue"
)

type Timestamp = time.Duration
type RawPacket = []byte

type PacketHandler interface {
	Handle(Timestamp, RawPacket)
}

type MetaPacketBlock = [1024]datatype.MetaPacket

type DataHandler struct {
	sync.Pool

	block       *MetaPacketBlock
	blockCursor int
	ip          datatype.IPv4Int
	queue       queue.MultiQueueWriter
}

func (h *DataHandler) preAlloc() *datatype.MetaPacket {
	metaPacket := &h.block[h.blockCursor]
	metaPacket.InPort = uint32(datatype.PACKET_SOURCE_ISP)
	metaPacket.Exporter = h.ip
	return metaPacket
}

func (h *DataHandler) confirmAlloc() {
	h.blockCursor++
	if h.blockCursor >= len(*h.block) {
		h.block = h.Get().(*MetaPacketBlock)
		*h.block = MetaPacketBlock{} // 回收的block包含脏数据，因此需要重新清空
		h.blockCursor = 0
	}
}

func (h *DataHandler) Handle(timestamp Timestamp, packet RawPacket) {
	metaPacket := h.preAlloc()
	metaPacket.Timestamp = timestamp
	metaPacket.PacketLen = uint16(len(packet))
	if !metaPacket.Parse(packet) {
		return
	}
	h.confirmAlloc()
	hash := metaPacket.InPort + metaPacket.IpSrc + metaPacket.IpDst +
		uint32(metaPacket.Protocol) + uint32(metaPacket.PortSrc) + uint32(metaPacket.PortDst)
	h.queue.Put(queue.HashKey(hash), metaPacket)
}

func (h *DataHandler) Init() *DataHandler {
	gc := func(b *MetaPacketBlock) { h.Put(b) }
	h.Pool.New = func() interface{} {
		block := new(MetaPacketBlock)
		runtime.SetFinalizer(block, gc)
		return block
	}
	h.block = new(MetaPacketBlock)
	return h
}

type TapHandler DataHandler

func (h *TapHandler) Handle(timestamp Timestamp, packet RawPacket) {
	metaPacket := (*DataHandler)(h).preAlloc()
	metaPacket.Timestamp = timestamp
	metaPacket.PacketLen = uint16(len(packet))
	tunnel := datatype.TunnelInfo{}
	if offset := tunnel.Decapsulate(packet); offset > 0 {
		packet = packet[offset:]
		metaPacket.Tunnel = &tunnel
	}
	if dedup.Lookup(packet, timestamp) {
		return
	}
	if !metaPacket.Parse(packet) {
		return
	}
	(*DataHandler)(h).confirmAlloc()
	hash := metaPacket.InPort + metaPacket.IpSrc + metaPacket.IpDst +
		uint32(metaPacket.Protocol) + uint32(metaPacket.PortSrc) + uint32(metaPacket.PortDst)
	h.queue.Put(queue.HashKey(hash), metaPacket)
}
