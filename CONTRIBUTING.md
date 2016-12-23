# Contributing

I love every form of contribution. By participating to this project, you agree
to abide to the deluminator [code of conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

Prerequisites for compilation are:

* `make`
* [Go 1.8+](http://golang.org/doc/install)

Clone `deluminator` from source:

```sh
$ git clone https://github.com/lucapette/deluminator.git
$ cd deluminator
```

Install the build dependencies:

``` sh
$ make setup
```

A good way of making sure everything is all right is running the test suite:

``` sh
$ make test
```

## Test your change

You can create a branch for your changes and try to build from the source as you go:

``` sh
$ make build
```

When you are satisfied with the changes, I suggest you run:

``` sh
$ make ci
```

Which runs all the tests and other checks to ensure code quality and consistency.

## Submit a pull request

Push your branch to your `deluminator` fork and open a pull request against
the master branch.
