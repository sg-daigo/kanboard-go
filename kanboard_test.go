package kanboard

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func getLogger() *slog.Logger {
	level := new(slog.LevelVar)
	level.Set(slog.LevelDebug)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)

}

var logger = getLogger()

// テスト用サーバーの起動
func NewTestServer(t *testing.T, mockResponse string) ServerConfig {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, mockResponse)
	}))

	t.Cleanup(func() {
		ts.Close()
	})

	// サーバーのURLをモックサーバーに向ける
	return ServerConfig{
		Server: ts.URL,
		Token:  "dummy",
	}
}

func TestTags_Mock(t *testing.T) {
	// 期待するレスポンスを定義
	mockResponse := `{
		"jsonrpc": "2.0",
		"id": "test-id",
		"result": {"1": "TagA", "2": "TagB"}
	}`

	// テスト用サーバーの起動
	conf := NewTestServer(t, mockResponse)

	// 実行
	request := NewTaskTagsRequest(1)
	request.ID = "test-id" // IDを固定して比較しやすくする

	res, err := SendRequest[TaskTagsParams, TaskTagsResult](request, conf, logger)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	// 検証
	if res["1"] != "TagA" {
		t.Errorf("Expected TagA, got %s", res["1"])
	}
}

func TestProjects(t *testing.T) {
	// 期待するレスポンスを定義
	mockResponse := `{
		"jsonrpc": "2.0",
		"id": "test-id",
		"result": [
			{"id":"1","name":"demo","is_active":"1","token":"","last_modified":"1777846898","is_public":"0","is_private":"0","is_everybody_allowed":"0","default_swimlane":"Default swimlane","show_default_swimlane":"1","description":null,"identifier":"DEMO","start_date":"","end_date":"","owner_id":"1","priority_default":"0","priority_start":"0","priority_end":"3","email":null,"predefined_email_subjects":null,"per_swimlane_task_limits":"0","task_limit":"0","enable_global_tags":"1","url":{"board":"http://localhost/board/1","list":"http://localhost/list/1"}}
		]
	}`

	// テスト用サーバーの起動
	conf := NewTestServer(t, mockResponse)

	// 実行
	request := NewRequest[struct{}]("getAllUsers", nil)
	request.ID = "test-id"

	res, err := SendRequest[struct{}, []UserResult](request, conf, logger)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	// 検証
	if res[0].ID != 1 {
		t.Errorf("Expected ID, got %d", res[0].ID)
	}
	if res[0].Name != "demo" {
		t.Errorf("Expected Name, got %s", res[0].Name)
	}
}

func TestUsers(t *testing.T) {
	// 期待するレスポンスを定義
	mockResponse := `{
		"jsonrpc": "2.0",
		"id": "test-id",
		"result": [
			{"id":"1","username":"admin","password":"secret","is_admin":"1","is_ldap_user":"0","name":null,"email":null,"google_id":null,"github_id":null,"notifications_enabled":"0","timezone":null,"language":null,"disable_login_form":"0","twofactor_activated":"0","twofactor_secret":null,"token":"","notifications_filter":"4","nb_failed_login":"0","lock_expiration_date":"0","is_project_admin":"0","gitlab_id":null,"role":"app-admin","is_active":"1","avatar_path":null,"api_access_token":null,"filter":null,"theme":"light"},
			{"id":"2","username":"user1","password":"secret","is_admin":"0","is_ldap_user":"0","name":"user1","email":"user1@example.com","google_id":null,"github_id":null,"notifications_enabled":"0","timezone":"","language":"","disable_login_form":"0","twofactor_activated":"0","twofactor_secret":null,"token":"","notifications_filter":"4","nb_failed_login":"0","lock_expiration_date":"0","is_project_admin":"0","gitlab_id":null,"role":"app-user","is_active":"1","avatar_path":null,"api_access_token":null,"filter":"","theme":"light"}
		]
	}`

	// テスト用サーバーの起動
	conf := NewTestServer(t, mockResponse)

	// 実行
	request := NewRequest[struct{}]("getAllUsers", nil)
	request.ID = "test-id"

	res, err := SendRequest[struct{}, []UserResult](request, conf, logger)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	// 検証
	if len(res) != 2 {
		t.Errorf("Expected length, got %d", len(res))
	}
	if res[0].Username != "admin" {
		t.Errorf("Expected Name, got %s", res[0].Username)
	}
}

func TestUser(t *testing.T) {
	// 期待するレスポンスを定義
	mockResponse := `{
		"jsonrpc": "2.0",
		"id": "test-id",
		"result": 
			{"id":"1","username":"admin","password":"secret","is_admin":"1","is_ldap_user":"0","name":null,"email":null,"google_id":null,"github_id":null,"notifications_enabled":"0","timezone":null,"language":null,"disable_login_form":"0","twofactor_activated":"0","twofactor_secret":null,"token":"","notifications_filter":"4","nb_failed_login":"0","lock_expiration_date":"0","is_project_admin":"0","gitlab_id":null,"role":"app-admin","is_active":"1","avatar_path":null,"api_access_token":null,"filter":null,"theme":"light"}
	}`

	// テスト用サーバーの起動
	conf := NewTestServer(t, mockResponse)

	// 実行
	request := NewUserRequest(1)
	request.ID = "test-id"

	res, err := SendRequest[UserParams, UserResult](request, conf, logger)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	// 検証
	if res.Username != "admin" {
		t.Errorf("Expected Name, got %s", res.Username)
	}
}

func TestColumns(t *testing.T) {
	// 期待するレスポンスを定義
	mockResponse := `{
		"jsonrpc": "2.0",
		"id": "test-id",
		"result": [
			{"id":"1","title":"Backlog","position":"1","project_id":"1","task_limit":"0","description":"","hide_in_dashboard":"0"},
			{"id":"2","title":"Ready","position":"2","project_id":"1","task_limit":"0","description":"","hide_in_dashboard":"0"},
			{"id":"3","title":"Work in progress","position":"3","project_id":"1","task_limit":"0","description":"","hide_in_dashboard":"0"},
			{"id":"4","title":"Done","position":"4","project_id":"1","task_limit":"0","description":"","hide_in_dashboard":"0"}
		]
	}`

	// テスト用サーバーの起動
	conf := NewTestServer(t, mockResponse)

	// 実行
	request := NewColumnsRequest(1)
	request.ID = "test-id"

	res, err := SendRequest[ColumnsParams, []ColumnsResult](request, conf, logger)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	// 検証
	if len(res) != 4 {
		t.Errorf("Expected length 4, got %d", len(res))
	}
	if res[0].Title != "Backlog" {
		t.Errorf("Expected Name, got %s", res[0].Title)
	}
}
