package config

import (
	"fmt"
	"regexp"
	"time"

	validator "github.com/theflyingcodr/govalidator"
)

// Environment variable constants.
const (
	EnvServerPort           = "server.port"
	EnvServerHost           = "server.host"
	EnvServerFQDN           = "server.fqdn"
	EnvServerSwaggerEnabled = "server.swagger.enabled"
	EnvEnvironment          = "env.environment"
	EnvBitcoinNetwork       = "env.bitcoin.network"
	EnvRegion               = "env.region"
	EnvVersion              = "env.version"
	EnvCommit               = "env.commit"
	EnvBuildDate            = "env.builddate"
	EnvLogLevel             = "log.level"
	EnvPaydHost             = "payd.host"
	EnvPaydPort             = "payd.port"
	EnvPaydSecure           = "payd.secure"
	EnvPaydCertPath         = "payd.cert.path"
	EnvPaydNoop             = "payd.noop"

	LogDebug = "debug"
	LogInfo  = "info"
	LogError = "error"
	LogWarn  = "warn"
)

var reNetworks = regexp.MustCompile(`^(regtest|stn|testnet|mainnet)$`)

// NetworkType is used to restrict the networks we can support.
type NetworkType string

// Supported bitcoin network types.
const (
	NetworkRegtest NetworkType = "mainnet"
	NetworkSTN     NetworkType = "stn"
	NetworkTestnet NetworkType = "testnet"
	NetworkMainet  NetworkType = "mainnet"
)

// Config returns strongly typed config values.
type Config struct {
	Logging    *Logging
	Server     *Server
	Deployment *Deployment
	PayD       *PayD
}

// Validate will ensure the config matches certain parameters.
func (c *Config) Validate() error {
	vl := validator.New()
	if c.Deployment != nil {
		vl = vl.Validate("deployment.network", validator.MatchString(string(c.Deployment.Network), reNetworks))
	}
	return vl.Err()
}

// Deployment contains information relating to the current
// deployed instance.
type Deployment struct {
	Environment string
	AppName     string
	Region      string
	Version     string
	Commit      string
	BuildDate   time.Time
	Network     NetworkType
}

// IsDev determines if this app is running on a dev environment.
func (d *Deployment) IsDev() bool {
	return d.Environment == "dev"
}

func (d *Deployment) String() string {
	return fmt.Sprintf("Environment: %s \n AppName: %s\n Region: %s\n Version: %s\n Commit:%s\n BuildDate: %s\n",
		d.Environment, d.AppName, d.Region, d.Version, d.Commit, d.BuildDate)
}

// Logging contains log configuration.
type Logging struct {
	Level string
}

// Server contains all settings required to run a web server.
type Server struct {
	Port     string
	Hostname string
	// FQDN - fully qualified domain name, used to form the paymentRequest
	// payment URL as this may be different from the hostname + port.
	FQDN string
	// SwaggerEnabled if true we will include an endpoint to serve swagger documents.
	SwaggerEnabled bool
}

// PayD is used to setup connection to a payd instance.
// In this case, we connect to only one merchant wallet
// implementors may need to connect to more.
type PayD struct {
	Host            string
	Port            string
	Secure          bool
	CertificatePath string
	Noop            bool
}

// ConfigurationLoader will load configuration items
// into a struct that contains a configuration.
type ConfigurationLoader interface {
	WithServer() ConfigurationLoader
	WithDeployment(app string) ConfigurationLoader
	WithLog() ConfigurationLoader
	WithPayD() ConfigurationLoader
	Load() *Config
}
