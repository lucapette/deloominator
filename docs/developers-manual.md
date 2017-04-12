# Developers' Manual

This is a _living_ document, its purpose is to describe the tooling available
to `deloominator` contributors regarding testing, building and releasing new
versions.

## Setup your machine

Our [Makefile](/Makefile) is the entry point for most of the activities you
will run into as a contributor. To get a basic understanding of what you can
do with it, you can run:

```sh
$ make help
```

Which shows all the documented targets. `deloominator` is written in
[Go](https://golang.org/) and JavaScript. Here is a list of prerequisites to
build and test the code:

* Build:
  * `make`
  * [Go 1.8+](http://golang.org/doc/install)
  * [Yarn](https://yarnpkg.com/en/)
  * [Node.js](https://nodejs.org/en/)
* Test:
  * [PostgreSQL](https://www.postgresql.org/)
  * [MySQL](https://www.mysql.com/)
  * [PhantomJS](http://phantomjs.org/)

Clone `deloominator` from source:

```sh
$ git clone https://github.com/lucapette/deloominator.git
$ cd deloominator
```

Install the build and lint dependencies:

```sh
$ make setup
```

A good way of making sure everything is all right is running the test suite:

```sh
$ make test
```

Please open an [issue](https://github.com/lucapette/deloominator/issues/new)
if you run into any problem.

## Building and running deloominator

`deloominator` has two main components:

- An API server written in [Go](https::/golang.org)
- A UI written in JavaScript that uses [React](https://facebook.github.io/react/)

You can build the entire application by running `make` without arguments:

```sh
make
```

since `build` is the default target. Both components have their own
building target:

```sh
$ make build-api # builds the API server
$ make build-ui # builds the UI application
```

You can run `deloominator` following the steps:

```sh
$ make
$ ./deloominator
```

These steps will do a full build of the application so it's a slower process.
If you're working only with one of the components, you can use a more
specialised command:

```sh
$ make run-api # runs the API server
$ make run-ui # runs the UI application
```

Please *note* that `make run-api` relies on a small bash script called
`bin/run`. There is an example file [here](bin/run.example). You can create
your own script with the following command:

```sh
cp bin/run{.example,}
```

And then edit it accordingly to your local configuration.

## Testing

We try to cover as much as we can with testing. The goal is having each single
feature covered by one or more tests. Adding more tests is a great way of
contributing to the project!

### Running the tests

Once you are [setup](#setup-your-machine), you can run the test suite with one
command:

```sh
$ make test
```

You can run only a subset of the tests using the `TEST_PATTERN` variable:

```sh
$ make test TEST_PATTERN=TheAnswerIsFortyTwo
```

You can use pass options to `go test` through the `TEST_OPTIONS` variable:

```sh
$ make test TEST_OPTIONS=-v
```

You can combine the two options which is very helpful if you are working on a
specific feature and want immediate feedback. Like so:

```sh
$ make test TEST_OPTIONS=-v TEST_PATTERN=TheAnswerIsFortyTwo
```

### Golden files

Our test suite has fixtures files that facilitate both test setup, as in the
case of the SQL scripts to create databases, and test assertion. Golden files
are a technique to handle fixtures files used in assertion. It works this way:

- You store complex output you expect (like a JSON response for example) in a
  file and use it to compare it to the actual outcome of a test
- You add a command line flag that updates your golden file so that it is easy
  to get a test passing when behaviour changes

You can find an example of it [here](/pkg/api/graphql_test.go).

### Use testutil
