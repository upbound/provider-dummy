/*
Copyright 2023 Upbound Inc.

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

import "errors"

// NewNotFound returns a new not found error.
func NewNotFound() error {
	return notFound{}
}

type notFound struct{}

func (n notFound) Error() string {
	return "not found"
}

// IsNotFound returns true if the given error is a not found error.
func IsNotFound(err error) bool {
	return errors.As(err, &notFound{})
}
