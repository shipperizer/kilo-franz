package refresh

import (
	"context"
	"sync"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
)

// AutoRefreshXInterface is the interface for the autorefresh functionality
type AutoRefreshXInterface interface {
	Configure(context.Context, config.TLSConfigInterface, config.SASLConfigInterface)
	Object(context.Context) (core.RefreshableInterface, error)
	Stats() interface{}
	Refresh(context.Context) (core.RefreshableInterface, error)
	Stop()
}

// AutoRefreshXConfigInterface is the refresh.AutoRefreshX config interface, embeds ConfigInterface
type AutoRefreshXConfigInterface interface {
	config.ConfigInterface
	GetMutexObj() *sync.RWMutex
}
