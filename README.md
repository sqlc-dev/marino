# Marino - A MySQL Compatible SQL Parser

## TODO

<!-- TODO: fill in fork-specific details (origin of the fork, scope of changes, maintainers, support channels, etc.) -->

## About

The goal of this project is to build a Golang parser that is fully compatible with MySQL syntax, easy to extend, and high performance. Currently, features supported by parser are as follows:

- Highly compatible with MySQL: it supports almost all features of MySQL.
- Extensible: the grammar is hand-written recursive descent — adding a new syntax is a few lines of ordinary Go.
- Good performance: a single-pass hand-written parser with a streaming lexer and no parser-generator runtime.

## Acknowledgments

Marino is a hard fork of [pingcap/parser](https://github.com/pingcap/parser),
the MySQL-compatible SQL parser developed as part of [TiDB](https://github.com/pingcap/tidb).
Huge thanks to the TiDB team and all the upstream contributors whose work this
project is built on.

## License

Marino is under the Apache 2.0 license. See the LICENSE file for details.
