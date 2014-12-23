package mongo

import (
	"fmt"
	"os"
	"strings"

	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/subliminl/tunnel"
	"arvika.subliminl.com/developers/subliminl/util"
)

// Copy the production database to local
func CopyDbToLocal(log *log.Logger, dbName, dbPassword string) error {
	scriptFormat := `use api2
	db.dropDatabase()
	db.copyDatabase('%s', '%s', 'localhost:27018', 'Subliminl', '%s')
	`
	dbName2 := dbName + "2"
	script := fmt.Sprintf(scriptFormat, dbName, dbName2, dbPassword)

	if err := tunnel.OpenMongoTunnel(log); err != nil {
		return err
	}
	defer tunnel.CloseMongoTunnel(log)

	cmd := util.PrepareCommand(log, "mongo", "--verbose")
	cmd.SetStdin(strings.NewReader(script))
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	cmd.Announce()
	if err := cmd.RunPrintError(); err != nil {
		return err
	}
	return nil
}

// Restore the local database from the copy from remote
func RestoreDbToLocal(log *log.Logger, dbName string) error {
	scriptFormat := `use api
	db.dropDatabase()
	db.copyDatabase('%s', '%s')
	`
	dbName2 := dbName + "2"
	script := fmt.Sprintf(scriptFormat, dbName2, dbName)
	cmd := util.PrepareCommand(log, "mongo", "--verbose")
	cmd.SetStdin(strings.NewReader(script))
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	cmd.Announce()
	if err := cmd.RunPrintError(); err != nil {
		return err
	}
	return nil
}
