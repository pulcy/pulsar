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

package golang

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/op/go-logging"

	"github.com/pulcy/pulsar/util"
)

type FlattenFlags struct {
	VendorDir string
}

// Flatten copies all directories found in the given vendor directory to the GOPATH
// and flattens all vendor directories found in the GOPATH.
func Flatten(log *log.Logger, flags *FlattenFlags) error {
	vendorDir, err := filepath.Abs(flags.VendorDir)
	if err != nil {
		return maskAny(err)
	}
	goSrcDir := filepath.Join(gopath, "src")
	if err := copyFromVendor(log, goSrcDir, vendorDir); err != nil {
		return maskAny(err)
	}
	if err := flattenGoDir(log, goSrcDir, goSrcDir); err != nil {
		return maskAny(err)
	}

	return nil
}

func copyFromVendor(log *log.Logger, goDir, vendorDir string) error {
	entries, err := ioutil.ReadDir(vendorDir)
	if err != nil {
		return maskAny(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		entryVendorDir := filepath.Join(vendorDir, entry.Name())
		entryGoDir := filepath.Join(goDir, entry.Name())
		if _, err := os.Stat(entryGoDir); os.IsNotExist(err) {
			// We must create a link
			log.Infof("copying %s", entryVendorDir)
			if err := os.MkdirAll(goDir, 0777); err != nil {
				return maskAny(err)
			}
			if err := util.ExecPrintError(nil, "rsync", "-a", entryVendorDir, goDir); err != nil {
				return maskAny(err)
			}

		} else if err != nil {
			return maskAny(err)
		} else {
			// entry already exists in godir, recurse into the directory
			if err := copyFromVendor(log, entryGoDir, entryVendorDir); err != nil {
				return maskAny(err)
			}
		}
	}

	return nil
}

func flattenGoDir(log *log.Logger, goSrcDir, curDir string) error {
	vendorDir := filepath.Join(curDir, "vendor")
	if _, err := os.Stat(vendorDir); err == nil {
		if err := util.ExecPrintError(nil, "rsync", "-a", "--ignore-existing", vendorDir+"/", goSrcDir); err != nil {
			return maskAny(err)
		}
		if err := os.RemoveAll(vendorDir); err != nil {
			return maskAny(err)
		}
	}

	// Recurse into sub-directories
	entries, err := ioutil.ReadDir(curDir)
	if err != nil {
		return maskAny(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if err := flattenGoDir(log, goSrcDir, filepath.Join(curDir, entry.Name())); err != nil {
			return maskAny(err)
		}
	}

	return nil
}
