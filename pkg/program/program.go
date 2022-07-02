package program

type Program struct {
	Hosts  []string `yaml:"hosts"`
	Tasks  []Task   `yaml:"tasks"`
	Config Config   `yaml:"config"`
}

type Task struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
	Out  bool   `yaml:"out"`
}

type Config struct {
	ConnectionTimeout int `yaml:"connection_timeout"`
}
