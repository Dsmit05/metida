package config

import (
	"encoding/json"
	"fmt"
	"github.com/Dsmit05/metida/internal/logger"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	// buildVersion sets on compile time.
	buildVersion = ""
)

// Database - contains all parameters database connection.
type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Table    string `yaml:"table"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// ApiServer - contains parameter for rest connection.
type ApiServer struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

// DebagServer - contains parameter for debag server.
type DebagServer struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

// CORS - contains parameter for cors settings.
type CORS struct {
	AllowedOrigins []string `yaml:"allowedOrigins"`
}

// Cryptography - contains secret for jwt token.
type Cryptography struct {
	Secret string `yaml:"secret"`
}

// Project - contains all parameters project information.
type Project struct {
	BuildVersion string
}

// Config contains all necessary params.
type Config struct {
	Database     Database     `yaml:"database"`
	ApiServer    ApiServer    `yaml:"apiServer"`
	DebagServer  DebagServer  `yaml:"debagServer"`
	CORS         CORS         `yaml:"cors"`
	Cryptography Cryptography `yaml:"cryptography"`
	Project
	CommandLineI
}

type CommandLineI interface {
	IfDebagOn() bool
}

// NewConfig return Config from...
func NewConfig(flagCmd CommandLineI) (*Config, error) {
	var cfg = new(Config)

	if flagCmd.IfDebagOn() {
		if err := cfg.initFromFile("config.yml"); err != nil {
			return nil, err
		}
	} else {
		// for the Prod version, we fill in the config from a safe place, for example from the consul
		// Todo: you will need to change:
		if err := cfg.initFromFile("config.yml"); err != nil {
			return nil, err
		}
	}

	cfg.CommandLineI = flagCmd
	cfg.Project.BuildVersion = buildVersion

	return cfg, nil
}

// initFromFile init Config from yml file.
func (o *Config) initFromFile(filePath string) error {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("config file close", err)
		}
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&o); err != nil {
		return err
	}

	return nil
}

// String representation Config settings.
func (o Config) String() string {
	return fmt.Sprintf(" BuildVersion: %+v\n Database: %+v\n ApiServer: %+v\n DebagServer: %+v\n",
		o.BuildVersion, o.Database, o.ApiServer, o.DebagServer)
}

func (o *Config) GetConnectDB() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		o.Database.User,
		o.Database.Password,
		o.Database.Host,
		o.Database.Port,
		o.Database.Table,
	)
}

// GetApiAddr return addr in format localhost:8080
func (o *Config) GetApiAddr() string {
	return fmt.Sprintf("%v:%v", o.ApiServer.Host, o.ApiServer.Port)
}

// GetApiReadTimeout in second.
func (o *Config) GetApiReadTimeout() time.Duration {
	timeout := time.Duration(o.ApiServer.ReadTimeout) * time.Second
	return timeout
}

// GetApiWriteTimeout in second.
func (o *Config) GetApiWriteTimeout() time.Duration {
	timeout := time.Duration(o.ApiServer.WriteTimeout) * time.Second
	return timeout
}

// GetDebagAddr return addr debag server.
func (o *Config) GetDebagAddr() string {
	return fmt.Sprintf("%v:%v", o.DebagServer.Host, o.DebagServer.Port)
}

// GetDebagReadTimeout in second.
func (o *Config) GetDebagReadTimeout() time.Duration {
	timeout := time.Duration(o.DebagServer.ReadTimeout) * time.Second
	return timeout
}

// GetDebagWriteTimeout in second.
func (o *Config) GetDebagWriteTimeout() time.Duration {
	timeout := time.Duration(o.DebagServer.WriteTimeout) * time.Second
	return timeout
}

// GetCorsAllowedOrigins return list of origins a cross-domain request can be executed from.
func (o *Config) GetCorsAllowedOrigins() []string {
	o.CORS.AllowedOrigins = append(o.CORS.AllowedOrigins, "http://"+o.GetDebagAddr())
	return o.CORS.AllowedOrigins
}

// GetConfigInfo handler info build.
func (o *Config) GetConfigInfo(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"BuildVersion": o.BuildVersion,
		"debug":        o.IfDebagOn(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Service information encoding error", err)
	}
}
