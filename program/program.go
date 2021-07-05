package program

type Program struct {
	Hosts []string `yaml:"hosts"`
	Tasks []Task   `yaml:"tasks"`
}

type Task struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
	Out  bool   `yaml:"out"`
}
