Library to generate test models, run test gametype builds, and test games

### Setup Runner
The setuprunner directory contains a script that builds a test game type and runs test games. Must be run as a pod in the kubernetes cluster

### Bind data
To generate the bindata.go file, which embeds the test files into the go source code you must install go-bindata. Read more [here](https://github.com/jteeuwen/go-bindata)

The template data must be binded to distribute the lib without any non source file dependencies.

Anytime the test files are changed you must run:
```
go-bindata -pkg testutils testfiles/
```
