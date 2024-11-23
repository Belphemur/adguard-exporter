package http

import (
	"context"
	"github.com/belphemur/adguard-exporter/internal/adguard"
	"github.com/belphemur/adguard-exporter/internal/config"
	"github.com/hellofresh/health-go/v5"
	"time"
)

func NewHealthCheck(config *config.Global) (*health.Health, error) {

	healthConfigs := make([]health.Config, len(config.Configs))
	for index, config := range config.Configs {
		healthConfigs[index] = health.Config{
			Name: config.Url,
			Check: func(ctx context.Context) error {
				client := adguard.NewClient(config)
				_, err := client.GetStatus(ctx)
				return err
			},

			SkipOnErr: true,
			Timeout:   5 * time.Second,
		}
	}

	return health.New(health.WithComponent(health.Component{
		Name:    "adguard-exporter",
		Version: config.Version,
	}), health.WithChecks(healthConfigs...))
}
