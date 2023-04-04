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

import "github.com/pkg/errors"

type notFound struct{}

func (n *notFound) Error() string {
	return "not found"
}

func (n *notFound) As(err interface{}) bool {
	_, ok := err.(*notFound)
	return ok
}

// IsNotFound returns true if the given error is a not found error.
func IsNotFound(err error) bool {
	return errors.As(err, &notFound{})
}
