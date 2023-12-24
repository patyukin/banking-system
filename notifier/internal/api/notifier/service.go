package notifier

import (
	desc "github.com/patyukin/banking-system/notifier/pkg/notifier_v1"
)

type Implementation struct {
	desc.UnimplementedNotifierV1Server
}

func NewImplementation() *Implementation {
	return &Implementation{}
}
