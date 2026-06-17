package services

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/models/value"
	"github.com/example/go-svc-boilerplate/pkg/utils"
)

// pricingService is an example external service client. It holds the shared
// HTTP client + config; a real implementation would call a remote API here.
// This stub computes a price locally to keep the boilerplate runnable.
type pricingService struct {
	client *http.Client
	cfg    *config.Config
	log    *zap.Logger
}

func NewPricingService(client *http.Client, cfg *config.Config, log *zap.Logger) *pricingService {
	return &pricingService{client: client, cfg: cfg, log: log}
}

func (s *pricingService) Quote(units int) (*value.CalculatedPrice, error) {
	base := s.cfg.Pricing.BaseFare
	unitFare := utils.RoundTo2(float64(units) * s.cfg.Pricing.PerUnit)
	return &value.CalculatedPrice{
		BaseFare: base,
		UnitFare: unitFare,
		Total:    utils.RoundTo2(base + unitFare),
	}, nil
}
