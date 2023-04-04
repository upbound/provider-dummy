/*
Copyright 2023 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package robots

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
)

// Robot is a simple client-side struct that represents a robot.
type Robot struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// NewClient returns a new client that can be used to interact with the
// robots service.
func NewClient(host *url.URL) *Client {
	return &Client{
		endpoint: host.JoinPath("robots"),
	}
}

// Client is a client that can be used to interact with the robots service.
type Client struct {
	endpoint *url.URL
}

// Get returns the robot with the given name.
func (c *Client) Get(ctx context.Context, name string) (*Robot, error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet,
		fmt.Sprintf("%s?&name=%s", c.endpoint.String(), name),
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get robot")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read body")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, NewNotFound()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(data))
	}
	r := &Robot{}
	if err := json.Unmarshal(data, r); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal response")
	}
	return r, nil
}

// Create creates the given robot.
func (c *Client) Create(ctx context.Context, r *Robot) error {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(r); err != nil {
		return errors.Wrap(err, "cannot encode robot")
	}
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, c.endpoint.String(), payload,
	)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot create robot")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "cannot read body")
		}
		return fmt.Errorf("%d: %s", resp.StatusCode, string(data))
	}
	return nil
}

// Update updates the given robot.
func (c *Client) Update(ctx context.Context, r *Robot) error {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(r); err != nil {
		return errors.Wrap(err, "cannot encode robot")
	}
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPut,
		fmt.Sprintf("%s?&name=%s", c.endpoint.String(), r.Name),
		payload,
	)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot update robot")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode == http.StatusNotFound {
		return NewNotFound()
	}
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "cannot read body")
		}
		return fmt.Errorf("%d: %s", resp.StatusCode, string(data))
	}
	return nil
}

// Delete deletes the robot with the given name.
func (c *Client) Delete(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodDelete,
		fmt.Sprintf("%s?&name=%s", c.endpoint.String(), name),
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot delete robot")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode == http.StatusNotFound {
		return NewNotFound()
	}
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "cannot read body")
		}
		return fmt.Errorf("%d: %s", resp.StatusCode, string(data))
	}
	return nil
}
