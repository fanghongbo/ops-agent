package g

import (
	"github.com/fanghongbo/ops-agent/common/rpc"
	"time"
)

var (
	HbsClient *rpc.SingleConnRpcClient
)

func InitHbsClient() {
	if config.Heartbeat != nil && config.Heartbeat.Enabled {
		HbsClient = &rpc.SingleConnRpcClient{
			RpcServer: config.Heartbeat.Addr,
			Timeout:   time.Duration(config.Heartbeat.Timeout) * time.Millisecond,
		}
	}
}
