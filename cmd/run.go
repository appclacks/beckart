package cmd

import (
	"fmt"
	"os"

	"github.com/appclacks/beckart/config"
	"github.com/appclacks/beckart/runner"
	"github.com/appclacks/beckart/store"
	"github.com/appclacks/beckart/transformers"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func RunCommand() *cobra.Command {
	var configFile string
	command := &cobra.Command{
		Use:   "run",
		Short: "Run the test scenario",
		Run: func(cmd *cobra.Command, args []string) {
			zapConfig := zap.NewProductionConfig()
			logger, err := zapConfig.Build()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			file, err := os.ReadFile(configFile)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			var config config.Test
			if err := yaml.Unmarshal(file, &config); err != nil {
				logger.Error("fail to unmarshal configuration file to YAML")
				os.Exit(1)
			}

			store := store.New(config.Variables)
			transformers, err := transformers.NewTransformers(config.Transformers)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			err = runner.Run(logger, store, transformers, config)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
		},
	}
	command.PersistentFlags().StringVar(&configFile, "config", "", "configuration file path")
	command.MarkPersistentFlagRequired("config") //nolint
	return command
}
