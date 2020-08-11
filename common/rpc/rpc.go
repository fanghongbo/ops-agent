package rpc

import (
	"errors"
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/net"
	"math"
	"net/rpc"
	"sync"
	"time"
)

type SingleConnRpcClient struct {
	sync.Mutex
	rpcClient *rpc.Client
	RpcServer string
	Timeout   time.Duration
}

func (u *SingleConnRpcClient) Close() {
	if u.rpcClient != nil {
		_ = u.rpcClient.Close()
		u.rpcClient = nil
	}
}

func (u *SingleConnRpcClient) serverConn() error {
	var (
		err   error
		retry int = 1
	)

	if u.rpcClient != nil {
		return nil
	}

	for {
		if u.rpcClient != nil {
			return nil
		}

		u.rpcClient, err = net.JsonRpcClient("tcp", u.RpcServer, u.Timeout)
		if err != nil {
			dlog.Errorf("dial %s fail: %v", u.RpcServer, err)
			if retry > 3 {
				return err
			}
			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
			retry++
			continue
		}
		return err
	}
}

func (u *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {
	var (
		err     error
		timeout time.Duration
		done    chan error
	)

	u.Lock()
	defer u.Unlock()

	err = u.serverConn()
	if err != nil {
		return err
	}

	timeout = time.Duration(10 * time.Second)
	done = make(chan error, 1)

	go func() {
		err := u.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		dlog.Errorf("rpc call timeout %v => %v", u.rpcClient, u.RpcServer)
		u.Close()
		return errors.New(u.RpcServer + " rpc call timeout")
	case err := <-done:
		if err != nil {
			u.Close()
			return err
		}
	}

	return nil
}

// code == 0 => success
// code == 1 => bad request
type SimpleRpcResponse struct {
	Code int `json:"code"`
}

func (u *SimpleRpcResponse) String() string {
	return fmt.Sprintf("<Code: %d>", u.Code)
}

type NullRpcRequest struct {
}
