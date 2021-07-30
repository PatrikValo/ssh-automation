# ssh-automation
![example workflow](https://github.com/PatrikValo/ssh-automation/actions/workflows/testing.yaml/badge.svg)

Simple, fast and lightweight automation tool for execution bash commands on multiple remote machines. It is inspired
by [ansible](https://github.com/ansible/ansible) and written in Go. For communication with remote machines it uses SSH.
It executes commands, which are defined in yaml file with pre-defined syntax.

This project was created for purpose of learning the Go language.

## YAML file syntax

```yaml
hosts:
  - "172.26.2.109"
  - "172.26.2.110"
  - "172.26.2.111"
  - "172.26.2.112:8080"
tasks:
  - name: "Show all"
    cmd: "ls -a"
    out: true
  - name: "Install python3"
    cmd: "yum install -y python3"
  - name: "Check python3"
    cmd: "which python3"
    out: true
```

Each file must contain `hosts` and `tasks` attributes.

- `hosts` - contains array of addresses of the machines. Port is optional and default port is **22**.

- `tasks` - contains array of the tasks, which will be executed on the machines from the `hosts` array. Each task
  contains:
    - `name` - name of the task, which is printed during execution on terminal
    - `cmd` - correct bash command
    - `out` - bool if you want to see the output of the command on the terminal

## Running app

```shell
go mod download
go run main.go [-flags]
```

### flags

- `-f [playbook_path]` - path to source file with commands
- `-u [user_name]` - user login, if is not provided **root** is used
- `-p` - authentication by password. (password is required in the next step)
- `-k [path_to_id_rsa]` - authentication by private key. If is not provided **~/.ssh/id_rsa** is used
- `-h` flag for help

## Output

Example of the successful output for above yaml file

```shell
TASK [Try to connect machines]
*******************************************
|-+ SUCCESS: 4
|-+ FAIL: 0

TASK [Show all]
*******************************************
HOST: [172.26.2.110:22]
.   anaconda-ks.cfg  .bash_profile  .cshrc           .ssh
..  .bash_logout     .bashrc        original-ks.cfg  .tcshrc

HOST: [172.26.2.111:22]
.   anaconda-ks.cfg  .bash_profile  .cshrc           .ssh
..  .bash_logout     .bashrc        original-ks.cfg  .tcshrc

HOST: [172.26.2.109:22]
.   anaconda-ks.cfg  .bash_profile  .cshrc           .ssh
..  .bash_logout     .bashrc        original-ks.cfg  .tcshrc

HOST: [172.26.2.112:22]
.   anaconda-ks.cfg  .bash_profile  .cshrc           .ssh
..  .bash_logout     .bashrc        original-ks.cfg  .tcshrc

|-+ SUCCESS: 4
|-+ FAIL: 0

TASK [Install python3]
*******************************************
|-+ SUCCESS: 4
|-+ FAIL: 0

TASK [Check python3]
*******************************************
HOST: [172.26.2.110:22]
/usr/bin/python3

HOST: [172.26.2.112:22]
/usr/bin/python3

HOST: [172.26.2.111:22]
/usr/bin/python3

HOST: [172.26.2.109:22]
/usr/bin/python3

|-+ SUCCESS: 4
|-+ FAIL: 0
```

Example of the failed output

```shell
TASK [Try to connect machines]
*******************************************
HOST: [172.26.2.113:22]
dial tcp 172.26.2.113:22: connect: no route to host

|-+ SUCCESS: 3
|-+ FAIL: 1
```

## Contributing

### Commit messages convention

`feat: add new feature`

`ref: refactor the chunk of the code`

`test: add test`

`doc: change documentation`
