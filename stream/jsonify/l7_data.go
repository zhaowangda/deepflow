package jsonify

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"gitlab.x.lan/yunshan/droplet-libs/ckdb"
	"gitlab.x.lan/yunshan/droplet-libs/datatype"
	"gitlab.x.lan/yunshan/droplet-libs/grpc"
	"gitlab.x.lan/yunshan/droplet-libs/pool"
	"gitlab.x.lan/yunshan/droplet-libs/utils"
	pf "gitlab.x.lan/yunshan/droplet/stream/platformdata"
)

type L7Base struct {
	// 知识图谱
	KnowledgeGraph

	// 网络层
	IP40   uint32 `json:"ip4_0"`
	IP41   uint32 `json:"ip4_1"`
	IP60   net.IP `json:"ip6_0"`
	IP61   net.IP `json:"ip6_1"`
	IsIPv4 bool   `json:"is_ipv4"`

	// 传输层
	ClientPort uint16 `json:"client_port"`
	ServerPort uint16 `json:"server_port"`

	// 流信息
	FlowIDStr uint64 `json:"flow_id_str"`
	TapType   uint16 `json:"tap_type"`
	TapPort   uint32 `json:"tap_port"` // 显示为固定八个字符的16进制如'01234567'
	VtapID    uint16 `json:"vtap_id"`
	Timestamp uint64 `json:"timestamp"` // us
	Time      uint32 `json:"time"`      // 秒，用来淘汰过期数据
}

func L7BaseColumns() []*ckdb.Column {
	columns := []*ckdb.Column{}
	// 知识图谱
	columns = append(columns, KnowledgeGraphColumns...)
	columns = append(columns,
		// 网络层
		ckdb.NewColumn("ip4_0", ckdb.IPv4),
		ckdb.NewColumn("ip4_1", ckdb.IPv4),
		ckdb.NewColumn("ip6_0", ckdb.IPv6),
		ckdb.NewColumn("ip6_1", ckdb.IPv6),
		ckdb.NewColumn("is_ipv4", ckdb.UInt8),

		// 传输层
		ckdb.NewColumn("client_port", ckdb.UInt16),
		ckdb.NewColumn("server_port", ckdb.UInt16),

		// 流信息
		ckdb.NewColumn("flow_id_str", ckdb.UInt64),
		ckdb.NewColumn("tap_type", ckdb.UInt16),
		ckdb.NewColumn("tap_port", ckdb.UInt32),
		ckdb.NewColumn("vtap_id", ckdb.UInt16),
		ckdb.NewColumn("timestamp", ckdb.UInt64).SetCodec(ckdb.CodecDoubleDelta),
		ckdb.NewColumn("time", ckdb.DateTime),
	)

	return columns
}

func (f *L7Base) WriteBlock(block *ckdb.Block) error {
	if err := f.KnowledgeGraph.WriteBlock(block); err != nil {
		return err
	}

	if err := block.WriteUInt32(f.IP40); err != nil {
		return err
	}
	if err := block.WriteUInt32(f.IP41); err != nil {
		return err
	}
	if len(f.IP60) == 0 {
		f.IP60 = net.IPv6zero
	}
	if err := block.WriteIP(f.IP60); err != nil {
		return err
	}
	if len(f.IP61) == 0 {
		f.IP61 = net.IPv6zero
	}
	if err := block.WriteIP(f.IP61); err != nil {
		return err
	}

	if err := block.WriteBool(f.IsIPv4); err != nil {
		return err
	}

	if err := block.WriteUInt16(f.ClientPort); err != nil {
		return err
	}
	if err := block.WriteUInt16(f.ServerPort); err != nil {
		return err
	}

	if err := block.WriteUInt64(f.FlowIDStr); err != nil {
		return err
	}
	if err := block.WriteUInt16(f.TapType); err != nil {
		return err
	}
	if err := block.WriteUInt32(f.TapPort); err != nil {
		return err
	}
	if err := block.WriteUInt16(f.VtapID); err != nil {
		return err
	}
	if err := block.WriteUInt64(f.Timestamp); err != nil {
		return err
	}
	if err := block.WriteUInt32(f.Time); err != nil {
		return err
	}

	return nil
}

type HTTPLogger struct {
	pool.ReferenceCount
	_id uint64

	L7Base

	// http应用层
	Type         uint8  `json:"type"` // 0: request  1: response 2: session
	Version      uint8  `json:"version"`
	Method       string `json:"method,omitempty"`
	ClientIP4    uint32 `json:"client_ip4,omitempty"`
	ClientIP6    net.IP `json:"client_ip6,omitempty"`
	ClientIsIPv4 bool   `json:"client_is_ipv4"`
	Host         string `json:"host,omitempty"`
	Path         string `json:"path,omitempty"`
	StreamID     uint32 `json:"stream_id,omitempty"`
	TraceID      string `json:"trace_id,omitempty"`
	StatusCode   uint16 `json:"status_code,omitempty"`

	// 指标量
	ContentLength int64  `json:"content_length"`
	Duration      uint64 `json:"duration,omitempty"` // us
}

func HTTPLoggerColumns() []*ckdb.Column {
	httpColumns := []*ckdb.Column{}
	httpColumns = append(httpColumns, ckdb.NewColumn("_id", ckdb.UInt64).SetCodec(ckdb.CodecDoubleDelta))
	httpColumns = append(httpColumns, L7BaseColumns()...)
	httpColumns = append(httpColumns,
		// 应用层HTTP
		ckdb.NewColumn("type", ckdb.UInt8),
		ckdb.NewColumn("version", ckdb.UInt8),
		ckdb.NewColumn("method", ckdb.String),
		ckdb.NewColumn("client_ip4", ckdb.IPv4),
		ckdb.NewColumn("client_ip6", ckdb.IPv6),
		ckdb.NewColumn("client_is_ipv4", ckdb.UInt8),

		ckdb.NewColumn("host", ckdb.String),
		ckdb.NewColumn("path", ckdb.String),
		ckdb.NewColumn("stream_id", ckdb.UInt32),
		ckdb.NewColumn("trace_id", ckdb.String),
		ckdb.NewColumn("status_code", ckdb.UInt16),

		// 指标量
		ckdb.NewColumn("content_length", ckdb.Int64),
		ckdb.NewColumn("duration", ckdb.UInt64),
	)
	return httpColumns
}

func (h *HTTPLogger) WriteBlock(block *ckdb.Block) error {
	index := 0
	err := block.WriteUInt64(h._id)
	if err != nil {
		return err
	}
	index++

	if err := h.L7Base.WriteBlock(block); err != nil {
		return nil
	}

	if err := block.WriteUInt8(h.Type); err != nil {
		return err
	}
	if err := block.WriteUInt8(h.Version); err != nil {
		return err
	}
	if err := block.WriteString(h.Method); err != nil {
		return err
	}
	if err := block.WriteUInt32(h.ClientIP4); err != nil {
		return err
	}
	if len(h.ClientIP6) == 0 {
		h.ClientIP6 = net.IPv6zero
	}
	if err := block.WriteIP(h.ClientIP6); err != nil {
		return err
	}
	if err := block.WriteBool(h.ClientIsIPv4); err != nil {
		return err
	}

	if err := block.WriteString(h.Host); err != nil {
		return err
	}
	if err := block.WriteString(h.Path); err != nil {
		return err
	}
	if err := block.WriteUInt32(h.StreamID); err != nil {
		return err
	}
	if err := block.WriteString(h.TraceID); err != nil {
		return err
	}
	if err := block.WriteUInt16(h.StatusCode); err != nil {
		return err
	}
	if err := block.WriteInt64(h.ContentLength); err != nil {
		return err
	}
	if err := block.WriteUInt64(h.Duration); err != nil {
		return err
	}

	return nil
}

func parseIP(ipStr string) (uint32, net.IP, bool) {
	var ip4 uint32
	var ip6 net.IP
	isIPv4 := true

	ip := net.ParseIP(ipStr)
	if ip != nil {
		to4 := ip.To4()
		if to4 != nil {
			isIPv4 = true
			ip4 = utils.IpToUint32(to4)
		} else {
			isIPv4 = false
			ip6 = ip
		}
	}

	return ip4, ip6, isIPv4
}

func parseVersion(str string) uint8 {
	// 对于1.0,1.1 解析为 10, 11
	rmDot := strings.ReplaceAll(str, ".", "")
	v, _ := strconv.Atoi(rmDot)
	// 对于 2，需要解析为20
	if v < 10 {
		v = v * 10
	}
	return uint8(v)
}

func (h *HTTPLogger) Fill(l *datatype.AppProtoLogsData) {
	h.L7Base.Fill(l)
	if l.Proto == datatype.PROTO_HTTP {
		if httpInfo, ok := l.Detail.(*datatype.HTTPInfo); ok {
			h.Version = parseVersion(httpInfo.Version)
			h.Method = strings.ToUpper(httpInfo.Method)
			h.ClientIP4, h.ClientIP6, h.ClientIsIPv4 = parseIP(httpInfo.ClientIP)
			h.Host = httpInfo.Host
			h.Path = httpInfo.Path
			h.StreamID = httpInfo.StreamID
			h.TraceID = httpInfo.TraceID
			h.ContentLength = int64(httpInfo.ContentLength)
		}
	}
	h.Type = uint8(l.MsgType)
	h.StatusCode = l.Code
	h.Duration = uint64(l.RRT / time.Microsecond)
}

func (h *HTTPLogger) Release() {
	ReleaseHTTPLogger(h)
}

func (h *HTTPLogger) EndTime() time.Duration {
	return time.Duration(h.Timestamp) * time.Microsecond
}

func (h *HTTPLogger) String() string {
	return fmt.Sprintf("HTTP: %+v\n", *h)
}

type DNSLogger struct {
	pool.ReferenceCount
	_id uint64

	L7Base

	// DNS应用层
	Type       uint8  `json:"type"` // 0: request  1: response 2: session
	ID         uint16 `json:"id"`
	DomainName string `json:"domain_name,omitempty"`
	QueryType  uint16 `json:"query_type,omitempty"`
	AnswerCode uint16 `json:"answer_code"`
	AnswerAddr string `json:"answer_addr,omitempty"`

	// 指标量
	Duration uint64 `json:"duration,omitempty"` // us
}

func DNSLoggerColumns() []*ckdb.Column {
	dnsColumns := []*ckdb.Column{}
	dnsColumns = append(dnsColumns, ckdb.NewColumn("_id", ckdb.UInt64).SetCodec(ckdb.CodecDoubleDelta))
	dnsColumns = append(dnsColumns, L7BaseColumns()...)
	dnsColumns = append(dnsColumns,
		// 应用层DNS
		ckdb.NewColumn("type", ckdb.UInt8).SetComment("0: request 1: response 2: session"),
		ckdb.NewColumn("id", ckdb.UInt16),
		ckdb.NewColumn("domain_name", ckdb.String),
		ckdb.NewColumn("query_type", ckdb.UInt16),
		ckdb.NewColumn("answer_code", ckdb.UInt16),
		ckdb.NewColumn("answer_addr", ckdb.String),

		// 指标量
		ckdb.NewColumn("duration", ckdb.UInt64),
	)
	return dnsColumns
}

func (d *DNSLogger) WriteBlock(block *ckdb.Block) error {
	if err := block.WriteUInt64(d._id); err != nil {
		return err
	}

	if err := d.L7Base.WriteBlock(block); err != nil {
		return nil
	}

	if err := block.WriteUInt8(d.Type); err != nil {
		return err
	}
	if err := block.WriteUInt16(d.ID); err != nil {
		return err
	}
	if err := block.WriteString(d.DomainName); err != nil {
		return err
	}
	if err := block.WriteUInt16(d.QueryType); err != nil {
		return err
	}
	if err := block.WriteUInt16(d.AnswerCode); err != nil {
		return err
	}
	if err := block.WriteString(d.AnswerAddr); err != nil {
		return err
	}

	if err := block.WriteUInt64(d.Duration); err != nil {
		return err
	}
	return nil
}

func (d *DNSLogger) Fill(l *datatype.AppProtoLogsData) {
	d.L7Base.Fill(l)

	// 应用层DNS信息
	if l.Proto == datatype.PROTO_DNS {
		if dnsInfo, ok := l.Detail.(*datatype.DNSInfo); ok {
			d.ID = dnsInfo.TransID
			d.DomainName = dnsInfo.QueryName
			d.QueryType = dnsInfo.QueryType
			d.AnswerAddr = dnsInfo.Answers
		}
	}
	d.Type = uint8(l.MsgType)
	d.AnswerCode = l.Code

	// 指标量
	d.Duration = uint64(l.RRT / time.Microsecond)
}

func (d *DNSLogger) Release() {
	ReleaseDNSLogger(d)
}

func (d *DNSLogger) EndTime() time.Duration {
	return time.Duration(d.Timestamp) * time.Microsecond
}

func (d *DNSLogger) String() string {
	return fmt.Sprintf("DNS: %+v\n", *d)
}

func (b *L7Base) Fill(l *datatype.AppProtoLogsData) {
	// 网络层
	if l.IsIPv6 {
		b.IsIPv4 = false
		b.IP60 = l.IP6Src[:]
		b.IP61 = l.IP6Dst[:]
	} else {
		b.IsIPv4 = true
		b.IP40 = l.IPSrc
		b.IP41 = l.IPDst
	}

	// 传输层
	b.ClientPort = l.PortSrc
	b.ServerPort = l.PortDst

	// 知识图谱
	b.KnowledgeGraph.FillL7(l)

	// 流信息
	b.FlowIDStr = l.FlowId
	b.TapType = l.TapType
	b.TapPort = l.TapPort
	b.VtapID = l.VtapId
	b.Timestamp = uint64(l.Timestamp / time.Microsecond)
	b.Time = uint32(l.Timestamp / time.Second)
}

func (k *KnowledgeGraph) FillL7(l *datatype.AppProtoLogsData) {
	var info0, info1 *grpc.Info
	l3EpcID0, l3EpcID1 := l.L3EpcIDSrc, l.L3EpcIDDst

	if l.IsIPv6 {
		info0, info1 = pf.PlatformData.QueryIPV6InfosPair(int16(l3EpcID0), net.IP(l.IP6Src[:]), int16(l3EpcID1), net.IP(l.IP6Dst[:]))
	} else {
		info0, info1 = pf.PlatformData.QueryIPV4InfosPair(int16(l3EpcID0), uint32(l.IPSrc), int16(l3EpcID1), uint32(l.IPDst))
	}

	if info0 != nil {
		k.RegionID0 = uint16(info0.RegionID)
		k.AZID0 = uint16(info0.AZID)
		k.HostID0 = uint16(info0.HostID)
		k.L3DeviceType0 = uint8(info0.DeviceType)
		k.L3DeviceID0 = info0.DeviceID
		k.PodNodeID0 = info0.PodNodeID
		k.PodNSID0 = uint16(info0.PodNSID)
		k.PodGroupID0 = info0.PodGroupID
		k.PodID0 = info0.PodID
		k.PodClusterID0 = uint16(info0.PodClusterID)
		k.SubnetID0 = uint16(info0.SubnetID)
	}
	if info1 != nil {
		k.RegionID1 = uint16(info1.RegionID)
		k.AZID1 = uint16(info1.AZID)
		k.HostID1 = uint16(info1.HostID)
		k.L3DeviceType1 = uint8(info1.DeviceType)
		k.L3DeviceID1 = info1.DeviceID
		k.PodNodeID1 = info1.PodNodeID
		k.PodNSID1 = uint16(info1.PodNSID)
		k.PodGroupID1 = info1.PodGroupID
		k.PodID1 = info1.PodID
		k.PodClusterID1 = uint16(info1.PodClusterID)
		k.SubnetID1 = uint16(info1.SubnetID)
	}
	k.L3EpcID0, k.L3EpcID1 = l3EpcID0, l3EpcID1
}

var poolHTTPLogger = pool.NewLockFreePool(func() interface{} {
	return new(HTTPLogger)
})

func AcquireHTTPLogger() *HTTPLogger {
	l := poolHTTPLogger.Get().(*HTTPLogger)
	l.ReferenceCount.Reset()
	return l
}

func ReleaseHTTPLogger(l *HTTPLogger) {
	if l == nil {
		return
	}
	if l.SubReferenceCount() {
		return
	}
	*l = HTTPLogger{}
	poolHTTPLogger.Put(l)
}

var L7HTTPCounter uint32

func ProtoLogToHTTPLogger(l *datatype.AppProtoLogsData, shardID int) *HTTPLogger {
	h := AcquireHTTPLogger()
	h._id = genID(uint32(l.Timestamp/time.Microsecond), &L7HTTPCounter, shardID)
	h.Fill(l)
	return h
}

var poolDNSLogger = pool.NewLockFreePool(func() interface{} {
	return new(DNSLogger)
})

func AcquireDNSLogger() *DNSLogger {
	l := poolDNSLogger.Get().(*DNSLogger)
	l.ReferenceCount.Reset()
	return l
}

func ReleaseDNSLogger(l *DNSLogger) {
	if l == nil {
		return
	}
	if l.SubReferenceCount() {
		return
	}
	*l = DNSLogger{}
	poolDNSLogger.Put(l)
}

var L7DNSCounter uint32

func ProtoLogToDNSLogger(l *datatype.AppProtoLogsData, shardID int) *DNSLogger {
	h := AcquireDNSLogger()
	h._id = genID(uint32(l.Timestamp/time.Microsecond), &L7DNSCounter, shardID)
	h.Fill(l)
	return h
}
