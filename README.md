# raccoon
WIP App orchestration, configuration and deployment

![Raccon logo](raccoon.jpg)

[![asciicast](https://asciinema.org/a/45363.png)](https://asciinema.org/a/45363)

### Features
- [x] Pretty Output in real time
- [ ] Dockerfile Syntax to ease learning path. WIP
    - [x] RUN
    - [x] ADD
    - [x] MAINTAINER
    - [x] ENV
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

```bash
    Raccon

    NAME:
       Raccoon - WIP Automation utility made easy with Dockerfile syntax

    USAGE:
       cli [global options] command [command options] [arguments...]

    VERSION:
       0.3.0

    COMMANDS:
         tasks   Execute a task list
         server  Launch a server to receive Commands JSON files
         show    Show special information about Raccoon

    GLOBAL OPTIONS:
       --help, -h     show help
       --version, -v  print the version
```

## Raccoon syntax

```bash
    NAME:
       cli tasks - Execute a task list

    USAGE:
       cli tasks [command options] [arguments...]

    OPTIONS:
       --tasks value, -t value           Tasks file
       --infrastructure value, -i value  Infrastructure file
       --dispatcher value, -d value      Dispatching method between 3 options: sequential (no concurrent dispatch). simple (a Goroutine for each host) and worker_pool (a fixed number of workers) (default: "simple")
       --workersNumber value, -w value   In case of worker_pool dispath method, define the maximum number of workers here (default: 5)
```

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

More info coming soon...