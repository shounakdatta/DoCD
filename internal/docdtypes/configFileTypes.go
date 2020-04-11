package docdtypes

// Command : Structure of service installation and build commands
type Command struct {
	Directory   string
	Command     string
	Environment []string
}

// Service : Structure of configuration services
type Service struct {
	ServiceName          string
	PackageManager       string
	Path                 string
	LogFilePath          string
	InstallationCommands []Command
	BuildCommands        []Command
}

// Config : Structure of DOCD-config.json
type Config struct {
	ProjectName        string
	BasePackageManager string
	InstallServices    bool
	Services           []Service
}

const (
	// ConfigFileName : Configuration file name
	ConfigFileName = "DoCD-config.json"
)
