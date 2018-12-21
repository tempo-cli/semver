Semver
======

Semver is a Go implementation of [Composer's Semver parser](https://github.com/composer/semver).

## Usage

To use Semver, add it to your project

```bash
$ go get github.com/tempo-cli/semver
```

Then you can import the package in your project:

```go
import "gitub.com/tempo-cli/semver"
```

### Parsing versions

To parse a version, use the `NewVersion` function which will return a `Version` instance

```go

version := semver.NewVersion("v1.2.3")

fmt.PrintLn(version.Major()) // Prints '1'
fmt.PrintLn(version.Minor()) // Prints '2'
fmt.PrintLn(version.Patch()) // Prints '3'

```

You can parse various different version constraints, including branch names and normalizing values.


## TODO

 - [ ] Update documentation with more use cases
 - [ ] Compare versions
 - [ ] Parse constraints

## Contributing

If you would like to contribute to this project, create a [fork](https://help.github.com/articles/fork-a-repo/) and send a [pull request](https://help.github.com/articles/creating-a-pull-request/) with your changes.
Alternatively you can open an [issue](https://github.com/tempo-cli/semver/issues/) if something is not working as expectec
