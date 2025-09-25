/* Started by AICoder, pid:68d39f49e7064bd08f0f101a2c0be37d */
package http

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"uapregistry/config"
	"uapregistry/leadermanager"
	"uapregistry/logger"
	"uapregistry/utils"
)

func (s *HTTPServer) addMiddlewares() {
	s.router.Use(loggingMiddleware)
	s.router.Use(httpServerLeaderMiddleware)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			logger.GetLogger().Infof("request url:%s ,method: %s, user-agent: %s", r.URL.Path, r.Method, r.Header.Get("User-Agent"))
		}
		next.ServeHTTP(w, r)
	})
}

/* Ended by AICoder, pid:68d39f49e7064bd08f0f101a2c0be37d */

/* Started by AICoder, pid:lbca32b5bbra17914596088af02aed66b8b9a599 */
func httpServerLeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" || (strings.HasPrefix(r.URL.Path, "/services") && strings.HasSuffix(r.URL.Path, "healthcheck")) {
			next.ServeHTTP(w, r)
			return
		}

		lm, err := leadermanager.GetLeaderManager()
		if err != nil {
			logger.GetLogger().Warnf("failed to GetLeaderManager:%v", err)
			next.ServeHTTP(w, r)
			return
		}

		leadAddr := lm.GetLeader()
		if leadAddr == "" {
			logger.GetLogger().Warnf("uapregistry leader is not found,r.RequestURI:%s", err, r.RequestURI)
			next.ServeHTTP(w, r)
			return
		}

		if leadAddr == utils.GetNodeIP() {
			next.ServeHTTP(w, r)
			return
		}

		// 复制请求体
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.GetLogger().Warnf("failed to read request body: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		defer r.Body.Close()

		// 重新设置请求体
		r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

		backendURL := "http://" + leadAddr + ":" + config.GetHTTPListenPort()
		req, err := http.NewRequest(r.Method, backendURL+r.RequestURI, bytes.NewReader(bodyBytes))
		if err != nil {
			logger.GetLogger().Warnf("failed to NewRequest: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		for key, value := range r.Header {
			req.Header.Set(key, value[0])
		}

		tr := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := net.Dialer{
					Timeout: 2 * time.Second, // 设置建链超时时间为2秒
				}
				return dialer.DialContext(ctx, network, addr)
			},
			IdleConnTimeout: 60 * time.Second,
		}

		client := &http.Client{
			Transport: tr,
		}
		resp, err := client.Do(req)
		if err != nil {
			logger.GetLogger().Warnf("failed to send request to leader uapregistry: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		defer resp.Body.Close()

		for key, value := range resp.Header {
			w.Header().Set(key, value[0])
		}
		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			logger.GetLogger().Warnf("io.Copy failed, err: %v", err)
		}
	})
}

/* Ended by AICoder, pid:lbca32b5bbra17914596088af02aed66b8b9a599 */
