package g

import (
	"bytes"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/model"
	"github.com/fanghongbo/ops-agent/common/rpc"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	TransferClientsLock *sync.RWMutex                       = new(sync.RWMutex)
	TransferClients     map[string]*rpc.SingleConnRpcClient = map[string]*rpc.SingleConnRpcClient{}
)

func initTransferClient(addr string) *rpc.SingleConnRpcClient {
	var c *rpc.SingleConnRpcClient = &rpc.SingleConnRpcClient{
		RpcServer: addr,
		Timeout:   time.Duration(config.Transfer.Timeout) * time.Millisecond,
	}

	TransferClientsLock.Lock()
	defer TransferClientsLock.Unlock()
	TransferClients[addr] = c

	return c
}

func updateMetrics(c *rpc.SingleConnRpcClient, metrics []*model.MetricValue, resp *model.TransferResponse) bool {
	err := c.Call("Transfer.Update", metrics, resp)
	if err != nil {
		dlog.Error("call Transfer.Update fail:", c, err)
		return false
	}
	return true
}

func getTransferClient(addr string) *rpc.SingleConnRpcClient {
	TransferClientsLock.RLock()
	defer TransferClientsLock.RUnlock()

	if c, ok := TransferClients[addr]; ok {
		return c
	}
	return nil
}

func SendMetrics(metrics []*model.MetricValue, resp *model.TransferResponse) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(Conf().Transfer.Addrs)) {
		addr := Conf().Transfer.Addrs[i]

		c := getTransferClient(addr)
		if c == nil {
			c = initTransferClient(addr)
		}

		if updateMetrics(c, metrics, resp) {
			break
		}
	}
}

func SendToTransfer(metrics []*model.MetricValue) {
	var (
		dt             map[string]string
		buf            bytes.Buffer
		defaultTagList []string
		defaultTagStr  string
		resp           model.TransferResponse
	)

	if len(metrics) == 0 {
		return
	}

	dt = Conf().DefaultTags
	if len(dt) > 0 {
		defaultTagList = []string{}
		for k, v := range dt {
			buf.Reset()
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v)
			defaultTagList = append(defaultTagList, buf.String())
		}

		defaultTagStr = strings.Join(defaultTagList, ",")

		for i, x := range metrics {
			buf.Reset()
			if x.Tags == "" {
				metrics[i].Tags = defaultTagStr
			} else {
				buf.WriteString(metrics[i].Tags)
				buf.WriteString(",")
				buf.WriteString(defaultTagStr)
				metrics[i].Tags = buf.String()
			}
		}
	}

	debug := Conf().Debug
	if debug {
		dlog.Infof("=> <Total=%d> %v", len(metrics), metrics[0])
	}

	SendMetrics(metrics, &resp)

	if debug {
		dlog.Infof("<= %v", &resp)
	}
}
