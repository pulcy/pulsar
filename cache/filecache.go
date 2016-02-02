package cache

import (
	"crypto/sha512"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/juju/errgo"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	cacheDir = "~/cache/pulcy"
)

var (
	maskAny    = errgo.MaskFunc(errgo.Any)
	cacheMutex sync.Mutex
)

func Clear(key string) error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	dir, err := dir(key)
	if err != nil {
		return maskAny(err)
	}
	if err := os.RemoveAll(dir); err != nil {
		return maskAny(err)
	}
	return nil
}

func ClearAll() error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	dir, err := rootDir()
	if err != nil {
		return maskAny(err)
	}
	if err := os.RemoveAll(dir); err != nil {
		return maskAny(err)
	}
	return nil
}

// Dir returns the cache directory for a given key.
// Returns: path, isValid, error
func Dir(key string, cacheValidHours int) (string, bool, error) {
	cachedir, err := dir(key)
	if err != nil {
		return "", false, maskAny(err)
	}

	// Lock
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if cache directory exists
	s, err := os.Stat(cachedir)
	isValid := false
	if err == nil {
		// Package cache directory exists, check age.
		if cacheValidHours > 0 && s.ModTime().Add(time.Hour*time.Duration(cacheValidHours)).Before(time.Now()) {
			// Cache has become invalid
			if err := os.RemoveAll(cachedir); err != nil {
				return "", false, maskAny(err)
			}
		} else {
			// Cache is still valid
			isValid = true
		}
	} else {
		// cache directory not found, create needed
		isValid = false
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cachedir, 0777); err != nil {
		return "", false, maskAny(err)
	}

	return cachedir, isValid, nil
}

// dir returns the cache directory for a given key.
// Returns: path, error
func dir(key string) (string, error) {
	cachedirRoot, err := rootDir()
	if err != nil {
		return "", maskAny(err)
	}

	// Create hash of key
	hashBytes := sha512.Sum512([]byte(key))
	hash := fmt.Sprintf("%x", hashBytes)
	cachedir := filepath.Join(cachedirRoot, hash)

	return cachedir, nil
}

func rootDir() (string, error) {
	cachedirRoot, err := homedir.Expand(cacheDir)
	if err != nil {
		return "", maskAny(err)
	}

	return cachedirRoot, nil
}
