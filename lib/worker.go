package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/process"
)

// go version	: 1.8+
// author		: chapin
// email		: chengbin@lycam.tv

// Worker struct.
type Worker struct {
	kapi      client.KeysAPI
	key       string
	extraInfo string
	active    bool
	stop      bool
}

// NewWorker create a worker
func NewWorker(etcdRequestInfo EtcdRequestInfo, endpoints []string) (*Worker, error) {
	cfg := client.Config{
		Endpoints:               endpoints,
		HeaderTimeoutPerRequest: TimeOut,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	extraData := getPayload(etcdRequestInfo.Host, etcdRequestInfo.Port)
	out, err := json.Marshal(extraData)
	if err != nil {
		panic(err)
	}
	extraInfo := string(out)

	key := getServiceKey(etcdRequestInfo.AppName, etcdRequestInfo.ServiceName, etcdRequestInfo.Version)
	ep := fmt.Sprintf("%s:%d", etcdRequestInfo.Host, etcdRequestInfo.Port)
	sha1 := sha1Encode(ep)

	worker := &Worker{
		kapi:      client.NewKeysAPI(c),
		key:       fmt.Sprintf("%s/%s", key, sha1),
		extraInfo: extraInfo,
		active:    false,
		stop:      false,
	}
	return worker, nil
}

// Register .
func (w *Worker) Register() {
	w.heartbeat()
	go w.heartbeatPeriod()
}

// Unregister .
func (w *Worker) Unregister() {
	w.stop = true
}

// IsActive .
func (w *Worker) IsActive() bool {
	return w.active
}

// IsStop .
func (w *Worker) IsStop() bool {
	return w.stop
}

func (w *Worker) heartbeatPeriod() {
	for !w.stop {
		w.heartbeat()
		time.Sleep(HeartBeatInterval)
	}
}

func (w *Worker) heartbeat() error {
	_, err := w.kapi.Set(context.Background(), w.key, w.extraInfo, &client.SetOptions{
		TTL: TTLTime,
	})
	w.active = err != nil
	return err
}

func getPayload(hostIP string, port uint64) *EtcdTransData {
	checkPid := os.Getpid()
	ret, _ := process.NewProcess(int32(checkPid))
	m, _ := ret.MemoryInfo()
	mPercent, _ := ret.MemoryPercent()
	cPercent, _ := ret.Percent(0)
	dockerID := getDockerID()

	memoryInfo := &MemoryInfo{
		Rss:   m.RSS,
		Vsize: m.VMS,
	}

	return &EtcdTransData{
		DockerID:   dockerID,
		Host:       hostIP,
		Port:       port,
		Weight:     100,
		CPU:        cPercent,
		Memory:     mPercent,
		MemoryInfo: memoryInfo,
	}
}

func getDockerID() string {
	d, _ := docker.GetDockerStat()
	if len(d) >= 1 {
		conatiner := d[0]
		return conatiner.ContainerID
	}

	return "0"
}
