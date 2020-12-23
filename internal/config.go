package internal

import (
	"gopkg.in/go-playground/validator.v9"
	"os"
)

const AlphaVantageApiKey = "ALPHA_VANTAGE_API_KEY"
const DBUsername = "DB_USERNAME"
const DBPassword = "DB_PASSWORD"

var config *EnvVars
var validate = validator.New()

type EnvVars struct {
	AlphaVantageKey string `validate:"required"`
	DBUsername      string `validate:"required"`
	DBPassword      string `validate:"required"`
}

func GetAppConfig() (*EnvVars, error) {
	if config == nil {
		config = &EnvVars{
			AlphaVantageKey: os.Getenv(AlphaVantageApiKey),
			DBUsername:      os.Getenv(DBUsername),
			DBPassword:      os.Getenv(DBPassword),
		}

		err := validate.Struct(config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	return config, nil
}
