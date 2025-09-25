package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"uapregistry/config"
)

const ModuleName = "uapregistry"

var (
	httpClient *http.Client
)

func getHTTPClient() *http.Client {
	if httpClient == nil {
		tr := &http.Transport{
			/* #nosec */
			TLSClientConfig: &tls.Config{InsecureSkipVerify: config.GetInsecureSkipVerify()},
			IdleConnTimeout: 60 * time.Second,
		}
		httpClient = &http.Client{Transport: tr}
	}

	return httpClient
}

func HTTPGetWithTime(url string, timeout time.Duration) ([]byte, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", ModuleName)

	res, err := getHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	buf, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func HTTPGet(url string) ([]byte, error) {
	return HTTPGetWithTime(url, 0)
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Printf("Dial errors:%v in GetOutboundIP", err)
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}

func GetNodeIP() string {
	nodeIP := os.Getenv("BIND_IP")
	if nodeIP == "" {
		nodeIP = GetOutboundIP()
	}

	return nodeIP
}

/* Ended by AICoder, pid:18047b2529744602b85f80e470620f5c */

func IsValidIntStr(intStr string) bool {
	if intStr == "" {
		return true
	}
	_, err := strconv.Atoi(intStr)
	return err == nil
}
