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
}

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
	return &ProjectInfo{
		Name:    project,
		Version: oldVersion,
		pkg:     pkg,
	}, nil
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
