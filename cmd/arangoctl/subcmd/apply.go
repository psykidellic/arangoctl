package subcmd

import (
	"context"
	"fmt"
	"github.com/psykidellic/arangoctl/pkg/arangoctl"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var applyCmd = cobra.Command{
	Use:   "apply",
	Short: "apply one of more configs",
	Args:  cobra.MinimumNArgs(1),
	RunE:   applyRun,
}

type applyCmdConfig struct {
	clusterConfig string
	dryRun        bool
}

var applyConfig applyCmdConfig

func init() {
	applyCmd.Flags().StringVar(
		&applyConfig.clusterConfig,
		"cluster-config",
		os.Getenv("ARANGOCTL_CLUSTER_CONFIG"),
		"Cluster config path",
	)

	applyCmd.Flags().BoolVar(
		&applyConfig.dryRun,
		"dry-run",
		false,
		"Do a dry run of",
	)

	rootCmd.AddCommand(&applyCmd)
}

func applyRun(cmd *cobra.Command, args []string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get full path for config path
	clusterConfigPath, err := getClusterConfigPath(applyConfig.clusterConfig)
	if err != nil {
		return err
	}

	// Load up the cluster config and connect to cluster
	clusterConfig, err := arangoctl.LoadClusterFile(clusterConfigPath)
	if err != nil {
		return err
	}

	adminClient, err := clusterConfig.NewAdminClient(ctx)
	if err != nil {
		return err
	}

	matchCount := 0
	for _, arg := range args {

		// We support passing of folders as argument so that
		// multiple resource specs can be just applied in one shot
		matches, err := filepath.Glob(arg)
		if err != nil {
			log.Errorf("%+v", err)
			continue
		}

		for _, match := range matches {

			// Load the relevant resource config from each type
			// and then get the resource object and apply them
			resourceConfig, err := arangoctl.LoadResourceFile(match)
			if err != nil {
				log.Errorf("%+v", err)
				continue
			}

			applyErr := resourceConfig.GetResource().Apply(adminClient)
			if applyErr != nil {
				log.Errorf("%+v", applyErr)
				continue
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return fmt.Errorf("Error applying any of the collection config provide by (%+v)", args)
	}

	return nil
}

func getClusterConfigPath(configPath string) (string, error) {
	return filepath.Abs(configPath)
}