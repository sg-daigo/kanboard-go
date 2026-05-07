package kanboard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type ServerConfig struct {
	Server string
	Token  string
}

func NewServerConfig() ServerConfig {
	return ServerConfig{
		Server: "http://localhost",
		Token:  os.Getenv("KB_TOKEN"),
	}
}

// NewRequest APIリクエストを生成（共通）
func NewRequest[T any](method string, params *T) Request[T] {
	return Request[T]{
		JSONRPC: "2.0",
		Method:  method,
		ID:      uuid.NewString(),
		Params:  params,
	}
}

// NewTaskTagsRequest getTaskTags APIリクエストを生成
func NewTaskTagsRequest(taskID int) Request[TaskTagsParams] {
	params := TaskTagsParams{
		TaskID: taskID,
	}
	return NewRequest[TaskTagsParams]("getTaskTags", &params)
}

// NewUserRequest getUser APIリクエストを生成
func NewUserRequest(userID int) Request[UserParams] {
	params := UserParams{
		UserID: userID,
	}
	return NewRequest[UserParams]("getUser", &params)
}

func NewColumnsRequest(projectID int) Request[ColumnsParams] {
	params := ColumnsParams{
		ProjectID: projectID,
	}
	return NewRequest("getColumns", &params)
}

var httpClient = &http.Client{
	Timeout: time.Second * 10, // タイムアウト10秒
}

// SendRequestRaw Kanboardにリクエストを送信してレスポンスを受け取る
func SendRequestRaw[T any](req Request[T], conf ServerConfig, logger *slog.Logger) (res Response, err error) {
	b, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("json marshal failed: %w", err)
	}

	url := conf.Server + "/jsonrpc.php"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	r.Header.Set("Content-Type", "application/json")
	r.SetBasicAuth("jsonrpc", conf.Token)
	resp, err := httpClient.Do(r)
	if err != nil {
		return res, fmt.Errorf("request post failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("http status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, fmt.Errorf("read response body failed: %w", err)
	}
	if logger != nil {
		logger.Debug("Response " + string(body))
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return res, fmt.Errorf("json unmarshal failed: %w", err)
	}

	if res.ID != req.ID {
		return res, fmt.Errorf("response id %s not equal to %s", res.ID, req.ID)
	}

	return res, nil
}

func SendRequest[P any, R any](req Request[P], conf ServerConfig, logger *slog.Logger) (res R, err error) {
	response, err := SendRequestRaw[P](req, conf, logger)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(response.Result, &res); err != nil {
		return res, fmt.Errorf("json unmarshal failed: %w", err)
	}

	return res, nil
}
