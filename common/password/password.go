package password

import (
	"crypto/md5"
	"encoding/hex"
)

const salt = "this is my mall"

// Hash 与历史种子/登录一致：MD5(password + salt)
func Hash(plain string) string {
	h := md5.New()
	h.Write([]byte(plain + salt))
	return hex.EncodeToString(h.Sum(nil))
}
