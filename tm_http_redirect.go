package tm_http_redirect

import (
	"context"
	"net/http"
)

type TmHttpRedirect struct {
	next   http.Handler
	name   string
	config *Config
	rules  *[]*Rule
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	log := NamedLogger(name)
	if config == nil {
		return nil, log.CollectedError(ErrNoConfigurationFound, true)
	}
	if config.Redirects == nil {
		return nil, log.CollectedError(ErrNoConfigurationFound, true)
	}
	rules, err := ParseRules(config.Redirects, log)
	if err != nil {
		return nil, log.CollectedError(err, true)
	}
	return &TmHttpRedirect{
		next:   next,
		name:   name,
		config: config,
		rules:  &rules,
	}, nil
}
