cccv
====

Check if git diff (commit, pr) contains copy pasted code.

## Usage

```
% go get github.com/artemave/cccv
% git checkout pr1
% git diff master | cccv
```

For fine tuning, drop `.cccv.yml` in the root of your project. Example:
```
exclude-lines:
  - "fmt|if err != nil"
  - WriteString

exclude-files:
  - "README.*" # this is regexp, NOT a glob

min-hunk-size: 3 # lines; defaults to 2

min-line-length: 15 # defaults to 10
```

## Limitations

Relies on _standard_ git diff format.
