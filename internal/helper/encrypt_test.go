package helper

import (
	"fmt"
	"testing"
)

const key = "O4dUc6NexJ8myBFh1jwSgAYkLn7ifHGr"

func TestEncrypt(t *testing.T) {
	pwd := "123456"
	encrypwd, err := Encrypt(pwd, key)

	fmt.Printf("pwd: %s, err:%v", encrypwd, err)

}

func TestDecrypt(t *testing.T) {
	encrypwd := "UNTsZF8BKfXEai72alMpcbDLCJOpAA=="
	pwd, err := Decrypt(encrypwd, key)
	fmt.Printf("pwd: %s, err:%v", pwd, err)
}

func TestHashPassword(t *testing.T) {
	pwd := "123456"
	fmt.Printf("pwd: %s, err:%v", HashPassword(pwd, key), nil)
}
