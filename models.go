package kanboard

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type FlexibleInt int

func (f *FlexibleInt) UnmarshalJSON(b []byte) error {
	// 数値として来た場合: 2
	var n int
	if err := json.Unmarshal(b, &n); err == nil {
		*f = FlexibleInt(n)
		return nil
	}

	// 文字列として来た場合: "1"
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*f = FlexibleInt(n)
		return nil
	}

	return fmt.Errorf("invalid int value: %s", string(b))
}

// ----- リクエスト -----

// Request Kanboard API リクエスト
type Request[T any] struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      string `json:"id"`
	Params  *T     `json:"params,omitempty"`
}

// TaskTagsParams... GetTaskTagsリクエストのパラメータ
type TaskTagsParams struct {
	TaskID int `json:"task_id"`
}

// UserParams GetUseリクエストのパラメータ
type UserParams struct {
	UserID int `json:"user_id"`
}

// ColumnsParams getColumnsリクエストのパラメータ
type ColumnsParams struct {
	ProjectID int `json:"project_id"`
}

// ----- レスポンス -----

// Response Kanboard API レスポンス
type Response struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result"`
}

// TaskTagsResult getTaskTags APIの結果
type TaskTagsResult map[string]string

type UrlResult struct {
	Board    string `json:"board"`
	Calendar string `json:"calendar"`
	List     string `json:"list"`
}

type ProjectResult struct {
	ID                  FlexibleInt `json:"id"`
	Name                string      `json:"name"`
	IsActive            FlexibleInt `json:"is_active"`
	Token               string      `json:"token"`
	LastModified        FlexibleInt `json:"last_modified"`
	IsPublic            FlexibleInt `json:"is_public"`
	IsPrivate           FlexibleInt `json:"is_private"`
	DefaultSwimlane     string      `json:"default_swimlane"`
	ShowDefaultSwimlane FlexibleInt `json:"show_default_swimlane"`
	Description         string      `json:"description"`
	Identifier          string      `json:"identifier"`
	Url                 UrlResult   `json:"url"`
}

type AllTagsResponse struct {
	ID        FlexibleInt `json:"id"`
	Name      string      `json:"name"`
	ProjectID FlexibleInt `json:"project_id"`
}

type UserResult struct {
	ID                   FlexibleInt `json:"id"`
	Username             string      `json:"username"`
	Password             string      `json:"password"`
	IsAdmin              FlexibleInt `json:"is_admin"`
	IsLdapUser           FlexibleInt `json:"is_ldap_user"`
	Name                 string      `json:"name"`
	Email                string      `json:"email"`
	GoogleID             FlexibleInt `json:"google_id"`
	GithubID             FlexibleInt `json:"github_id"`
	NotificationsEnabled FlexibleInt `json:"notifications_enabled"`
	Timezone             string      `json:"timezone"`
	Language             string      `json:"language"`
	DisableLoginForm     string      `json:"disable_login_form"`
	TwofactorActivated   FlexibleInt `json:"twofactor_activated"`
	TwofactorSecret      string      `json:"twofactor_secret"`
	Token                string      `json:"token"`
	NotificationsFilter  string      `json:"notifications_filter"`
	NbFailedLogin        string      `json:"nb_failed_login"`
	LockExpirationDate   string      `json:"lock_expiration_date"`
	IsProjectAdmin       FlexibleInt `json:"is_project_admin"`
	GitlabID             FlexibleInt `json:"gitlab_id"`
	Role                 string      `json:"role"`
	IsActive             FlexibleInt `json:"is_active"`
	AvatarPath           string      `json:"avatar_path"`
	APIAccessToken       string      `json:"api_access_token"`
	Filter               string      `json:"filter"`
	Theme                string      `json:"theme"`
}

type ColumnsResult struct {
	ID              FlexibleInt `json:"id"`
	Title           string      `json:"title"`
	Position        string      `json:"position"`
	ProjectID       FlexibleInt `json:"project_id"`
	TaskLimit       FlexibleInt `json:"task_limit"`
	Description     string      `json:"description"`
	HideInDashboard string      `json:"hide_in_dashboard"`
}
