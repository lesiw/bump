# bump

`bump` bumps versions.

It accepts a version from standard input and prints the bumped version to
standard output.

`bump` is:

* A small, unix-like utility, designed for use in build scripts.
* Agnostic to the number of segments.
* Agnostic to version prefixes, like "v".

## Installation

```sh
wget -O /usr/local/bin/bump https://github.com/lesiw/bump/releases/latest/download/bump-$(uname -s)-$(uname -m)
chmod +x /usr/local/bin/bump
```

## Usage

```text
Usage of bump:
  -s string
        index of segment to bump
```

## Examples

```sh
# default
echo "1.0.0" | bump  # => 1.0.1
echo "v1.2" | bump  # => v1.3
echo "version 42" | bump  # => version 43

# by index
echo "1.2.3" | bump -s 0  # => 2.0.0
echo "1.2.3" | bump -s 1  # => 1.3.0
echo "1.2.3" | bump -s 2  # => 1.2.4, same as default

# semver aliases
echo "1.2.3" | bump -s major  # => 2.0.0
echo "1.2.3" | bump -s minor  # => 1.3.0
echo "1.2.3" | bump -s patch  # => 1.2.4, same as default
```

## Cookbook

### Bump git tag

```sh
git tag "$(git describe --abbrev=0 --tags | bump)"
git push origin --tags
```
