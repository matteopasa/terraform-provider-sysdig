package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var GroupMappingNotFound = errors.New("group mapping not found")

type GroupMapper interface {
	CreateGroupMapping(ctx context.Context, request *GroupMapping) (*GroupMapping, error)
	UpdateGroupMapping(ctx context.Context, request *GroupMapping, id int) (*GroupMapping, error)
	DeleteGroupMapping(ctx context.Context, id int) error
	GetGroupMapping(ctx context.Context, id int) (*GroupMapping, error)
}

func (client *client) CreateGroupMapping(ctx context.Context, request *GroupMapping) (*GroupMapping, error) {
	payload, err := request.ToJSON()
	if err != nil {
		return nil, err
	}

	response, err := client.DoSysdigRequest(ctx, http.MethodPost, client.CreateGroupMappingUrl(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var groupMapping GroupMapping
	err = json.Unmarshal(body, &groupMapping)
	if err != nil {
		return nil, err
	}

	return &groupMapping, nil
}

func (client *client) UpdateGroupMapping(ctx context.Context, request *GroupMapping, id int) (*GroupMapping, error) {
	payload, err := request.ToJSON()
	if err != nil {
		return nil, err
	}

	response, err := client.DoSysdigRequest(ctx, http.MethodPut, client.UpdateGroupMappingUrl(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var groupMapping GroupMapping
	err = json.Unmarshal(body, &groupMapping)
	if err != nil {
		return nil, err
	}

	return &groupMapping, nil
}

func (client *client) DeleteGroupMapping(ctx context.Context, id int) error {
	response, err := client.DoSysdigRequest(ctx, http.MethodDelete, client.DeleteGroupMappingUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *client) GetGroupMapping(ctx context.Context, id int) (*GroupMapping, error) {
	response, err := client.DoSysdigRequest(ctx, http.MethodGet, client.GetGroupMappingUrl(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, GroupMappingNotFound
		}
		return nil, client.ErrorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var groupMapping GroupMapping
	err = json.Unmarshal(body, &groupMapping)
	if err != nil {
		return nil, err
	}

	return &groupMapping, nil
}
func (client *client) GetGroupMappingUrl(id int) string {
	return fmt.Sprintf("%s/api/groupmappings/%d", client.URL, id)
}

func (client *client) CreateGroupMappingUrl() string {
	return fmt.Sprintf("%s/api/groupmappings", client.URL)
}

func (client *client) UpdateGroupMappingUrl(id int) string {
	return fmt.Sprintf("%s/api/groupmappings/%d", client.URL, id)
}

func (client *client) DeleteGroupMappingUrl(id int) string {
	return fmt.Sprintf("%s/api/groupmappings/%d", client.URL, id)
}