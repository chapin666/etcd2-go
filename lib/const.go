package lib

import "time"

// go version	: 1.8+
// author		: chapin
// email		: chengbin@lycam.tv

const (

	// RootKey .
	RootKey = "rpc"

	// HeartBeatInterval .
	HeartBeatInterval = time.Second * 5

	// TTLTime .
	TTLTime = time.Second * 10

	// TimeOut .
	TimeOut = time.Second * 5

	// DefaultVersion .
	DefaultVersion = "1.0"

	// DefaultServiceName .
	DefaultServiceName = "default"
)

// EtcdRequestInfo .
type EtcdRequestInfo struct {
	AppName     string
	ServiceName string
	Version     string
	Host        string
	Port        uint64
}

// MemoryInfo .
type MemoryInfo struct {
	Rss   uint64 `json:"rss"`
	Vsize uint64 `json:"vsize"`
}

// EtcdTransData .
type EtcdTransData struct {
	DockerID   string      `json:"cid"`
	Host       string      `json:"host"`
	Port       uint64      `json:"port"`
	Weight     uint64      `json:"weight"`
	CPU        float64     `json:"cpu"`
	Memory     float32     `json:"memory"`
	MemoryInfo *MemoryInfo `json:"memoryInfo"`
}
