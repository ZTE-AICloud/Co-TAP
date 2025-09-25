package types

// logger config
type Logger struct {
	Console ConsoleOutput
	File    FileOutput
}

type ConsoleOutput struct {
	Level string
}

type FileOutput struct {
	Filename string
	Level    string
	Maxlines int
	Maxsize  int
	Daily    bool
	Maxdays  int64
	Rotate   bool `env:"LOG_FILE_ROTATE"`
	Perm     string
}

type HealthCheckResult struct {
	Status string `json:"status"`
}
