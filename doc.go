/*
Raccoon is a utility for simple orchestration of machines by using Dockerfile
syntax. You use ADD, RUN, MAINTAINER, etc. commands as you would do in Docker
but to provision hosts instead of containers using SSH. In this aspect, you
could need to have root privileges to a machine to do some specific actions

Raccoon uses JSON syntax to define its behaviour. This way you can generate
Raccoon jobs by simply writing a JSON with whichever language you want and
it provides a nice way to interact with the Raccoon server.

You can use Raccoon as a standalone app, running as server or as a library.
Ideally the three modes should work nicely but some errors are expected as
the heavy test of the package is done as a standalone app.
*/
package raccoon
