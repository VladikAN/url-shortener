package config

type config struct {
	Log     LogSettings
	Host    HostSettings
	Service ServiceSettings
}

// LogSettings holds log specific settings
type LogSettings struct {
	Level string
}

// HostSettings holds host specific settings
type HostSettings struct {
	Port      int
	Ssl       bool
	Whitelist []string
}

// ServiceSettings holds service specific settings
type ServiceSettings struct {
	Allowed string
}
