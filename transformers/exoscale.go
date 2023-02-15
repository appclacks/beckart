package transformers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/appclacks/beckart/config"
	"github.com/appclacks/beckart/store"
	"github.com/appclacks/beckart/template"
	"github.com/exoscale/egoscale/v2/api"
	"github.com/google/uuid"
)

type Exoscale struct {
	provider *api.SecurityProviderExoscale
}

func NewExoscale(config config.TransformerExoscaleConfig, store *store.Store) (*Exoscale, error) {
	apiKey, err := template.GenTemplate(store, fmt.Sprintf("exokey-%s", uuid.NewString()), config.APIKey)
	if err != nil {
		return nil, errors.New("fail to template exoscale transformer api key")
	}
	apiSecret, err := template.GenTemplate(store, fmt.Sprintf("exosecret-%s", uuid.NewString()), config.APISecret)
	if err != nil {
		return nil, errors.New("fail to template exoscale transformer api secret")
	}
	provider, err := api.NewSecurityProvider(apiKey.String(), apiSecret.String())
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
