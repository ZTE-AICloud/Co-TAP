package graph

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"uapregistry/logger"
)

func parsePaginationParams(pageStr, limitStr string) (page, limit int, err error) {
	// 无需分页
	if pageStr == "" || limitStr == "" {
		page = -1
		limit = -1
		return
	}

	page, err = strconv.Atoi(pageStr)
	if err != nil {
		logger.GetLogger().Errorf("pagination param is invalid, page:%s, limit:%s, %s", pageStr, limitStr, err.Error())
		return
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		logger.GetLogger().Errorf("pagination param is invalid, page:%s, limit:%s, %s", pageStr, limitStr, err.Error())
		return
	}

	if page < 1 || limit <= 0 {
		err = fmt.Errorf("pagination param is invalid, page should be started from 1,page: %s, limit:%s", pageStr, limitStr)
		logger.GetLogger().Error(err.Error())
	}

	return
}

func ResponseCodeBody(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")

	if body != "" || body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			logger.GetLogger().Errorf("failed to marshal data when build response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			if _, err := w.Write([]byte(err.Error())); err != nil {
				logger.GetLogger().Errorf("failed to write response body:%v", err)
			}
			return
		}

		w.WriteHeader(code)
		if _, err := w.Write(bs); err != nil {
			logger.GetLogger().Errorf("failed to write response body， code:%d, body: %s, :%v", code, string(bs), err)
		}
	}
}
