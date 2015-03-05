package release

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type ProjectInfo struct {
	Name    string
	Version string
	pkg     map[string]interface{}
	Image   string
	NoGrunt bool // If set, grunt won't be called even if there is a Gruntfile.js
	Targets struct {
		CleanTarget string
	}
}

type ProjectSettings struct {
	Image   string `json:"image"`    // Docker image name
	NoGrunt bool   `json:"no-grunt"` // If set, grunt won't be called even if there is a Gruntfile.js
	Targets struct {
		CleanTarget string `json:"clean"`
	} `json:"targets"`
}

const (
	projectSettingsFile = ".devtool"
)

func GetProjectInfo() (*ProjectInfo, error) {
	// Read the current version and name
	project := ""
	pkg, err := readPackageJson()
	if err != nil {
		return nil, err
	}
	var oldVersion string
	if pkg != nil {
		oldVersion = pkg[versionKey].(string)
		project = pkg[nameKey].(string)
	}
	if oldVersion == "" {
		// Read version from VERSION file
		oldVersion, err = readVersion()
		if err != nil {
			return nil, err
		}
	}
	if oldVersion == "" {
		oldVersion = "0.0.1"
	}
	if project == "" {
		// Take current directory as name
		if dir, err := os.Getwd(); err != nil {
			return nil, err
		} else {
			project = path.Base(dir)
		}
	}

	// Read project settings (if any)
	image := project
	noGrunt := false
	settings, err := readProjectSettings()
	if err != nil {
		return nil, err
	}
	if settings != nil {
		if settings.Image != "" {
			image = settings.Image
		}
		noGrunt = settings.NoGrunt
	}

	result := &ProjectInfo{
		Name:    project,
		Image:   image,
		NoGrunt: noGrunt,
		Version: oldVersion,
		pkg:     pkg,
	}
	result.Targets.CleanTarget = "clean"
	if settings != nil && settings.Targets.CleanTarget != "" {
		result.Targets.CleanTarget = settings.Targets.CleanTarget
	}

	return result, nil
}

// Try to read package.json
func readPackageJson() (packageJson, error) {
	if data, err := ioutil.ReadFile(packageJsonFile); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		result := make(packageJson)
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
}

// Try to read VERSION
func readVersion() (string, error) {
	if data, err := ioutil.ReadFile(versionFile); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", err
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
			return nil, err
		}
	} else {
		result := &ProjectSettings{}
		if err := json.Unmarshal(data, result); err != nil {
			return nil, err
		}
		return result, nil
	}
}
