package sidecarsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// AuthGoogleInstanceIdentity authenticates using a Google Cloud instance identity token.
func (s *Service) AuthGoogleInstanceIdentity(ctx context.Context, req GoogleInstanceIdentityToken) (AuthenticateResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentsidecars/google-instance-identity", req)
	if err != nil {
		return AuthenticateResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AuthenticateResponse{}, client.ReadBodyAsError(resp)
	}

	var result AuthenticateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return AuthenticateResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// AuthAWSInstanceIdentity authenticates using an AWS instance identity document and signature.
func (s *Service) AuthAWSInstanceIdentity(ctx context.Context, req AWSInstanceIdentityToken) (AuthenticateResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentsidecars/aws-instance-identity", req)
	if err != nil {
		return AuthenticateResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AuthenticateResponse{}, client.ReadBodyAsError(resp)
	}

	var result AuthenticateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return AuthenticateResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// AuthAzureInstanceIdentity authenticates using an Azure instance identity token.
func (s *Service) AuthAzureInstanceIdentity(ctx context.Context, req AzureInstanceIdentityToken) (AuthenticateResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentsidecars/azure-instance-identity", req)
	if err != nil {
		return AuthenticateResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AuthenticateResponse{}, client.ReadBodyAsError(resp)
	}

	var result AuthenticateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return AuthenticateResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}
