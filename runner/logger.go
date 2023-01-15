package runner

import (
	"github.com/appclacks/beckart/config"
	"github.com/appclacks/beckart/store"
	"github.com/appclacks/beckart/template"
	"go.uber.org/zap"
)

func ExecuteLog(logger *zap.Logger, store *store.Store, action config.Action) error {
	message, err := template.GenTemplate(store, action.Name, action.Log.Message)
	if err != nil {
		return err
	}
	logger.Info(message.String())
	return nil
}
