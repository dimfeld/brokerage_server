package types

import (
	"fmt"
	"time"
)

type Account struct {
	Id          string
	Broker      string
	Description string
}

func (a Account) String() string {
	return fmt.Sprintf("%s:%s - %s", a.Broker, a.Id, a.Description)
}

type ConnectionStatus struct {
	Connected bool
	Error     error
}

type Quote struct {
	High  int
	Low   int
	Open  int
	Close int
	Time  time.Time
}
