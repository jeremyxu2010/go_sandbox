package main

import (
	"golang.org/x/crypto/blowfish"
	"os"
	"fmt"
	"encoding/base64"
)

func main() {
	key := []byte("my key")
	cipher, err := blowfish.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
	}
	enc := make([]byte, 255)

	cipher.Encrypt(enc, []byte("hello\n\n\n"))

	fmt.Println(base64.StdEncoding.EncodeToString(enc))

	decrypt := make([]byte, 8)
	cipher.Decrypt(decrypt, enc)
	fmt.Println(string(decrypt))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
