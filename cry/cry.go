package cry

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func GetMd5(str string ) string {
	h := md5.New()
	h.Write([]byte(str )) // 需要加密的字符串为
	cipherStr := h.Sum(nil)
	return fmt.Sprintf("%s\n", hex.EncodeToString(cipherStr))
}


func GetMd52(str string) string{
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}