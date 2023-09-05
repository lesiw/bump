# bump

`bump` bumps versions.

It accepts a version from standard input and prints the bumped version to
standard output.

`bump` is:

* A small, unix-like utility, designed for use in build scripts.
* Agnostic to the number of segments.
* Agnostic to version prefixes, like "v".
* Compatible with semantic versioning, but does not require it.

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
echo "1.0.0" | bump      # => 1.0.1
echo "v1.2" | bump       # => v1.3
echo "version 42" | bump # => version 43

# by index
echo "1.2.3" | bump -s 0 # => 2.0.0
echo "1.2.3" | bump -s 1 # => 1.3.0
echo "1.2.3" | bump -s 2 # => 1.2.4, same as default
echo "1.2.3" | bump -s 3 # => 1.2.4-rc.1

# semver aliases
echo "1.2.3" | bump -s major # => 2.0.0
echo "1.2.3" | bump -s minor # => 1.3.0
echo "1.2.3" | bump -s patch # => 1.2.4, same as default
echo "1.2.3" | bump -s pre   # => 1.2.4-rc.1
```

## Cookbook

### Change default segment to bump

By default, `bump` will bump the rightmost segment. The idiomatic way to
override this behavior is with an alias.

```sh
alias bump='bump -s 1' # bump the minor segment by default
```

### Bump git tag

```sh
git tag "$(git describe --abbrev=0 --tags | bump)"
git push origin --tags
```

### Bump tag based on latest commit keyword

If `+major`, `+minor`, `+patch`, or `+pre` are in the most recent commit, bump
the version according to the keyword, otherwise bump the patch segment.

```sh
SEGMENT=$(git show -s --format=%s | awk -F'+' 'BEGIN{RS=" "} /\+/ {print $2}')
bump -s "${SEGMENT:-patch}"
```
