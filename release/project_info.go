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

package release

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type ProjectInfo struct {
	Name      string
	Version   string
	Manifests []Manifest
	Image     string
	Registry  string
	NoGrunt   bool // If set, grunt won't be called even if there is a Gruntfile.js
	TagLatest bool `json:"tag-latest"` // If set, a latest tag will be set of the docker image
	Targets   struct {
		CleanTarget string
	}
}

type ProjectSettings struct {
	Image     string `json:"image"`      // Docker image name
	Registry  string `json:"registry"`   // Docker registry prefix
	NoGrunt   bool   `json:"no-grunt"`   // If set, grunt won't be called even if there is a Gruntfile.js
	TagLatest bool   `json:"tag-latest"` // If set, a latest tag will be set of the docker image
	Targets   struct {
		CleanTarget string `json:"clean"`
	} `json:"targets"`
	ManifestFiles []string `json:"manifest-files"` // Additional manifest files
}

const (
	projectSettingsFile = ".pulcy"
)

func GetProjectInfo() (*ProjectInfo, error) {
	// Read the current version and name
	project := ""
	manifests := []Manifest{}
	mf, err := tryReadManifest(packageJsonFile)
	if err != nil {
		return nil, maskAny(err)
	}
	var oldVersion string
	if mf != nil {
		manifests = append(manifests, *mf)
		oldVersion = mf.Data[versionKey].(string)
		project = mf.Data[nameKey].(string)
	}
	if oldVersion == "" {
		// Read version from VERSION file
		oldVersion, err = readVersion()
		if err != nil {
			return nil, maskAny(err)
		}
	}
	if oldVersion == "" {
		oldVersion = "0.0.1"
	}
	if project == "" {
		// Take current directory as name
		if dir, err := os.Getwd(); err != nil {
			return nil, maskAny(err)
		} else {
			project = path.Base(dir)
		}
	}

	// Read project settings (if any)
	image := project
	registry := ""
	noGrunt := false
	tagLatest := false
	settings, err := readProjectSettings()
	if err != nil {
		return nil, maskAny(err)
	}
	if settings != nil {
		if settings.Image != "" {
			image = settings.Image
		}
		if settings.Registry != "" {
			registry = settings.Registry
		}
		noGrunt = settings.NoGrunt
		tagLatest = settings.TagLatest

		for _, path := range settings.ManifestFiles {
			mf, err := tryReadManifest(path)
			if err != nil {
				return nil, maskAny(err)
			} else if mf == nil {
				return nil, maskAny(fmt.Errorf("manifest '%s' not found", path))
			}
			manifests = append(manifests, *mf)
		}
	}

	result := &ProjectInfo{
		Name:      project,
		Image:     image,
		Registry:  registry,
		NoGrunt:   noGrunt,
		TagLatest: tagLatest,
		Version:   oldVersion,
		Manifests: manifests,
	}
	result.Targets.CleanTarget = "clean"
	if settings != nil && settings.Targets.CleanTarget != "" {
		result.Targets.CleanTarget = settings.Targets.CleanTarget
	}

	return result, nil
}

// Try to read VERSION
func readVersion() (string, error) {
	if data, err := ioutil.ReadFile(versionFile); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", maskAny(err)
		}
	} else {
		return strings.TrimSpace(string(data)), nil
	}
}

// Try to read .devtool file
func readProjectSettings() (*ProjectSettings, error) {
	if data, err := ioutil.ReadFile(projectSettingsFile); err != nil {
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
