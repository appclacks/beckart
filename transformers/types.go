package transformers

import (
	"net/http"

	"github.com/appclacks/beckart/store"
)

type Transformer interface {
	Transform(req *http.Request, store *store.Store) error
}
