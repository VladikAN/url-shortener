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
	Addr      string
	Ssl       bool
	Whitelist []string
}

// ServiceSettings holds service specific settings
type ServiceSettings struct {
	Chars string
}

// GetBase returns number of allowed chars
func (s *ServiceSettings) GetBase() int {
	return len(s.Chars)
}
