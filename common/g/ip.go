package g

import (
	"fmt"
	"github.com/fanghongbo/dlog"
	"net"
	"os"
	"strings"
	"time"
)

var LocalIp string

func InitLocalIp() {
	if config.Heartbeat != nil && config.Heartbeat.Enabled {
		conn, err := net.DialTimeout("tcp", config.Heartbeat.Addr, time.Second*10)
		if err != nil {
			dlog.Errorf("connect to %s err: %s", config.Heartbeat.Addr, err)
		} else {
			defer func() {
				_ = conn.Close()
			}()
			LocalIp = strings.Split(conn.LocalAddr().String(), ":")[0]
		}
	} else {
		if LocalIp, err := GetLocalIP(); err != nil {
			dlog.Errorf("get local ip fail %s", err)
		} else {
			dlog.Infof("local ip found: %s", LocalIp)
		}
	}

	if LocalIp == "" {
		dlog.Fatal("init local ip failed")
	} else {
		dlog.Infof("local ip found: %s", LocalIp)
	}
}

func LocalHostname() (string, error) {
	var (
		hostname string
		err      error
	)

	hostname, err = os.Hostname()
	if err != nil {
		return hostname, err
	}
	return hostname, nil
}

func GetLocalIP() (string, error) {
	var (
		addresses []net.Addr
		err       error
	)

	addresses, err = net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addresses {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("cannot get local ip address")
}

func Hostname() (string, error) {
	var (
		hostname string
		err      error
	)

	hostname = config.Hostname
	if hostname != "" {
		// use hostname in configuration
		return strings.TrimSpace(hostname), nil
	}

	hostname, err = LocalHostname()
	if err != nil {
		dlog.Errorf("get system hostname err: %s", err)
	}
	return strings.TrimSpace(hostname), err
}

func IP() string {
	var ip string

	ip = config.IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIp) > 0 {
		ip = LocalIp
	}

	return strings.TrimSpace(ip)
}
