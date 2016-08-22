// Package raccoon is a utility for simple automation of hosts by using Dockerfile
// syntax. You use ADD, RUN, MAINTAINER, etc. commands as you would do in Docker
// but to provision hosts instead of containers using SSH. In this aspect, you
// could need to have root privileges in a machine to do some specific actions.
//
// Raccoon uses JSON syntax to define its behaviour. This way you can generate
// Raccoon jobs by simply writing a JSON with whichever language you want and
// it provides a nice way to interact with the Raccoon server.
//
// You can use Raccoon as a standalone app, running as server or as a library.
// Ideally the three modes should work nicely but some errors are expected as
// the heavy test of the package is done as a standalone app.
//
// How to use Raccoon as a library
//
// Raccoon can be easily used as a library. For each cluster and task, you need to
// create a hierarchy of Jobs. A Job is a type with two members: Cluster and Task.
// A Cluster is a list of Hosts and a
//
// How to user Raccoon as server
//
// Raccoon can also be used as server. Doing POST request you can launch some
// automation on some infrastructure.
package raccoon
