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

For more information on suported versions and constraints, visit the [Composer documentation](https://getcomposer.org/doc/articles/versions.md) 

### Parsing versions

To parse a version, use the `NewVersion` function which will return a `Version` instance

```go

version := semver.NewVersion("v1.2.3")

fmt.PrintLn(version.Major) // Prints '1'
fmt.PrintLn(version.Minor) // Prints '2'
fmt.PrintLn(version.Patch) // Prints '3'

```

You can parse various different version constraints, including branch names and normalizing values.

### Parsing constraints

To parse a version constraint, use the `NewConstraint` function which will return an array of `Constraint` structs which containts the lower and upper bound for the constraint

```go

constraint := semver.NewConstraint("^2.0")

fmt.PrintLn(constraint[0].String) // Prints '>= 2.0.0.0-dev'
fmt.PrintLn(constraint[1].String) // Prints '< 3.0.0.0-dev'

```

## TODO

 - [ ] Update documentation with more use cases
 - [ ] Compare versions

## Contributing

If you would like to contribute to this project, create a [fork](https://help.github.com/articles/fork-a-repo/) and send a [pull request](https://help.github.com/articles/creating-a-pull-request/) with your changes.
Alternatively you can open an [issue](https://github.com/tempo-cli/semver/issues/) if something is not working as expectec
