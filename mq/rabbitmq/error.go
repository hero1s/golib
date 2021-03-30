package rabbitmq

import (
	"errors"
)

var (
	ErrBadProducer     = errors.New("MQ: Bad producer")
	ErrNoExchangeBinds = errors.New("MQ: No exchangeBinds found. You should SetExchangeBinds before open.")
	ErrStateOpened     = errors.New("MQ: Producer had been opened")
)
