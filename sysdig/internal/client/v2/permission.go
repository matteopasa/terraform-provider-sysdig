package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	MonitorPermissions = "%s/api/permissions/monitor/dependencies?requestedPermissions=%s"
	SecurePermissions  = "%s/api/permissions/secure/dependencies?requestedPermissions=%s"
)

type PermissionInterface interface {
	Base

	GetMonitorPermissionsWithDependencies(ctx context.Context, permissions []string) ([]Permission, error)

	GetSecurePermissionsWithDependencies(ctx context.Context, permissions []string) ([]Permission, error)
}

func (client *Client) GetMonitorPermissionsWithDependencies(ctx context.Context, permissions []string) ([]Permission, error) {
	return client.getPermissionsWithDependencies(ctx, client.GetMonitorPermissionsFor(permissions))

}

func (client *Client) GetSecurePermissionsWithDependencies(ctx context.Context, permissions []string) ([]Permission, error) {
	return client.getPermissionsWithDependencies(ctx, client.GetSecurePermissionsFor(permissions))
}

func (client *Client) getPermissionsWithDependencies(ctx context.Context, url string) ([]Permission, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []Permission{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []Permission{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[permissionListWrapper](response.Body)

	if err != nil {
		return []Permission{}, err
	}

	return wrapper.Permissions, nil
}

func (client *Client) GetMonitorPermissionsFor(permissions []string) string {
	return fmt.Sprintf(MonitorPermissions, client.config.url, strings.Join(permissions, ","))
}

func (client *Client) GetSecurePermissionsFor(permissions []string) string {
	return fmt.Sprintf(SecurePermissions, client.config.url, strings.Join(permissions, ","))
}
