package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	proxyURI, dstURI, authEncrypted, err := parseArgs()
	if err != nil {
		log.Println("usage: http-proxy-for-ssh <proxyhost> <proxyport> <dsthost> <dstport> [authfile]")
		log.Fatalln(err)
	}

	request, err := http.NewRequest("CONNECT", dstURI, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Host = dstURI
	if authEncrypted != "" {
		request.Header.Add("Proxy-Authorization", "Basic "+authEncrypted)
	}
	request.Header.Add("Host", dstURI)
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Proxy-Connection", "keep-alive")

	conn, err := net.Dial("tcp", proxyURI)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	if err = request.Write(conn); err != nil {
		log.Fatalln(err)
	}

	br := bufio.NewReader(conn)
	response, err := http.ReadResponse(br, request)
	if err != nil {
		log.Fatalln(err)
	}

	if response.StatusCode != 200 {
		log.Fatalln("HTTP Status: " + response.Status)
	}
	log.Println("HTTP Status: " + response.Status)

	done := make(chan bool)
	go ioCopy(conn, os.Stdin, done)
	go ioCopy(os.Stdout, conn, done)
	<-done
}

func ioCopy(dst io.Writer, src io.Reader, done chan bool) {
	buf := make([]byte, 8*1024)
	_, err := io.CopyBuffer(dst, src, buf)
	if err != nil {
		log.Println(err)
		done <- true
	}
}

func parseArgs() (proxyURI string, dstURI string, authEncrypted string, err error) {
	if len(os.Args) != 5 && len(os.Args) != 6 {
		err = errors.New("param error")
		return
	}

	proxyHost := os.Args[1]
	proxyPort := os.Args[2]
	dstHost := os.Args[3]
	dstPort := strings.Replace(os.Args[4], ":", "", -1)

	proxyURI = proxyHost + ":" + proxyPort
	dstURI = "//" + dstHost + ":" + dstPort

	if len(os.Args) == 6 {
		auth, newErr := ioutil.ReadFile(os.Args[5])
		if newErr != nil {
			err = newErr
			return
		}
		authStr := strings.TrimSuffix(string(auth), "\n")
		authEncrypted = base64.StdEncoding.EncodeToString([]byte(authStr))
	}

	return
}
