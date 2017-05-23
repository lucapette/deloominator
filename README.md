# Deloominator

`deloominator` helps you explore your data with SQL queries. You can visualize
charts and group them in dashboard to share with your team. You can schedule
reports or just play around with your data. `deloominator` is in a very alpha
stage at the moment. You can have a look at [product
roadmap](https://github.com/lucapette/deloominator/projects/1) and at our
[milestones](https://github.com/lucapette/deloominator/milestones?direction=desc&sort=completeness&state=open)
to better understand what's the current status of the project.

- [Documentation](#documentation)
  - [Quick Start](#quick-start)
  - [User Manual](#user-manual)
  - [Demo](#demo)
- [Installation Guide](#installation-guide)
  - [Homebrew](#homebrew)
  - [Standalone](#standalone)
  - [Source Code](#standalone)
- [How To Contribute](#how-to-contribute)
- [Code Of Conduct](#code-of-conduct)

# Documentation

Latest version of our documentation is available [here](/docs).

## Quick start

Once you [installed](#installation-guide) `deloominator`, you can find a list of
the available options running the following command:

```sh
$ deloominator --help
```

`deloominator` uses only enviroment variables for configuration. The easiest way
to learn how to configure and run `deloominator` is reading our
[run.example](/bin/run.example).

## User Manual

Please refer to our [user manual](/docs/user-manual.md) to learn more about
`deloominator` features. **Note** that the document serves as an entry
point for the [product
roadmap](https://github.com/lucapette/deloominator/projects/1) at the moment.

## Demo

TODO

# Installation guide

`deloominator` has **zero** external dependencies, what you need is a binary for
your operating system and we offer multiple installation procedures.

## Homebrew

`deloominator` can be installed through Homebrew:

``` sh
$ brew tap lucapette/tap/deloominator
```

## Standalone

`deloominator` can be installed as an executable. Download the latest [compiled
binaries](https://github.com/lucapette/deloominator/releases) and put it
anywhere in your executable path.

## Source

Please refer to our [contributing guidelines](/CONTRIBUTING.md) to build and
install `deloominator` from the source.

# How to contribute

We welcome (and love) every form of contribution! Good entry points to the project are:

- Our [contributing guidelines](/CONTRIBUTING.md) document
- Our [developers' manual](/docs/developers-manual.md)
- Issues with the tag
  [gardening](https://github.com/lucapette/deloominator/issues?q=is%3Aissue+is%3Aopen+label%3Agardening)
- Issues with the tag [good first
  patch](https://github.com/lucapette/deloominator/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+patch%22)

If you're still not sure where to start, please open a [new
issue](https://github.com/lucapette/deloominator/issues/new) and we'll gladly
help you get started.

# Code of Conduct

You are expected to follow our [code of conduct](/CODE_OF_CONDUCT.md) when
interacting with the projects via issues, pull requests or in any other form.
Many thanks to the awesome [contributor
covenant](http://contributor-covenant.org/) initiative!

# License

[MIT License](/LICENSE) Copyright (c) [2017] [Luca Pette](http://lucapette.me)
