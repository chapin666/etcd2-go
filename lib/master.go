package lib

// go version	: 1.8+
// author		: chapin
// email		: chengbin@lycam.tv

import (
	"context"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
)

// Master .
type Master struct {
	sync.RWMutex
	kapi   client.KeysAPI
	key    string
	nodes  map[string]string
	active bool
}

// NewMaster create a master
func NewMaster(appName string, serviceName string, version string, endpoints []string) (*Master, error) {
	cfg := client.Config{
		Endpoints:               endpoints,
		HeaderTimeoutPerRequest: TimeOut,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	master := &Master{
		kapi:   client.NewKeysAPI(c),
		key:    getServiceKey(appName, serviceName, version),
		nodes:  make(map[string]string),
		active: true,
	}

	master.fetch()

	go master.watch()

	return master, err
}

// GetNodesStrictly .
func (m *Master) GetNodesStrictly() map[string]string {
	if !m.active {
		return nil
	}
	return m.GetNodes()
}

// GetNodes .
func (m *Master) GetNodes() map[string]string {
	m.RLock()
	defer m.RUnlock()
	return m.nodes
}

// addNode .
func (m *Master) addNode(node, extInfo string) {
	m.Lock()
	defer m.Unlock()
	node = strings.TrimLeft(node, m.key)
	m.nodes[node] = extInfo
}

// delNode .
func (m *Master) delNode(node string) {
	m.Lock()
	defer m.Unlock()
	node = strings.TrimLeft(node, m.key)
	delete(m.nodes, node)
}

// watch .
func (m *Master) watch() {
	watcher := m.kapi.Watcher(m.key, &client.WatcherOptions{
		Recursive: true,
	})
	for {
		resp, err := watcher.Next(context.Background())
		if err != nil {
			m.active = false
			continue
		}
		m.active = true

		switch resp.Action {
		case "set", "update":
			m.addNode(resp.Node.Key, resp.Node.Value)
			break
		case "expire", "delete":
			m.delNode(resp.Node.Key)
			break
		default:
			log.Info("watchme!!!", "resp ->", resp)
		}
	}
}

// fetch.
func (m *Master) fetch() error {
	resp, err := m.kapi.Get(context.Background(), m.key, nil)
	if err != nil {
		return nil
	}
	if resp.Node.Dir {
		for _, v := range resp.Node.Nodes {
			m.addNode(v.Key, v.Value)
		}
	}
	return err
}
