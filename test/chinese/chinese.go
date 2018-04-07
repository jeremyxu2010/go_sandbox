package main

import (
	"os"
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"io"
)

func main() {
	fileName := "test.txt"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		_, err := os.Create(fileName)
		checkErr(err)
	}

	f, err := os.Open(fileName)
	checkErr(err)

	writer := transform.NewWriter(f, simplifiedchinese.GB18030.NewEncoder())
	io.WriteString(writer, "百度一下，你就知道")
	writer.Close()

	reader := transform.NewReader(f, simplifiedchinese.GB18030.NewDecoder())
	bb, err := ioutil.ReadAll(reader)
	checkErr(err)
	fmt.Println(string(bb))

	f.Close()

}
func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
