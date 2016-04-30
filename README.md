# raccoon
WIP App orchestration, configuration and deployment

## Try it

* Zombiebook example inside examples folder

* Prerequesites: Having 2 Centos Vagrant image running in 192.168.1.44 and 45

```bash
go build; ./raccoon zombiebook -z exampleBook.json -m exampleMansion.json
```

### Features
- [x] Pretty Output in real time
- [ ] Dockerfile Syntax to ease learning path. WIP
    - [x] RUN
    - [ ] ADD (WIP)
    - [ ] MAINTAINER (WIP)
- [ ] API REST.
- [x] Support for JSON syntax parsing
- [ ] Support for TOML syntax
- [x] CLI
- [ ] Automation tests
- [ ] Templating
- [ ] Target information retrieval
- [ ] Target "gathering facts"
- [ ] Identity file auth.

## Raccoon syntax
> In Raccoon, you'll have groups of remote commands (chapters of instructions)
with a group ID. Then you'll make group of hosts (rooms) with references to this
groups IDs (via chapters titles).

In Raccoon, a list of commands to execute in a remote machine is called a
Zombiebook. Syntax is similar to the one used for Docker in a Dockerfile.

In short, you'll write a **Zombiebook** which is composed with **chapters** which
are groups of **instructions** to survive a zombie apocalypse.

Then, you'll have a **Mansion** file that will be composed of many **rooms** each
of them associated with a **chapter** in the **zombiebook**.

A Mansion are one or many group of hosts where you want to execute some
commands. To know which command to execute, you associate the "chapter title" of
the group of hosts with the "chapter" of the book. A chapter on the book are
a group of commands.

### An example Zombiebook
Following Zombiebook shows all possible syntax. Every Zombiebook is an array of
**chapters** which instructions to run on grouped machines. A zombiebook must
start with a `chapter_title` and a `maintainer`. Then, a list of `instructions`
are passed as an array that will be parsed.

Each instruction as a different syntax but all have in common two parameters:
* **`name`**: as the command name like Dockerfile syntax (**`RUN`** at the moment)
* **`description`**: Information about the command to execute

```json
[
  {
    "chapter_title": "chapter1",
    "maintainer": "Burkraith",
    "instructions": [
      {
        "name": "RUN",
        "description": "Install htop",
        "instruction": "sudo yum install -y htop"
      },
      {
        "name": "ADD",
        "description": "Copying conf file",
        "sourcePath": "/tmp/asdfad",
        "destPath": "/tmp/folder"
      }
    ]
  },
  {
    "chapter_title": "chapter2",
    "maintainer": "Mario",
    "instructions": [
      {
        "name": "RUN",
        "description": "Install wget",
        "instruction": "sudo yum install -y wget"
      }
    ]
  }
]
```

## Mansion syntax
A mansion is a set of rooms (groups of remote hosts with a name). Each room
groups a set of remote hosts to execute commands on. Syntax is also easy and
straightforward:

```json
{
  "mansion_name":"A name",
  "rooms":[
    {
      "name":"some room",
      "chapter":"chapter1",
      "hosts":[
        {
          "ip":"192.168.1.44",
          "username":"vagrant",
          "password":"vagrant"
        }
      ]
    },
    {
      "name":"some room2",
      "chapter":"chapter2",
      "hosts":[
        {
          "ip":"192.168.1.45",
          "username":"vagrant",
          "password":"vagrant"
        }
      ]
    }
  ]
}
```

Each mansion has a **`mansion_name`** and an array of **`rooms`**. The Mansion
name is simply to identify the file and give it some description.

Each room has a **`name`**, a **`chapter`** of instructions to execute (that
will be taken from the **`zombiebook`**) and an array of hosts to execute to
commands on.

Each host has 3 parameters: **`ip`** of the host, **`username`** to access the
host and **`password`** to provide access to the host.