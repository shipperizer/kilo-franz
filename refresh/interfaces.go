package refresh

import (
	"context"

	"github.com/shipperizer/kilo-franz/config"
	kiloTLS "github.com/shipperizer/kilo-franz/tls"
)

// AutoRefreshXInterface is the interface for the autorefresh functionality
type AutoRefreshXInterface interface {
	Configure(ctx context.Context, config *kiloTLS.TLSConfig)
	Object(ctx context.Context) (config.RefreshableInterface, error)
	Stats() interface{}
	Refresh(ctx context.Context) (config.RefreshableInterface, error)
	Stop()
}
