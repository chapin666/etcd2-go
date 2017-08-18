package lib

import (
	"crypto/sha1"
	"fmt"
)

// go version	: 1.8+
// author		: chapin
// email		: chengbin@lycam.tv

func getServiceKey(appName, serviceName, version string) string {
	if serviceName == "" {
		serviceName = DefaultServiceName
	}
	if version == "" {
		version = DefaultVersion
	}

	return fmt.Sprintf("/%s/%s/%s/%s", RootKey, appName, serviceName, version)
}

func sha1Encode(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
