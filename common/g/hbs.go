package g

import (
	"github.com/fanghongbo/ops-agent/common/rpc"
	"time"
)

var (
	HbsClient *rpc.SingleConnRpcClient
)

func InitHbsClient() {
	if Conf().Heartbeat.Enabled {
		HbsClient = &rpc.SingleConnRpcClient{
			RpcServer: Conf().Heartbeat.Addr,
			Timeout:   time.Duration(Conf().Heartbeat.Timeout) * time.Millisecond,
		}
	}
}
