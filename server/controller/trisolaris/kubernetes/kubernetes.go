/*
 * Copyright (c) 2022 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kubernetes

import (
	"sync"
	"time"

	"github.com/op/go-logging"
	"gorm.io/gorm"

	. "github.com/deepflowys/deepflow/server/controller/common"
	models "github.com/deepflowys/deepflow/server/controller/db/mysql"
	"github.com/deepflowys/deepflow/server/controller/model"
	"github.com/deepflowys/deepflow/server/controller/service"
	"github.com/deepflowys/deepflow/server/controller/trisolaris/config"
	"github.com/deepflowys/deepflow/server/controller/trisolaris/dbmgr"
)

var log = logging.MustGetLogger("trisolaris.kubernetes")

type KubernetesInfo struct {
	mutex             sync.RWMutex
	clusterIDToDomain map[string]string
	db                *gorm.DB
	cfg               *config.Config
}

func NewKubernetesInfo(db *gorm.DB, cfg *config.Config) *KubernetesInfo {
	return &KubernetesInfo{cfg: cfg, db: db}
}

func (k *KubernetesInfo) TimedRefreshClusterID() {
	ticker := time.NewTicker(time.Duration(60) * time.Second).C
	for {
		k.refresh()
		select {
		case <-ticker:
			k.refresh()
		}
	}
}

func (k *KubernetesInfo) refresh() {
	log.Infof("refresh cache cluster_id started")
	k.mutex.Lock()
	DomainMgr := dbmgr.DBMgr[models.Domain](k.db)
	dbDomains, _ := DomainMgr.GetBatchFromTypes([]int{KUBERNETES})
	k.clusterIDToDomain = make(map[string]string)
	for _, dbDomain := range dbDomains {
		k.clusterIDToDomain[dbDomain.ClusterID] = dbDomain.Lcuuid
	}
	k.mutex.Unlock()
	log.Infof("refresh cache cluster_id completed")
	return
}

func (k *KubernetesInfo) CacheClusterID(clusterID string) {
	log.Infof("start cache cluster_id (%s)", clusterID)
	k.mutex.Lock()
	_, ok := k.clusterIDToDomain[clusterID]
	if !ok {
		k.clusterIDToDomain[clusterID] = ""
		log.Infof("cache cluster_id (%s)", clusterID)
		go func() {
			for k.clusterIDToDomain[clusterID] == "" {
				domainLcuuid, err := k.createDomain(clusterID)
				if err != nil {
					log.Errorf("auto create domain failed: %v", err)
					time.Sleep(time.Second * 30)
				} else {
					k.clusterIDToDomain[clusterID] = domainLcuuid
				}
			}
		}()
	}
	k.mutex.Unlock()
	return
}

func (k *KubernetesInfo) createDomain(clusterID string) (domainLcuuid string, err error) {
	log.Infof("auto create domain (cluster_id: %s)", clusterID)
	azConMgr := dbmgr.DBMgr[models.AZControllerConnection](k.db)
	azConn, err := azConMgr.GetFromControllerIP(k.cfg.NodeIP)
	if err != nil {
		log.Errorf("controller (%s) az connection not in DB", k.cfg.NodeIP)
		return
	}
	domainConf := map[string]interface{}{
		"controller_ip":              k.cfg.NodeIP,
		"pod_net_ipv4_cidr_max_mask": 16,
		"pod_net_ipv6_cidr_max_mask": 64,
		"port_name_regex":            DEFAULT_PORT_NAME_REGEX,
		"region_uuid":                azConn.Region,
		"vtap_id":                    "",
	}
	domainCreate := model.DomainCreate{
		Name:                "k8s-" + clusterID,
		Type:                KUBERNETES,
		KubernetesClusterID: clusterID,
		ControllerIP:        k.cfg.NodeIP,
		Config:              domainConf,
	}
	domain, err := service.CreateDomain(domainCreate, nil)
	if err != nil {
		return
	}
	return domain.Lcuuid, nil
}
