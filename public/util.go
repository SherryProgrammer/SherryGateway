package public

import (
	"crypto/sha256"
	"fmt"
)

func GenSaltPassword(salt, password string) string { //该函数的作用是将密码和盐值作为输入，生成一个哈希字符串作为密码的安全存储形式。
	s1 := sha256.New()
	s1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}
