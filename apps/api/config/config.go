// Package config provides an access point to env vars commonly use throught the app
package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Endpoint struct {
	URI string `env:"DEV_URI" json:",omitempty"`
}

type MongoDBConfig struct {
	URI string `env:"URI" json:",omitempty"`
	DB  string `env:"DB" json:",omitempty"`
}

type AWSConfig struct {
	Region               string `env:"REGION" json:",omitempty"`
	AccessKey            string `env:"ACCESS_KEY_ID" json:",omitempty"`
	SecretKey            string `env:"SECRET_ACCESS_KEY" json:",omitempty"`
	ExecutionEnv         string `env:"EXECUTION_ENV" json:",omitempty"`
	SubscriptionEndpoint string `env:"SUBSCRIPTION_ENDPOINT" json:",omitempty"`
}

type MailConfig struct {
	Host     string `env:"HOST" json:",omitempty"`
	Port     int    `env:"PORT" json:",omitempty"`
	Username string `env:"USERNAME" json:",omitempty"`
	Password string `env:"PASSWORD" json:",omitempty"`
	From     string `env:"FROM" json:",omitempty"`
}

type Config struct {
	Endpoint     *Endpoint      `env:", prefix=ENDPOINT_"`
	MongoDB      *MongoDBConfig `env:", prefix=MONGODB_" json:",omitempty"`
	AWS          *AWSConfig     `env:", prefix=AWS_" json:",omitempty"`
	Mail         *MailConfig    `env:", prefix=MAIL_" json:",omitempty"`
	APIAuthToken string         `env:"API_AUTH_TOKEN" json:",omitempty"`
}

func SetUp(cxt context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(cxt, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) MongoDBConfig() *MongoDBConfig {
	return c.MongoDB
}

func (c *Config) EndpointConfig() *Endpoint {
	return c.Endpoint
}

func (c *Config) AWSConfig() *AWSConfig {
	return c.AWS
}

func (c *Config) MailConfig() *MailConfig {
	return c.Mail
}
