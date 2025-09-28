package http

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	MaxBufferSz = 10 * 1024 * 1024
)

var defaultClient *http.Client
var once sync.Once

type ReqParams struct {
	Url     string
	Method  string
	Headers map[string]string
	Queries map[string]string
	Body    []byte
}

func getDefaultClient() *http.Client {
	once.Do(func() {
		if defaultClient == nil {
			// 创建自定义的Transport
			tr := &http.Transport{
				MaxIdleConns:          32,
				MaxIdleConnsPerHost:   20,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig: &tls.Config{
					//nolint:errcheck,gosec // need insecure TLS option for testing and development
					InsecureSkipVerify: true, /* #nosec G402 */
				},
			}

			// 使用自定义的Transport创建HTTP客户端
			defaultClient = &http.Client{
				Transport: tr,
				Timeout:   time.Second * 300, // 设置请求超时时间
			}
		}
	})

	return defaultClient
}

func (r *ReqParams) Call() ([]mcp.Content, error) {
	//创建请求
	client := getDefaultClient()
	req, err := http.NewRequest(r.Method, r.Url, bytes.NewBuffer(r.Body))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		output := make([]mcp.Content, 0)
		scanner := bufio.NewScanner(resp.Body)
		// 设置缓冲区大小（如果需要处理大行）
		buf := make([]byte, 0, 8*1024)
		scanner.Buffer(buf[:0], MaxBufferSz) // 最大10MB的行

		for scanner.Scan() {
			line := scanner.Text()

			if line == "" || strings.HasPrefix(line, "event:") {
				continue
			} else {
				line = strings.TrimPrefix(line, "data:")
				line = strings.TrimLeft(line, " ")
			}

			output = append(output, &mcp.TextContent{Text: line})
		}

		if err = scanner.Err(); err != nil {
			return nil, err
		}

		return output, nil
	} else {
		respBody, errR := io.ReadAll(resp.Body)
		if errR != nil {
			return nil, errR
		}

		return []mcp.Content{
			&mcp.TextContent{Text: string(respBody)},
		}, nil
	}
}
