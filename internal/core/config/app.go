package core_config

import (
	"fmt"
	"os"
	"time"
)

type App struct {
	TimeZone *time.Location
}

func newAppConfig() (App, error) {
	envVar := os.Getenv("TIME_ZONE")
	if envVar == "" {
		envVar = "UTC"
	}

	timeZone, err := time.LoadLocation(envVar)
	if err != nil {
		return App{}, fmt.Errorf(
			"load time zone: %s: %w",
			envVar,
			err,
		)
	}

	config := App{
		TimeZone: timeZone,
	}

	return config, nil
}
