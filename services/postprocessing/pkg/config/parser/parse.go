package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cs3org/reva/v2/pkg/events"
	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config/defaults"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
)

// ParseConfig loads configuration from known paths.
func ParseConfig(cfg *config.Config) error {
	_, err := ociscfg.BindSourcesToStructs(cfg.Service.Name, cfg)
	if err != nil {
		return err
	}

	defaults.EnsureDefaults(cfg)

	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	defaults.Sanitize(cfg)

	return Validate(cfg)
}

// Validate validates the config
func Validate(cfg *config.Config) error {
	if cfg.Postprocessing.Virusscan {
		if !contains(cfg.Postprocessing.Steps, events.PPStepAntivirus) {
			cfg.Postprocessing.Steps = append(cfg.Postprocessing.Steps, string(events.PPStepAntivirus))
			fmt.Printf("ATTENTION: POSTPROCESSING_VIRUSSCAN is deprecated. Use `POSTPROCESSING_STEPS=%v` in the future\n", strings.Join(cfg.Postprocessing.Steps, ","))
		}
	}

	if cfg.Postprocessing.Delayprocessing != 0 {
		if !contains(cfg.Postprocessing.Steps, events.PPStepDelay) {
			if len(cfg.Postprocessing.Steps) > 0 {
				s := strings.Join(append(cfg.Postprocessing.Steps, string(events.PPStepDelay)), ",")
				fmt.Printf("Added delay step to the list of postprocessing steps. NOTE: Use envvar `POSTPROCESSING_STEPS=%s` to suppress this message and choose the order of postprocessing steps.\n", s)
			}

			cfg.Postprocessing.Steps = append(cfg.Postprocessing.Steps, string(events.PPStepDelay))
		}
	}
	return nil
}

func contains(all []string, candidate events.Postprocessingstep) bool {
	for _, s := range all {
		if s == string(candidate) {
			return true
		}
	}
	return false
}
