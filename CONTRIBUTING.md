# Contributing

We love every form of contribution. By participating to this project, you
agree to abide to the `deloominator` [code of conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`deloominator` is written in [Go](https://golang.org/).

Prerequisites are:

* Build:
  * `make`
  * [Go 1.8+](http://golang.org/doc/install) (of course! :))
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

``` sh
$ make setup
```

A good way of making sure everything is all right is running the test suite:

``` sh
$ make test
```

Please open an [issue](https://github.com/lucapette/deloominator/issues/new)
if you run into any problem.

A good starting point to learn more about `deloominator` code is our
[developers' manual](/docs/developers-manual.md).

## Test your change

You can create a branch for your changes and try to build from the source as
you go:

``` sh
$ make build
```

When you are satisfied with the changes, we suggest running:

``` sh
$ make ci
```

This command runs all the linters and runs all the tests.

## Submit a pull request

Push your branch to your `deloominator` fork and open a pull request against
the master branch.
