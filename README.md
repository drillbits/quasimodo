Quasimodo
=========

Quasimodo execute command at a specified time.

## Requirements

- Go

## Installation

```
$ go get github.com/drillbits/quasimodo/qsmd
```

## Usage

### Run qsmd

```
$ qsmd [OPTIONS]
```

**Options**

- --conf="/etc/qsmd/qsmd.conf" Configuration file path

### Add task

Quasimodo provides http interface. You can add task by `curl`, etc.

```
$ curl http://127.0.0.1:1831/_qsmd/ -X POST -d "command=sleep 100000" -d "from=2014/05/23 12:06" -d "to=2014/05/23 12:09"
{"tasks":[{"from":"2014/05/23 12:06","to":"2014/05/23 12:09","command":"sleep 100000"}]}%
```

Run `sleep` after the specified `from` time.

```
$ ps ax|grep sleep|grep -v grep
(pid) (...) sleep 100000
```

And it ends after the specified `to` time.

```
$ ps ax|grep sleep|grep -v grep
```

## Show tasks

You can show task list by GET request.

```
$ curl http://127.0.0.1:1831/_qsmd/"
{"tasks":[{"from":"2014/05/23 12:06","to":"2014/05/23 12:09","command":"sleep 100000"}]}%
```

## LICENCE

MIT License.