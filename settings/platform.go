// Copyright (c) 2018 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package settings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Platform is an os+architecture pair that follows the GOOS/GOARCH
// naming convention.
type Platform struct {
	OS   string `json:"os"`   // Name of operating system
	Arch string `json:"arch"` // Name of architecture
}

// String converts a platform to string
func (p Platform) String() string {
	return p.OS + "/" + p.Arch
}

// MarshalJSON marshals a Platform into a JSON string
func (p Platform) MarshalJSON() ([]byte, error) {
	s := p.String()
	encoded, err := json.Marshal(s)
	if err != nil {
		return nil, maskAny(err)
	}
	return encoded, nil
}

// UnmarshalJSON unmarshals a JSON string into a Platform.
func (p *Platform) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return maskAny(err)
	}
	parts := strings.Split(s, "/")
	switch len(parts) {
	case 1:
		p.OS = "linux"
		p.Arch = parts[0]
	case 2:
		p.OS = parts[0]
		p.Arch = parts[1]
	default:
		return maskAny(fmt.Errorf("Invalid platform '%s', expected '<os>/<architecture>'", s))
	}

	return nil
}
