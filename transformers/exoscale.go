package transformers

import (
	"context"
	"net/http"

	"github.com/appclacks/beckart/config"
	"github.com/appclacks/beckart/store"
	"github.com/exoscale/egoscale/v2/api"
)

type Exoscale struct {
	provider *api.SecurityProviderExoscale
}

func NewExoscale(config config.TransformerExoscaleConfig) (*Exoscale, error) {
	provider, err := api.NewSecurityProvider(config.APIKey, config.APISecret)
	if err != nil {
		return nil, err
	}
	return &Exoscale{
		provider: provider,
	}, nil
}

func (t *Exoscale) Transform(req *http.Request, store *store.Store) error {
	return t.provider.Intercept(context.Background(), req)
}
