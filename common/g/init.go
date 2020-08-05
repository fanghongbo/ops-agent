package g

import (
	"context"
	"github.com/fanghongbo/dlog"
)

func InitAll() {
	InitConfig()    // 初始化配置文件
	InitRuntime()   // 初始化运行环境
	InitLocalIp()   // 初始化本地ip
	InitLog()       // 初始化日志
	InitHbsClient() // 初始化hbs rpc客户端
}

func Shutdown(ctx context.Context) error {
	defer ctx.Done()

	// 关闭hbs rpc连接
	if HbsClient != nil {
		HbsClient.Close()
	}

	// 关闭transfer rpc连接
	for _, rpcClient := range TransferClients {
		rpcClient.Close()
	}

	// 关闭日志，刷新缓存
	dlog.Close()

	return nil
}
