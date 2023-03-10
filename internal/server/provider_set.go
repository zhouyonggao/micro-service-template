package server

import (
	"github.com/google/wire"
)

// ProviderSetServer is server providers.
var ProviderSetServer = wire.NewSet(NewGRPCServer, NewConsumerServer, NewMonitorLogServer, NewCliServer)
