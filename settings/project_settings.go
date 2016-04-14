// Copyright (c) 2016 Pulcy.
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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/errgo"
)

var (
	maskAny = errgo.MaskFunc(errgo.Any)
)

type ProjectSettings struct {
	Image     string `json:"image"`      // Docker image name
	Registry  string `json:"registry"`   // Docker registry prefix
	Namespace string `json:"namespace"`  // Docker namespace prefix
	NoGrunt   bool   `json:"no-grunt"`   // If set, grunt won't be called even if there is a Gruntfile.js
	TagLatest bool   `json:"tag-latest"` // If set, a latest tag will be set of the docker image
	Targets   struct {
		CleanTarget string `json:"clean"`
	} `json:"targets"`
	ManifestFiles    []string `json:"manifest-files"`     // Additional manifest files
	GoVendorDir      string   `json:"go-vendor-dir"`      // If set, use this instead of `./vendor` as vendor directory.
	GradleConfigFile string   `json:"gradle-config-file"` // If set, creates a file with this path containing the current version
}

const (
	projectSettingsFile = ".pulcy"
)

// Read tries to read .pulcy settings file.
// If found the unmarshaled settings are returned, if not found nil is returned.
func Read(dir string) (*ProjectSettings, error) {
	if data, err := ioutil.ReadFile(filepath.Join(dir, projectSettingsFile)); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, maskAny(err)
		}
	} else {
		result := &ProjectSettings{}
		if err := json.Unmarshal(data, result); err != nil {
			return nil, maskAny(err)
		}
		return result, nil
	}
}
