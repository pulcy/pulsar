package docker

import (
	"fmt"

	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/subliminl/util"
)

// Push a docker image to the arvika-ssh registry
func Push(log *log.Logger, image, dockerRegistry string) error {
	registryTag := fmt.Sprintf("%s/%s", dockerRegistry, image)
	if err := util.ExecPrintError(log, "docker", "tag", image, registryTag); err != nil {
		return err
	}
	// Push
	if err := util.ExecPrintError(log, "docker", "push", registryTag); err != nil {
		return err
	}
	// Remove registry tag
	if err := util.ExecPrintError(log, "docker", "rmi", registryTag); err != nil {
		return err
	}
	return nil
}

// Pull a docker image from the arvika-ssh registry
func Pull(log *log.Logger, image, dockerRegistry string) error {
	registryTag := fmt.Sprintf("%s/%s", dockerRegistry, image)
	// Pull
	if err := util.ExecPrintError(log, "docker", "pull", registryTag); err != nil {
		return err
	}
	if err := util.ExecPrintError(log, "docker", "tag", registryTag, image); err != nil {
		return err
	}
	// Remove registry tag
	if err := util.ExecPrintError(log, "docker", "rmi", registryTag); err != nil {
		return err
	}
	return nil
}
