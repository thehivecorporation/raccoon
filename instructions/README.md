# Running tests on the `instructions` package

Unit tests in this package uses Docker to launch containers and try the SSH functions on them.

The nature of Docker forces you to launch the container using sudo privileges. This means that the tests must also be run using sudo:

```bash
# Change /usr/local/go to your GOROOT location
sudo /usr/local/go/bin/go test -v .
```