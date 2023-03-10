package eventimpl

import (
	"github.com/google/wire"
)

var ProviderEvent = wire.NewSet(
	NewConsumerServiceImpl,
	NewSendEventImpl,
)
