This is a _living_ document, its purpose is to describe the tooling available
to `deluminator` contributors regarding testing, building and releasing new
versions.

# Makefile

Our [Makefile](/Makefile) is the entry point for most of the activities you
will run into as a contributor. To get a basic understanding of what you can
do with it, you can run:

```sh
make
```

As our default build is the help target, you should be seeing all the
documented targets the Makefile offers.

# Testing

We try to cover as much as we can with testing. The goal is having each single
feature covered by one or more tests.

## Running the tests

Once you are [setup](/CONTRIBUTING.md#setup-your-machine), you can run the test
suite with one command:

```sh
make test
```

You can run only a subset of the tests using the `TEST_PATTERN` variable:

```sh
make test TEST_PATTERN=TheAnswerIsFortyTwo
```

You can use `go test` through the `TEST_OPTIONS` variable:

```sh
make test TEST_OPTIONS=-v
```

You can combine the two options which is very helpful if you are working on a
specific feature and want immediate feedback.

## Golden files

## Use testutil
