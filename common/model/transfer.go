package model

import "fmt"

type TransferResponse struct {
	Message string
	Total   int
	Invalid int
	Latency int64
}

func (u *TransferResponse) String() string {
	return fmt.Sprintf(
		"<Total=%v, Invalid:%v, Latency=%vms, Message:%s>",
		u.Total,
		u.Invalid,
		u.Latency,
		u.Message,
	)
}
