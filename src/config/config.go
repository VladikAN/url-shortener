package config

var cfg config

// Log gives log settings
func Log() *LogSettings {
	return &cfg.Log
}

// Host gives host settings
func Host() *HostSettings {
	return &cfg.Host
}

// Service gives service settings
func Service() *ServiceSettings {
	return &cfg.Service
}
