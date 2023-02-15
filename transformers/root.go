package transformers

import (
	"fmt"

	"github.com/appclacks/beckart/config"
	"github.com/appclacks/beckart/store"
)

func NewTransformers(transformersConfig map[string]config.Transformer, store *store.Store) (map[string]Transformer, error) {
	result := make(map[string]Transformer)
	for name, config := range transformersConfig {
		if config.Exoscale.APIKey != "" {
			transformer, err := NewExoscale(config.Exoscale, store)
			if err != nil {
				return nil, fmt.Errorf("fail to build Exoscale transformer %s: %w", name, err)
			}
			result[name] = transformer
		} else {
			return nil, fmt.Errorf("invalid configuration for transformer %s", name)

		}
	}
	return result, nil
}
