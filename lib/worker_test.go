package lib

import (
	"flag"
	"testing"
	"time"
)

// go version	: 1.8+
// author		: chapin
// email		: chengbin@lycam.tv

func TestWorker(t *testing.T) {
	appName := flag.String("app", "lycamplus-etcd-go-test", "app name")
	serviceName := flag.String("service", "default", "service name")
	version := flag.String("version", "1.0", "version")

	flag.Parse()

	etcdRequestInfo := EtcdRequestInfo{
		AppName:     *appName,
		ServiceName: *serviceName,
		Host:        "127.0.0.1",
		Port:        1337,
		Version:     *version,
	}
	w, err := NewWorker(etcdRequestInfo, []string{
		"http://54.222.215.0:2379",
		"http://54.222.170.116:2379",
		"http://54.222.214.254:2379",
	})
	if err != nil {
		t.Fatal(err)
	}
	w.Register()

	go func() {
		time.Sleep(time.Second * 50)
		w.Unregister()
	}()

	for {
		t.Log("isActive -> ", w.IsActive())
		t.Log("isStop ->", w.IsStop())
		time.Sleep(time.Second * 2)

		if w.IsStop() {
			return
		}
	}
}
