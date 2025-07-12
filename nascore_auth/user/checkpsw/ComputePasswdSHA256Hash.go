package checkpsw

import (
	"crypto/sha256"
	"encoding/hex"
)

func ComputePasswdSHA256Hash(password, salt string) string { // computeSHA256Hash 计算 SHA256 哈希值
	hasher := sha256.New()
	hasher.Write([]byte("NasCore.eu.org" + password + salt + "NasCore.eu.org"))
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}
