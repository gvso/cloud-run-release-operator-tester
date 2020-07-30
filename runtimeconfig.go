package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
)

type runtimeConfig struct {
	service    *runtimeconfig.Service
	project    string
	configName string

	// If the webservice should respect the latency and error rate configuration.
	respectVariables bool
	lastCheck        time.Time
}

func newRuntimeConfig(ctx context.Context, project, configName string) (*runtimeConfig, error) {
	runtimeconfigService, err := runtimeconfig.NewService(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client for Runtime Config")
	}

	return &runtimeConfig{
		service:    runtimeconfigService,
		project:    project,
		configName: configName,
	}, nil
}

func (r *runtimeConfig) shouldRespectVariables() (bool, error) {
	diff := time.Now().Sub(r.lastCheck)

	// If 1 minute has elapsed since last check.
	if diff-(1*time.Minute) > 0 {
		name := fmt.Sprintf("projects/%s/configs/%s/variables/%s", r.project, r.configName, "respect-variables")
		value, err := r.service.Projects.Configs.Variables.Get(name).Do()
		if err != nil {
			return false, errors.Wrap(err, "failed to get variable")
		}

		log.Printf("variable value retrieved, value=%s", value.Text)
		r.lastCheck = time.Now()
		r.respectVariables = value.Text == "yes"
	} else {
		log.Println("using config variable from cache")
	}

	return r.respectVariables, nil
}
