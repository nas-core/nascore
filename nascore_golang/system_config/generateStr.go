package system_config

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func GenerateStr(typeInt int) string {
	baseStr, err := os.Hostname()
	tmpHash := "NasCore.eu.org"
	if err != nil {
		baseStr = tmpHash
	}
	h := md5.New()
	switch typeInt {
	case 1:
		io.WriteString(h, baseStr)
		baseStr = fmt.Sprintf("%x", h.Sum(nil))
	case 2:
		io.WriteString(h, baseStr+tmpHash)
		baseStr = fmt.Sprintf("%x", h.Sum(nil))
	}
	return baseStr + "-" + tmpHash
}
