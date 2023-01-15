package runner

import (
	"errors"
	"fmt"

	"github.com/appclacks/beckart/config"
	"github.com/appclacks/beckart/store"
	"github.com/appclacks/beckart/transformers"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger, store *store.Store, transformers map[string]transformers.Transformer, config config.Test) error {
	logger.Info("starting test scenario")
	for _, action := range config.Actions {
		logger.Info("start executing action", zap.String("action", action.Name))
		var err error
		if action.HTTP.Target != "" {
			err = ExecuteHTTP(logger, store, transformers, action)
		} else if action.Log.Message != "" {
			err = ExecuteLog(logger, store, action)
		} else {
			logger.Error("invalid configuration for action", zap.String("action", action.Name))
			return errors.New("invalid configuration for action")
		}
		if err != nil {
			logger.Error(fmt.Sprintf("action failed: %s", err.Error()), zap.String("action", action.Name))
			return err
		}
	}
	logger.Info("test scenario finished successfully")
	return nil
}
