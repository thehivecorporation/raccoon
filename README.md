# raccoon
WIP App orchestration, configuration and deployment

![Raccoon logo](raccoon.jpg)

[![asciicast](https://asciinema.org/a/45363.png)](https://asciinema.org/a/45363)

### Features
- [x] Pretty Output in real time
- [ ] Dockerfile Syntax to ease learning path. WIP
    - [x] RUN
    - [x] ADD
    - [x] MAINTAINER
    - [x] ENV
- [ ] Pure Dockerfile parsing. Reuse your Dockerfiles in normal hosts.
- [ ] 1 file JSON parsing of jobs (a file with the infrastructure connected with the task list)
- [x] 2 files JSON parsing of jobs (one file with infrastructure pointing to tasks in a second file of tasks list).
- [ ] 3 files JSON parsing of jobs (a file for infrastructure, a file for tasks list and a file to connect them both).
- [x] Array based tasks for cluster
- [x] API REST.
- [x] Support for JSON syntax parsing
- [ ] Support for TOML syntax
- [x] CLI
- [ ] Automation tests
- [ ] Templating
- [ ] Target information retrieval
- [ ] Target "gathering facts"
- [ ] Identity file auth.

## Index
* [Raccoon CLI syntax](#raccoon-cli-syntax)
* [How to use Raccoon as a standalone app](#how-to-use-raccoon-as-a-standalone-app)
* [What's an infrastructure](#what's-an-infrastructure)
* [What's a task list](#what's-a-task-list)
* [Commands](#commands)
  * [RUN](#run)
  * [ENV](#env)
  * [ADD](#add)
  * [MAINTAINER](#maintainer)

## Raccoon CLI syntax
When you execute Raccoon from the command line without options, the following help will appear:
```
    NAME:
       raccoon job - Execute a job

    USAGE:
       raccoon job [command options] [arguments...]

    OPTIONS:
       --tasks value, -t value           Tasks file
       --infrastructure value, -i value  Infrastructure file
       --dispatcher value, -d value      Dispatching method between 3 options: sequential (no concurrent dispatch). simple (a Goroutine for each host) and worker_pool (a fixed number of workers) (default: "simple")
       --workersNumber value, -w value   In case of worker_pool dispath method, define the maximum number of workers here (default: 5)
```

## How to use Raccoon as a standalone app

You use Raccoon similar of how you use tools like Ansible. You mainly work with two things:

* **An Infrastructure**: An infrastructure is a list of clusters and a cluster is a list of hosts.
* **A Task list**: A task is a group of commands to achieve some task. You can have tasks of only one command (for example, 'apt-get install htop' if you just want to install `htop`) or a tasks composed of many commands (like a list of installation commands to take some host to a desired state).

It's quite easy, you just need to define your infrastructure and the tasks you want to execute in every cluster, then you launch the application.

## What's a Infrastructure: 
An infrastructure definition looks like the following:

```json
{
  "name":"Your infrastructure name (your project, for example)",
  "infrastructure":[
    {
      "name":"A nem for a cluster (like Kafka hosts)",
      "tasks":["Install Kafka", "Open Kafka ports"],
      "hosts":[
        {
          "ip":"172.17.42.1",
          "sshPort":32768,
          "username":"root",
          "password":"root"
        }
      ]
    }
  ]
}
```

* `name`: Is a name for your infrastructure, it is optional and its purpose is to identify the file between many infrastructure files. It could be your project's name or your company name.
* `infrastructure`: List of clusters. Take a closer look that is a JSON array. 
  * `name`: Name of the cluster. They should describe the cluster grouping in some way like Cassandra machines or QA machines.
  * `tasks`: List of tasks that will be executed on this cluster
  * `hosts`: List of hosts that compose this cluser.
    * `ip`: IP address to access the host.
    * `sshPort`: If needed, SSH port to make the SSH connection to.
    * `username`: Username to access the machine.
    * `password`: Password for the provided username.


## What's a Task list
A task is a group of commands that will be executed on the targeted host. A command refers, for example, to a shell command. See the commands index below.

Each task can be paired to a cluster, and **that task will be executed on the cluster**. The syntax of the task list is the following:

```json
[
  {
    "title": "task1",
    "maintainer": "Burkraith",
    "commands": [
      {
        "name": "ADD",
        "sourcePath": "doc.go",
        "destPath": "/tmp",
        "description": "Raccoon.go to /tmp"
      },
      {
        "name": "RUN",
        "description": "Removing htop",
        "instruction": "sudo yum remove -y htop"
      },
      {
        "name": "ADD",
        "sourcePath": "main.go",
        "destPath": "/tmp",
        "description": "copying raccoon.go to /tmp"
      }
    ]
  },
  {
    "title": "task2",
    "maintainer": "Mario",
    "commands": [
      {
        "name": "RUN",
        "description": "Removing htop",
        "instruction": "sudo apt-get remove -y htop"
      },
      {
        "name":"ADD",
        "description":"Adding the main.go file",
        "sourcePath":"main.go",
        "destPath":"/tmp"
      }
    ]
  }
]
```

First of all, take a close look because this JSON doesn't have root key. In this example we have two tasks and each task is composed of:

* `title`: This is a very important piece. The title of this task is referred from the insfrastructure definition. So if we have in a cluster:

```json
{
  "name":"myCluster",
  "tasks":["task1"]
  ...
}
```

The task with a title "task1" will be executed in cluster "myCluster". This is how Raccoon recognizes how to pair a task with a cluster.

* `maintainer`: Optional description for the maintainer of the task.
* `commands`: The list of commands that this task is composed of
  * `name`: Command name. This must be one of the commands described below. ADD, RUN, MAINTAINER are all valid commands.
  * `description`: An optional description to print to stdout when executing this command.
  * *specific key-values for each command*: All commands have a name and description in common. Then they have specific key-value pairs that contains the information to run the command. For example `RUN` command has the key `instruction` with the shell command you want to execute. Refer to the commands section for more information.

## Commands:

Until now we have developed the following commands with their corresponding
JSON formats:

### RUN

```json
      {
        "name": "RUN",
        "description": "Removing htop",
        "instruction": "sudo yum remove -y htop"
      }
```

Like in Docker, `RUN` will execute the command on the target machine, the
parameters are:
* name: The name of the command for the parser, for a `RUN` instruction is always `"RUN"`
* description: A description of the instruction. This has no effect and it is
only for logging purposes
* instruction: The command to execute in the target machine

### ADD

```json
      {
        "name": "ADD",
        "sourcePath": "raccoon.go",
        "destPath": "/tmp",
        "description": "copying raccoon.go to /tmp"
      }
```

`ADD` uses `scp` to send a file to the destination machine. It just supports
single files yet until we work in a folder solution (to send an entire folder):
* name: "ADD" must always be placed here so that the parser recognizes the
  instruction
* sourcePath: The full source path and file name of the file to send. For
  example: /tmp/fileToSend.log
* destPath: The full path to leave the file into the target machine. The file
  will have the same name.
* description: Optional description parameters for logging purposes.

### ENV

Sets an environment variable on the target machine:

```json
{
    "name": "ENV",
    "description": "Sets the variable GOROOT to /usr/local/go",
    "environment": "GOROOT=/usr/local/go"
}
```

The parser will look for a "=" in the "environment" value to split it into two
pieces and set the environment variable.
* name: "ENV" must always go here to use the ENV instruction
* description: Optional description parameters for logging purposes.
* environment: A key-value separated by a "=" to set in the target machine. Left
  side will be the environment name. Right side its value.

### MAINTAINER

Prints the name of a maintainer as a specific command:

```json
{
  "name":"MAINTAINER",
  "description":"Mario Castro"
}
```

* name: Must always contain "MAINTAINER" if you want to use this feature.
* description: the name of the maintainer.