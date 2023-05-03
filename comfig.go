package healthman

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"time"
)

const (
	// healthmanYamlKey is a key in .config file corresponding to the Healthman client configuration
	healthmanYamlKey = "healthman"
)

// HealthmanSettings contains configurable data of a health checker
type HealthmanSettings struct {
	Period time.Duration `fig:"period"`
}

// defaultHealthmanSettings is a healthman config settings used if config is not specified
var defaultHealthmanSettings = HealthmanSettings{Period: time.Second}

// HealthmanConfiger is an interface having methods needed to access healthman settings
type HealthmanConfiger interface {
	// HealthmanConfig returns a HealthmanConfig structure
	HealthmanConfig() *HealthmanSettings
}

// NewHealthmanConfiger returns an instance of HealthmanConfiger structure that gets healthman settings from .config file
func NewHealthmanConfiger(getter kv.Getter) HealthmanConfiger {
	return &healthmaner{
		getter: getter,
	}
}

// healthmaner is a struct implementing Healthmaner interface
type healthmaner struct {
	getter kv.Getter
	once   comfig.Once
}

// HealthmanConfig returns a HealthmanConfig structure from a .config file
func (h *healthmaner) HealthmanConfig() *HealthmanSettings {
	return h.once.Do(func() any {
		var (
			cfg = defaultHealthmanSettings
			raw = kv.MustGetStringMap(h.getter, healthmanYamlKey)
		)

		if err := figure.Out(&cfg).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out healthman config"))
		}

		return &cfg
	}).(*HealthmanSettings)
}
