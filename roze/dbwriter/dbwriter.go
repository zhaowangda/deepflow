package dbwriter

import (
	logging "github.com/op/go-logging"

	"gitlab.x.lan/yunshan/droplet-libs/ckdb"
	"gitlab.x.lan/yunshan/droplet-libs/zerodoc"
	"gitlab.x.lan/yunshan/droplet/pkg/ckwriter"
	"gitlab.x.lan/yunshan/droplet/roze/config"
	"gitlab.x.lan/yunshan/droplet/roze/msg"
)

var log = logging.MustGetLogger("roze.dbwriter")

const (
	CACHE_SIZE = 10240
)

type DbWriter struct {
	ckwriters []*ckwriter.CKWriter
}

func NewDbWriter(primaryAddr, secondaryAddr, user, password string, replicaEnabled bool, ckWriterCfg config.CKWriterConfig) (*DbWriter, error) {
	ckwriters := []*ckwriter.CKWriter{}
	engine := ckdb.MergeTree
	if replicaEnabled {
		engine = ckdb.ReplicatedMergeTree
	}
	tables := zerodoc.GetMetricsTables(engine)
	for _, table := range tables {
		counterName := "metrics_1m"
		if table.ID >= uint8(zerodoc.VTAP_FLOW_1S) {
			counterName = "metrics_1s"
		}
		ckwriter, err := ckwriter.NewCKWriter(primaryAddr, secondaryAddr, user, password, counterName, table, replicaEnabled,
			ckWriterCfg.QueueCount, ckWriterCfg.QueueSize, ckWriterCfg.BatchSize, ckWriterCfg.FlushTimeout)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		ckwriter.Run()
		ckwriters = append(ckwriters, ckwriter)
	}

	return &DbWriter{
		ckwriters: ckwriters,
	}, nil
}

func (w *DbWriter) Put(items ...interface{}) error {
	caches := [zerodoc.VTAP_DB_ID_MAX][]interface{}{}
	for i := range caches {
		caches[i] = make([]interface{}, 0, CACHE_SIZE)
	}
	for _, item := range items {
		doc, ok := item.(*msg.RozeDocument)
		if !ok {
			log.Warningf("receive wrong type data %v", item)
			continue
		}
		id, err := doc.TableID()
		if err != nil {
			log.Warningf("doc table id not found. %v", doc)
			continue
		}
		caches[id] = append(caches[id], doc)
	}

	for i, cache := range caches {
		if len(cache) > 0 {
			w.ckwriters[i].Put(cache...)
		}
	}
	return nil
}

func (w *DbWriter) Close() {
	for _, ckwriter := range w.ckwriters {
		ckwriter.Close()
	}
}
