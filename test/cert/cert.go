package main

import (
	"fmt"
	"crypto/rsa"
	"os"
	"io/ioutil"
	"encoding/pem"
	"crypto/x509"
	"time"
	"math/big"
	"crypto/x509/pkix"
	"crypto/rand"
	"crypto/tls"
	"net"
)

func main() {
	//generateCert("test.cer.pem")
	//cert := loadCert("test.cer.pem")
	//fmt.Printf("%v\n", *cert)

	go startTCPServerWithTLS()
	time.Sleep(time.Second * 5)
	go startTCPClientWithTLS()
	time.Sleep(time.Second * 30)
}
func startTCPClientWithTLS() {
	conn, err := tls.Dial("tcp", "127.0.0.1:1200", nil)
	checkErr(err)
	handleClient(conn)
}
func startTCPServerWithTLS() {
	cert, err := tls.LoadX509KeyPair("test.cer.pem", "pri.key")
	checkErr(err)
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	now := time.Now()
	config.Time = func() time.Time { return now }
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", ":1200", &config)
	checkErr(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}
func handleClient(conn net.Conn) {
	//do nothing
}
func loadCert(path string) *x509.Certificate {
	certCerFile, err := os.Open(path)
	checkErr(err)
	bb, err := ioutil.ReadAll(certCerFile)
	checkErr(err)
	certCerFile.Close()

	block, _ := pem.Decode(bb)

	// trim the bytes to actual length in call
	cert, err := x509.ParseCertificate(block.Bytes)
	checkErr(err)
	return cert;
}
func generateCert(path string) {
	random := rand.Reader
	privateKey := loadPEMKey("pri.key")
	now := time.Now()
	then := now.Add(60 * 60 * 24 * 365 * 1000 * 1000 * 1000) // one year
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "jan.newmarch.name",
			Organization: []string{"Jan Newmarch"},
		},
		NotBefore: now,
		NotAfter:  then,

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,

		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"jan.newmarch.name", "localhost"},
	}
	derBytes, err := x509.CreateCertificate(random, &template,
		&template, &(privateKey.PublicKey), privateKey)
	checkErr(err)

	block := &pem.Block{
		Type: "CERTIFICATE",
		Bytes: derBytes,
	}

	certCerFile, err := os.Create(path)
	checkErr(err)
	pem.Encode(certCerFile, block)
	certCerFile.Close()
}

func loadPEMKey(path string) *rsa.PrivateKey {
	file, err := os.Open(path)
	checkErr(err)
	defer file.Close()
	bb, err := ioutil.ReadAll(file)
	block, _ := pem.Decode(bb)
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey
}
func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(-1)
	}
}
