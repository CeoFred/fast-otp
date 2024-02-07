package weavychat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	httpclient "github.com/CeoFred/weavychat/lib"
)

type AccessTokenResponse struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}
type WeavyServer struct {
	ServerURL string
	APIKey    string
	client    HttpClient
}

type Metadata map[string]string

type ErrorResponse struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status uint   `json:"status"`
}

type AppRequest struct {
	ID          int      `json:"id"`
	Type        AppType  `json:"type"`
	UID         string   `json:"uid"`
	DisplayName string   `json:"display_name"`
	Metadata    Metadata `json:"metadata"`
	Tags        []string `json:"tags"`
}

type WeavyApp struct {
	AppRequest
	Name         string `json:"name"`
	Description  string `json:"description"`
	ArchiveURL   string `json:"archive_url"`
	AvatarURL    string `json:"avatar_url"`
	CreatedAt    string `json:"created_at"`
	CreatedByID  int    `json:"created_by_id"`
	ModifiedAt   string `json:"modified_at"`
	ModifiedByID int    `json:"modified_by_id"`
	IsStarred    bool   `json:"is_starred"`
	IsSubscribed bool   `json:"is_subscribed"`
	IsTrashed    bool   `json:"is_trashed"`
}

type User struct {
	UID         string   `json:"uid"`
	Email       string   `json:"email"`
	GivenName   string   `json:"given_name"`
	MiddleName  string   `json:"middle_name"`
	Name        string   `json:"name"`
	FamilyName  string   `json:"family_name"`
	Nickname    string   `json:"nickname"`
	PhoneNumber string   `json:"phone_number"`
	Comment     string   `json:"comment"`
	Picture     string   `json:"picture"`
	Directory   string   `json:"directory"`
	Metadata    Metadata `json:"metadata"`
	Tags        []string `json:"tags"`
	IsSuspended bool     `json:"is_suspended"`
}

type UserProfile struct {
	ID          int    `json:"id"`
	UID         string `json:"uid"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	GivenName   string `json:"given_name"`
	MiddleName  string `json:"middle_name"`
	Name        string `json:"name"`
	FamilyName  string `json:"family_name"`
	Nickname    string `json:"nickname"`
	PhoneNumber string `json:"phone_number"`
	Comment     string `json:"comment"`
	Directory   struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		MemberCount int    `json:"member_count"`
	} `json:"directory"`
	DirectoryID int `json:"directory_id"`
	Picture     struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		MediaType    string `json:"media_type"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		Size         int    `json:"size"`
		ThumbnailURL string `json:"thumbnail_url"`
	} `json:"picture"`
	PictureID   int               `json:"picture_id"`
	AvatarURL   string            `json:"avatar_url"`
	Metadata    map[string]string `json:"metadata"`
	Tags        []string          `json:"tags"`
	Presence    string            `json:"presence"`
	CreatedAt   string            `json:"created_at"`
	ModifiedAt  string            `json:"modified_at"`
	IsSuspended bool              `json:"is_suspended"`
	IsTrashed   bool              `json:"is_trashed"`
}

func NewWeavyServer(server, apiKey string) *WeavyServer {
	return &WeavyServer{
		ServerURL: server,
		APIKey:    apiKey,
		client:    httpclient.NewAPIClient(server+"/api", apiKey),
	}
}

func (s *WeavyServer) NewApp(ctx context.Context, a *AppRequest) (*WeavyApp, error) {
	if a.UID == "" {
		return nil, errors.New("app uid is required")
	}

	if a.Type == "" {
		return nil, errors.New("app Type is required")
	}

	resp, err := s.client.Post(ctx, "/apps", a)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("API error: Type: %s, Title: %s", errorResponse.Type, errorResponse.Title)
	}

	var newAppResponse WeavyApp
	if err := json.NewDecoder(resp.Body).Decode(&newAppResponse); err != nil {
		return nil, err
	}

	return &newAppResponse, nil
}

func (s *WeavyServer) GetApp(ctx context.Context, id string) (*WeavyApp, error) {
	if id == "" {
		return nil, errors.New("app id is required")
	}

	resp, err := s.client.Get(ctx, "/apps/"+id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("API error: Type: %s, Title: %s", errorResponse.Type, errorResponse.Title)
	}

	var newAppResponse WeavyApp
	if err := json.NewDecoder(resp.Body).Decode(&newAppResponse); err != nil {
		return nil, err
	}

	return &newAppResponse, nil
}

func (s *WeavyServer) NewUser(ctx context.Context, a *User) (*UserProfile, error) {
	if a.UID == "" {
		return nil, errors.New("user uid is required")
	}

	resp, err := s.client.Post(ctx, "/users", a)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("API error: Type: %s, Title: %s", errorResponse.Type, errorResponse.Title)
	}

	var newUserResponse UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&newUserResponse); err != nil {
		return nil, err
	}

	return &newUserResponse, nil
}

func (s *WeavyServer) AddUserToApp(ctx context.Context, app_id int, users []uint) error {
	if app_id == 0 {
		return errors.New("app_id is required")
	}

	resp, err := s.client.Post(ctx, fmt.Sprintf("/apps/%d/members", app_id), users)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return err
		}

		return fmt.Errorf("API error: Type: %s, Title: %s", errorResponse.Type, errorResponse.Title)
	}

	return nil
}

func (s *WeavyServer) RemoveUserFromApp(ctx context.Context, app_id int, users []uint) error {
	if app_id == 0 {
		return errors.New("app_id is required")
	}

	resp, err := s.client.Delete(ctx, fmt.Sprintf("/apps/%d/members", app_id), users)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return err
		}

		return fmt.Errorf("API error: Type: %s, Title: %s", errorResponse.Type, errorResponse.Title)
	}

	return nil
}

func (s *WeavyServer) AppInit(ctx context.Context, app *WeavyApp) error {

	if app.UID == "" {
		return errors.New("app uid is required")
	}
	if app.Type == "" {
		return errors.New("app type is required")
	}
	requestBody := struct{ App *WeavyApp }{
		App: app,
	}

	resp, err := s.client.Post(ctx, "/apps/init", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return err
		}

		return fmt.Errorf("API error: Type: %s, Title: %s", errorResponse.Type, errorResponse.Title)
	}
	return nil

}

// Issues an access token for a user. If a user with the with the specified uid does not exists, this endpoint first creates the user and then issues an access_token
func (s *WeavyServer) GetAccessToken(ctx context.Context, uid string, ttl int) (*AccessTokenResponse, error) {

	requestBody := struct {
		ExpiresIn int `json:"expires_in"`
	}{
		ExpiresIn: ttl,
	}
	resp, err := s.client.Post(ctx, fmt.Sprintf("/users/%s/tokens", uid), requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("API error: Type: %s, Title: %s", errorResponse.Type, errorResponse.Title)
	}
	var accessToken AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&accessToken); err != nil {
		return nil, err
	}

	return &accessToken, nil

}
