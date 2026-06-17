package services

import "github.com/example/go-svc-boilerplate/internal/models/value"

// PricingService is the interface core code depends on. Defining the interface
// here (consumer side) and providing a concrete impl via fx lets you swap or
// mock the client without touching callers.
type PricingService interface {
	Quote(units int) (*value.CalculatedPrice, error)
}

// SrvHolder aggregates every external service client into one injectable
// struct, mirroring StoHolder. Add a field per new service.
type SrvHolder struct {
	PricingService PricingService
}

func NewSrvHolder(pricing PricingService) *SrvHolder {
	return &SrvHolder{
		PricingService: pricing,
	}
}
