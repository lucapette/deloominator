# Developers' manual

This is a _living_ document, its purpose is to describe the tooling available
to `deloominator` contributors regarding testing, building and releasing new
versions.

## Set up your machine

Our [Makefile](/Makefile) is the entry point for most of the activities you
will run into as a contributor. To get a basic understanding of what you can
do with it, you can run:

```sh
$ make help
```

`deloominator` is written in [Go](https://golang.org/) and JavaScript. Here is a
list of prerequisites to build and test the project:

- Build:
  - `make`
  - [Go 1.9+](http://golang.org/doc/install)
  - [Node.js](https://nodejs.org/en/)
  - [npm](https://www.npmjs.com/)
- Test:
  - [PostgreSQL](https://www.postgresql.org/)
  - [MySQL](https://www.mysql.com/)

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
$ cp bin/run-test{.example,}
$ $EDITOR bin/run-test # edit according to your needs
$ make test
```

Please open an [issue](https://github.com/lucapette/deloominator/issues/new)
if you run into any problem.

## Building and running deloominator

`deloominator` has two main components:

- An API server written in [Go](https::/golang.org)
- A UI written in JavaScript that uses
  [React](https://facebook.github.io/react/)

You can build the entire application by running `make` without arguments

```sh
make
```

since `build` is the default target. Both components have their own building
target:

```sh
$ make build-server # builds the API server
$ make build-ui # builds the UI application
```

You can run `deloominator` following the steps:

```sh
$ make
$ ./deloominator
```

These steps will do a full build of the application, so it's a slower process.
If you're working with only one of the components, you can use a more
specialized command:

```sh
$ make run-server # runs the API server
$ make run-ui # runs the UI application
```

Please *note* that both `make run-server` and `make run-ui` rely on a small bash
script called `./bin/run`. There is an example file [here](bin/run.example). You
can create your own script with the following command:

```sh
cp bin/run{.example,}
```

and then edit it according to your local configuration.

## Testing

We try to cover with testing as much as we can. The goal is to have each single
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
specific feature and want immediate feedback:

```sh
$ make test TEST_OPTIONS=-v TEST_PATTERN=TheAnswerIsFortyTwo
```

Moreover, you can run tests only for the server:

```sh
$ make test-server
```

or the frontend tests:

```sh
$ make test-ui
```

### Golden files

Our test suite has fixture files that facilitate both the test setup (as in the
case of the SQL scripts to create databases) and test assertion. Golden files
is a technique to handle fixture files used in assertion. It works this way:

- You store complex output you expect (like a JSON response for example) in a
  file and use it to compare it to the actual outcome of a test
- You add a command line flag that updates your golden file so that it is easy
  to get a test passing when behavior changes

You can find an example of it [here](/pkg/api/graphql_test.go).

### Use testutil

The [testutil](pkg/testutil) package contains a number of utilities for testing
our Go code. In the spirit of [Advanced testing with
Go](https://speakerdeck.com/mitchellh/advanced-testing-with-go), we follow these
guidelines:

- Each public function takes `t *testing.T` as a first parameter
- Functions that use temporary resources (like databases, files, and so on)
  return a callback function so that the caller can decide when and how to
  clean up test execution
- For more complex tests, we use fixtures. Add new fixtures in
  [testutil/fixtures](pkg/testutil/fixtures) so that existing helpers can load
  them.
