package shared

import (
	"time"
)

type (
	DomainEvent interface {
		CreateAt() time.Time
		Identity() string
	}
)
