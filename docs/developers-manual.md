# Developer Manual

This is a _living_ document, its purpose is to describe the tooling available
to `deloominator` contributors regarding testing, building and releasing new
versions.

## Makefile

Our [Makefile](/Makefile) is the entry point for most of the activities you
will run into as a contributor. To get a basic understanding of what you can
do with it, you can run:

```sh
make help
```

Which shows all the documented targets.

## Testing

We try to cover as much as we can with testing. The goal is having each single
feature covered by one or more tests. Adding more tests is a great way of
contributing to the project!

### Running the tests

Once you are [setup](/CONTRIBUTING.md#setup-your-machine), you can run the test
suite with one command:

```sh
make test
```

You can run only a subset of the tests using the `TEST_PATTERN` variable:

```sh
make test TEST_PATTERN=TheAnswerIsFortyTwo
```

You can use pass options to `go test` through the `TEST_OPTIONS` variable:

```sh
make test TEST_OPTIONS=-v
```

You can combine the two options which is very helpful if you are working on a
specific feature and want immediate feedback. Like so:

```sh
make test TEST_OPTIONS=-v TEST_PATTERN=TheAnswerIsFortyTwo
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
