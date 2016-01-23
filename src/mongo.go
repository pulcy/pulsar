package main

import (
	"os"

	"github.com/spf13/cobra"

	"git.pulcy.com/pulcy/pulcy/mongo"
)

const (
	defaultDatabase  = "api"
	dbPasswordEnvKey = "MONGODB_DB_PASSWORD"
)

var (
	dbName string
	dbCmd  = &cobra.Command{
		Use:   "db",
		Short: "Perform database operations",
		Long:  "Perform database operations",
		Run:   UsageFunc,
	}
	dbBackupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Copy remote database to local",
		Long:  "Copy remote database to local (remote-api -> api2)",
		Run:   runDbBackup,
	}
	dbRestoreToLocalCmd = &cobra.Command{
		Use:   "restore-to-local",
		Short: "Restore a db copied from remote to local",
		Long:  "Restore a db copied from remote to local (api2 -> api)",
		Run:   runDbRestoreToLocal,
	}
)

func init() {
	dbCmd.Flags().StringVarP(&dbName, "database", "d", defaultDatabase, "Specify database name")
	dbCmd.AddCommand(dbRestoreToLocalCmd)
	dbCmd.AddCommand(dbBackupCmd)
	mainCmd.AddCommand(dbCmd)
}

func runDbRestoreToLocal(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		err := mongo.RestoreDbToLocal(log, dbName)
		if err != nil {
			Quitf("%s\n", err)
		}
	default:
		CommandError(cmd, "Too many arguments\n")
	}
}

func runDbBackup(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		dbPassword := os.Getenv(dbPasswordEnvKey)
		if dbPassword == "" {
			Quitf("Set %s first. [%s]\n", dbPasswordEnvKey, dbPassword)
		}
		err := mongo.CopyDbToLocal(log, dbName, dbPassword)
		if err != nil {
			Quitf("%s\n", err)
		}
	default:
		CommandError(cmd, "Too many arguments\n")
	}
}
