package main

import (
	"crypto/rsa"
	"crypto/rand"
	"os"
	"fmt"
	"encoding/pem"
	"crypto/x509"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	checkErr(err)

	publicKey := &(privateKey.PublicKey)

	savePEMKey("pri.key", privateKey)
	savePEMPublicKey("pub.key", publicKey)

	privateKey2 := loadPEMKey("pri.key")
	fmt.Printf("%v\n", *privateKey2)
	publicKey2 := loadPEMPublicKey("pub.key")
	fmt.Printf("%v\n", *publicKey2)

}
func savePEMKey(path string, key *rsa.PrivateKey) {
	file, err := os.Create(path)
	checkErr(err)
	defer file.Close()
	block := &pem.Block{
		Type : "RSA PRIVATE KEY",
		Bytes : x509.MarshalPKCS1PrivateKey(key),
	}
	pem.Encode(file, block)
}
func savePEMPublicKey(path string, key *rsa.PublicKey) {
	file, err := os.Create(path)
	checkErr(err)
	defer file.Close()
	publicKey, err := ssh.NewPublicKey(key)
	checkErr(err)
	file.Write(ssh.MarshalAuthorizedKey(publicKey))
}
func loadPEMKey(path string) *rsa.PrivateKey {
	file, err := os.Open(path)
	checkErr(err)
	defer file.Close()
	bb, err := ioutil.ReadAll(file)
	block, _ := pem.Decode(bb)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	checkErr(err)
	return privateKey
}
func loadPEMPublicKey(path string) *rsa.PublicKey{
	file, err := os.Open(path)
	checkErr(err)
	defer file.Close()
	bb, err := ioutil.ReadAll(file)
	checkErr(err)
	pkey, _, _, _, err:= ssh.ParseAuthorizedKey(bb)
	checkErr(err)
	if pkey, ok := pkey.(ssh.CryptoPublicKey); ok {
		publicKey := pkey.CryptoPublicKey().(*rsa.PublicKey)
		return publicKey
	}
	return nil
}
func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(-1)
	}
}