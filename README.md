# raccoon
WIP App orchestration, configuration and deployment

## Try it

*Zombiebook example inside examples folder

*Prerequesites: Having 1 Centos Vagrant image running in 192.168.33.10

```bash
go build; ./raccoon -file <PATH_TO_YOUR_ZOMBIEBOOK_FILE>
```

### Features
* Pretty Output in real time
* Dockerfile Syntax to ease learning path
* API REST
* Support for JSON syntax parsing. Nice YAML and TOML
* CLI

### TODO
* Automation tests
* JSON Parsing
* Templating
* Target information retrieval
* Target "gathering facts"
* CLI with codegansta