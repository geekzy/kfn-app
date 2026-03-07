package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config represents the configuration schema for the KFn platform
type Config struct {
	Registry        RegistryConfig        `mapstructure:"registry"`
	Container       ContainerConfig       `mapstructure:"container"`
	Sandbox         SandboxConfig         `mapstructure:"sandbox"`
	ImageRegistry   ImageRegistryConfig   `mapstructure:"image-registry"`
	Observability   ObservabilityConfig   `mapstructure:"observability"`
	Events          EventsConfig          `mapstructure:"events"`
	Auth            AuthConfig            `mapstructure:"auth"`
	Secrets         SecretsConfig         `mapstructure:"secrets"`
	Build           BuildConfig           `mapstructure:"build"`
	APIServer       APIServerConfig       `mapstructure:"api-server"`
	Scheduler       SchedulerConfig       `mapstructure:"scheduler"`
	RuntimeWrappers RuntimeWrappersConfig `mapstructure:"runtime-wrappers"`
}

// RegistryConfig holds etcd registry configuration
type RegistryConfig struct {
	Endpoints   []string      `mapstructure:"etcd-endpoints"`
	Timeout     time.Duration `mapstructure:"timeout"`
	DialTimeout time.Duration `mapstructure:"dial-timeout"`
	Prefix      string        `mapstructure:"prefix"`
}

// ContainerConfig holds container runtime configuration
type ContainerConfig struct {
	ContainerdSocket string `mapstructure:"containerd-socket"`
	Namespace        string `mapstructure:"namespace"`
	Snapshotter      string `mapstructure:"snapshotter"`
}

// SandboxConfig holds sandbox configuration
type SandboxConfig struct {
	Type string `mapstructure:"type"` // gVisor or Kata
}

// ImageRegistryConfig holds image registry configuration
type ImageRegistryConfig struct {
	HarborURL    string `mapstructure:"harbor-url"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Insecure     bool   `mapstructure:"insecure"`
	CertFile     string `mapstructure:"cert-file"`
	KeyFile      string `mapstructure:"key-file"`
	CARootFile   string `mapstructure:"ca-root-file"`
}

// ObservabilityConfig holds observability configuration
type ObservabilityConfig struct {
	PrometheusURL string `mapstructure:"prometheus-url"`
	LogLevel      string `mapstructure:"log-level"`
	LogFormat     string `mapstructure:"log-format"`
}

// EventsConfig holds event system configuration
type EventsConfig struct {
	KafkaBrokers []string `mapstructure:"kafka-brokers"`
	NATSServers  []string `mapstructure:"nats-servers"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string `mapstructure:"jwt-secret"`
	JWTIssuer     string `mapstructure:"jwt-issuer"`
	JWTAudience   string `mapstructure:"jwt-audience"`
	JWTExpiration string `mapstructure:"jwt-expiration"`
	RBACEnabled   bool   `mapstructure:"rbac-enabled"`
	MTLSEnabled   bool   `mapstructure:"mtls-enabled"`
	MTLSCertFile  string `mapstructure:"mtls-cert-file"`
	MTLSKeyFile   string `mapstructure:"mtls-key-file"`
	MTLSCAFile    string `mapstructure:"mtls-ca-file"`
}

// SecretsConfig holds secrets management configuration
type SecretsConfig struct {
	VaultAddress   string `mapstructure:"vault-address"`
	VaultToken     string `mapstructure:"vault-token"`
	VaultNamespace string `mapstructure:"vault-namespace"`
}

// BuildConfig holds image build configuration
type BuildConfig struct {
	Builder string `mapstructure:"builder"` // buildah or docker
	Context string `mapstructure:"context"`
}

// APIServerConfig holds API server configuration
type APIServerConfig struct {
	Address         string `mapstructure:"address"`
	Port            int    `mapstructure:"port"`
	TLSEnabled      bool   `mapstructure:"tls-enabled"`
	TLSCertFile     string `mapstructure:"tls-cert-file"`
	TLSKeyFile      string `mapstructure:"tls-key-file"`
	RequestTimeout  int    `mapstructure:"request-timeout"`
	ShutdownTimeout int    `mapstructure:"shutdown-timeout"`
}

// SchedulerConfig holds scheduler configuration
type SchedulerConfig struct {
	Concurrency     int           `mapstructure:"concurrency"`
	HealthCheckPort int           `mapstructure:"health-check-port"`
	WarmPoolSize    int           `mapstructure:"warm-pool-size"`
	IdleTimeout     time.Duration `mapstructure:"idle-timeout"`
}

// RuntimeWrappersConfig holds runtime wrapper configuration
type RuntimeWrappersConfig struct {
	NodeJSPath   string `mapstructure:"nodejs-path"`
	PythonPath   string `mapstructure:"python-path"`
	Entrypoint   string `mapstructure:"entrypoint"`
	WrapperPath  string `mapstructure:"wrapper-path"`
	Dependencies string `mapstructure:"dependencies"`
}

var (
	configInstance *Config
	viperInstance  *viper.Viper
)

// Init initializes the configuration system
func Init(configPaths ...string) error {
	viperInstance = viper.New()

	// Set default values
	setDefaults(viperInstance)

	// Set config file name and type
	viperInstance.SetConfigType("yaml")
	viperInstance.SetConfigName("default")

	// Add config paths
	if len(configPaths) == 0 {
		viperInstance.AddConfigPath(".")
		viperInstance.AddConfigPath("./configs")
		viperInstance.AddConfigPath("/etc/kfn/")
	} else {
		for _, path := range configPaths {
			viperInstance.AddConfigPath(path)
		}
	}

	// Enable environment variable support
	viperInstance.AutomaticEnv()
	viperInstance.SetEnvPrefix("KFN")
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Read config file
	if err := viperInstance.ReadInConfig(); err != nil {
		// If config file is not found, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Watch config file for changes
	viperInstance.WatchConfig()
	viperInstance.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
		if err := viperInstance.Unmarshal(&configInstance); err != nil {
			fmt.Printf("Failed to unmarshal config: %v\n", err)
		}
	})

	// Unmarshal config into struct
	configInstance = &Config{}
	if err := viperInstance.Unmarshal(configInstance); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Get returns the current configuration instance
func Get() *Config {
	if configInstance == nil {
		panic("config not initialized, call Init() first")
	}
	return configInstance
}

// GetViper returns the viper instance
func GetViper() *viper.Viper {
	if viperInstance == nil {
		panic("config not initialized, call Init() first")
	}
	return viperInstance
}

// Reload reloads the configuration from file
func Reload() error {
	if viperInstance == nil {
		return fmt.Errorf("config not initialized, call Init() first")
	}

	if err := viperInstance.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viperInstance.Unmarshal(configInstance); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Registry defaults
	v.SetDefault("registry.etcd-endpoints", []string{"localhost:2379"})
	v.SetDefault("registry.timeout", "5s")
	v.SetDefault("registry.dial-timeout", "5s")
	v.SetDefault("registry.prefix", "/kfn")

	// Container defaults
	v.SetDefault("container.containerd-socket", "/run/containerd/containerd.sock")
	v.SetDefault("container.namespace", "kfn")
	v.SetDefault("container.snapshotter", "overlayfs")

	// Sandbox defaults
	v.SetDefault("sandbox.type", "gVisor")

	// Image registry defaults
	v.SetDefault("image-registry.harbor-url", "http://localhost:8080")
	v.SetDefault("image-registry.insecure", false)

	// Observability defaults
	v.SetDefault("observability.prometheus-url", "http://localhost:9090")
	v.SetDefault("observability.log-level", "info")
	v.SetDefault("observability.log-format", "json")

	// Events defaults
	v.SetDefault("events.kafka-brokers", []string{"localhost:9092"})
	v.SetDefault("events.nats-servers", []string{"nats://localhost:4222"})

	// Auth defaults
	v.SetDefault("auth.jwt-secret", "kfn-jwt-secret")
	v.SetDefault("auth.jwt-issuer", "kfn-platform")
	v.SetDefault("auth.jwt-audience", "kfn-users")
	v.SetDefault("auth.jwt-expiration", "24h")
	v.SetDefault("auth.rbac-enabled", false)
	v.SetDefault("auth.mtls-enabled", false)

	// Secrets defaults
	v.SetDefault("secrets.vault-address", "http://localhost:8200")

	// Build defaults
	v.SetDefault("build.builder", "buildah")
	v.SetDefault("build.context", ".")

	// API server defaults
	v.SetDefault("api-server.address", "0.0.0.0")
	v.SetDefault("api-server.port", 8080)
	v.SetDefault("api-server.tls-enabled", false)
	v.SetDefault("api-server.request-timeout", 30)
	v.SetDefault("api-server.shutdown-timeout", 30)

	// Scheduler defaults
	v.SetDefault("scheduler.concurrency", 10)
	v.SetDefault("scheduler.health-check-port", 8081)
	v.SetDefault("scheduler.warm-pool-size", 5)
	v.SetDefault("scheduler.idle-timeout", "300s")

	// Runtime wrappers defaults
	v.SetDefault("runtime-wrappers.nodejs-path", "/usr/bin/node")
	v.SetDefault("runtime-wrappers.python-path", "/usr/bin/python3")
	v.SetDefault("runtime-wrappers.entrypoint", "entrypoint.sh")
	v.SetDefault("runtime-wrappers.wrapper-path", "wrapper.js")
	v.SetDefault("runtime-wrappers.dependencies", "package.json")
}