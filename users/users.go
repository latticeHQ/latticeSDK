package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Me is used as a replacement for the authenticated user's own identifier.
const Me = "me"

// Service provides user management operations against the Lattice Runtime API.
type Service struct {
	client *client.Client
}

// New creates a new users Service backed by the given client.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// HasFirstUser returns whether the first user has been created on the deployment.
func (s *Service) HasFirstUser(ctx context.Context) (bool, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/first", nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, client.ReadBodyAsError(resp)
	}
	return true, nil
}

// CreateFirstUser creates the initial superadmin user on a new Lattice deployment.
// If any users already exist, the request will fail.
func (s *Service) CreateFirstUser(ctx context.Context, req CreateFirstUserRequest) (CreateFirstUserResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users/first", req)
	if err != nil {
		return CreateFirstUserResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return CreateFirstUserResponse{}, client.ReadBodyAsError(resp)
	}

	var result CreateFirstUserResponse
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// CreateUser creates a new user.
//
// Deprecated: Use CreateUserWithOrgs instead.
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	if req.DisableLogin {
		req.UserLoginType = LoginTypeNone
	}
	return s.CreateUserWithOrgs(ctx, CreateUserRequestWithOrgs{
		Email:           req.Email,
		Username:        req.Username,
		Name:            req.Name,
		Password:        req.Password,
		UserLoginType:   req.UserLoginType,
		OrganizationIDs: []uuid.UUID{req.OrganizationID},
	})
}

// CreateUserWithOrgs creates a new user with membership in the specified organizations.
func (s *Service) CreateUserWithOrgs(ctx context.Context, req CreateUserRequestWithOrgs) (User, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users", req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return User{}, client.ReadBodyAsError(resp)
	}

	var user User
	return user, json.NewDecoder(resp.Body).Decode(&user)
}

// GetUser returns a user by ID or username.
func (s *Service) GetUser(ctx context.Context, userIdent string) (User, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/users/%s", userIdent), nil)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, client.ReadBodyAsError(resp)
	}

	var user User
	return user, json.NewDecoder(resp.Body).Decode(&user)
}

// ListUsers returns users matching the given request filters.
// If no filters are set, all users are returned in a single page.
func (s *Service) ListUsers(ctx context.Context, req UsersRequest) (GetUsersResponse, error) {
	opts := []client.RequestOption{
		paginationOption(req),
		searchOption(req),
	}

	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users", nil, opts...)
	if err != nil {
		return GetUsersResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetUsersResponse{}, client.ReadBodyAsError(resp)
	}

	var result GetUsersResponse
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// DeleteUser deletes a user by ID.
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/users/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateUserProfile updates a user's username and display name.
func (s *Service) UpdateUserProfile(ctx context.Context, user string, req UpdateUserProfileRequest) (User, error) {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/users/%s/profile", user), req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, client.ReadBodyAsError(resp)
	}

	var result User
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// UpdateUserStatus sets the user's account status (active or suspended).
func (s *Service) UpdateUserStatus(ctx context.Context, user string, status UserStatus) (User, error) {
	path := fmt.Sprintf("/api/v2/users/%s/status/", user)
	switch status {
	case UserStatusActive:
		path += "activate"
	case UserStatusSuspended:
		path += "suspend"
	default:
		return User{}, fmt.Errorf("status %q is not supported", status)
	}

	resp, err := s.client.Request(ctx, http.MethodPut, path, nil)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, client.ReadBodyAsError(resp)
	}

	var result User
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// UpdateUserPassword changes a user's password.
func (s *Service) UpdateUserPassword(ctx context.Context, user string, req UpdateUserPasswordRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/users/%s/password", user), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ValidateUserPassword validates the complexity of a password.
func (s *Service) ValidateUserPassword(ctx context.Context, req ValidateUserPasswordRequest) (ValidateUserPasswordResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users/validate-password", req)
	if err != nil {
		return ValidateUserPasswordResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ValidateUserPasswordResponse{}, client.ReadBodyAsError(resp)
	}

	var result ValidateUserPasswordResponse
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// UpdateUserAppearanceSettings updates the user's theme preference.
func (s *Service) UpdateUserAppearanceSettings(ctx context.Context, user string, req UpdateUserAppearanceSettingsRequest) (User, error) {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/users/%s/appearance", user), req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, client.ReadBodyAsError(resp)
	}

	var result User
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// UpdateUserRoles replaces a user's site-wide roles with the specified set.
// Include ALL roles the user should have.
func (s *Service) UpdateUserRoles(ctx context.Context, user string, req UpdateRoles) (User, error) {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/users/%s/roles", user), req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, client.ReadBodyAsError(resp)
	}

	var result User
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// GetUserRoles returns all roles assigned to the user, both site-wide and per-organization.
func (s *Service) GetUserRoles(ctx context.Context, user string) (UserRoles, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/users/%s/roles", user), nil)
	if err != nil {
		return UserRoles{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserRoles{}, client.ReadBodyAsError(resp)
	}

	var roles UserRoles
	return roles, json.NewDecoder(resp.Body).Decode(&roles)
}

// LoginWithPassword authenticates a user with email and password and returns a session token.
// Call client.SetSessionToken() to apply the token for subsequent requests.
func (s *Service) LoginWithPassword(ctx context.Context, req LoginWithPasswordRequest) (LoginWithPasswordResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users/login", req)
	if err != nil {
		return LoginWithPasswordResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return LoginWithPasswordResponse{}, client.ReadBodyAsError(resp)
	}

	var result LoginWithPasswordResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return LoginWithPasswordResponse{}, err
	}
	return result, nil
}

// RequestOneTimePasscode sends a one-time passcode to the user's email for password reset.
func (s *Service) RequestOneTimePasscode(ctx context.Context, req RequestOneTimePasscodeRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users/otp/request", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ChangePasswordWithOneTimePasscode resets a user's password using a one-time passcode.
func (s *Service) ChangePasswordWithOneTimePasscode(ctx context.Context, req ChangePasswordWithOneTimePasscodeRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users/otp/change-password", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ConvertLoginType initiates conversion of the authenticated user's login method
// (e.g. from password to OAuth). The response contains the OAuth state needed to
// complete the conversion flow.
func (s *Service) ConvertLoginType(ctx context.Context, req ConvertLoginRequest) (OAuthConversionResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users/me/convert-login", req)
	if err != nil {
		return OAuthConversionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return OAuthConversionResponse{}, client.ReadBodyAsError(resp)
	}

	var result OAuthConversionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return OAuthConversionResponse{}, err
	}
	return result, nil
}

// Logout invalidates the current session.
// Call client.SetSessionToken("") to clear the local token after logout.
func (s *Service) Logout(ctx context.Context) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users/logout", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// GetAuthMethods returns the authentication methods enabled on the deployment.
func (s *Service) GetAuthMethods(ctx context.Context) (AuthMethods, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/authmethods", nil)
	if err != nil {
		return AuthMethods{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AuthMethods{}, client.ReadBodyAsError(resp)
	}

	var methods AuthMethods
	return methods, json.NewDecoder(resp.Body).Decode(&methods)
}

// GetUserQuietHoursSchedule returns the quiet hours schedule for a user.
// This endpoint is only available in enterprise editions.
func (s *Service) GetUserQuietHoursSchedule(ctx context.Context, user string) (UserQuietHoursScheduleResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/users/%s/quiet-hours", user), nil)
	if err != nil {
		return UserQuietHoursScheduleResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserQuietHoursScheduleResponse{}, client.ReadBodyAsError(resp)
	}

	var result UserQuietHoursScheduleResponse
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// UpdateUserQuietHoursSchedule updates the quiet hours schedule for a user.
// This endpoint is only available in enterprise editions.
func (s *Service) UpdateUserQuietHoursSchedule(ctx context.Context, user string, req UpdateUserQuietHoursScheduleRequest) (UserQuietHoursScheduleResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/users/%s/quiet-hours", user), req)
	if err != nil {
		return UserQuietHoursScheduleResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserQuietHoursScheduleResponse{}, client.ReadBodyAsError(resp)
	}

	var result UserQuietHoursScheduleResponse
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// GetUserAutofillParameters returns recently used template parameters for the user.
func (s *Service) GetUserAutofillParameters(ctx context.Context, user string, templateID uuid.UUID) ([]UserParameter, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/users/%s/autofill-parameters?template_id=%s", user, templateID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var params []UserParameter
	return params, json.NewDecoder(resp.Body).Decode(&params)
}

// ListOrganizationsByUser returns all organizations the user is a member of.
func (s *Service) ListOrganizationsByUser(ctx context.Context, user string) ([]Organization, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/users/%s/organizations", user), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var orgs []Organization
	return orgs, json.NewDecoder(resp.Body).Decode(&orgs)
}

// GetOrganizationByUserAndName returns a specific organization by user identifier and org name.
func (s *Service) GetOrganizationByUserAndName(ctx context.Context, user string, name string) (Organization, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/users/%s/organizations/%s", user, name), nil)
	if err != nil {
		return Organization{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Organization{}, client.ReadBodyAsError(resp)
	}

	var org Organization
	return org, json.NewDecoder(resp.Body).Decode(&org)
}

// paginationOption builds a RequestOption that applies pagination query parameters from a UsersRequest.
func paginationOption(req UsersRequest) client.RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		if req.AfterID != uuid.Nil {
			q.Set("after_id", req.AfterID.String())
		}
		if req.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", req.Limit))
		}
		if req.Offset > 0 {
			q.Set("offset", fmt.Sprintf("%d", req.Offset))
		}
		r.URL.RawQuery = q.Encode()
	}
}

// searchOption builds a RequestOption that applies search filters from a UsersRequest.
func searchOption(req UsersRequest) client.RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		var params []string
		if req.Search != "" {
			params = append(params, req.Search)
		}
		if req.Status != "" {
			params = append(params, "status:"+string(req.Status))
		}
		if req.Role != "" {
			params = append(params, "role:"+req.Role)
		}
		if req.SearchQuery != "" {
			params = append(params, req.SearchQuery)
		}
		if len(params) > 0 {
			q.Set("q", strings.Join(params, " "))
		}
		r.URL.RawQuery = q.Encode()
	}
}
