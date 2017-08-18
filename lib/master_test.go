package lib

import (
	"log"
	"testing"
	"time"
)

// go version	: 1.8+
// author		: chapin
// email		: chengbin@lycam.tv

func TestMaster(t *testing.T) {
	m, err := NewMaster("lycamplus-etcd-go-test", "", "", []string{
		"http://54.222.215.0:2379",
		"http://54.222.170.116:2379",
		"http://54.222.214.254:2379",
	})
	if err != nil {
		t.Fatal(err)
	}
	for {
		log.Println("all ->", m.GetNodes())
		log.Println("all(strictly) ->", m.GetNodesStrictly())
		time.Sleep(time.Second * 2)
	}
}
