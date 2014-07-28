cccv
====

Check if git diff (commit, pr) contains copy pasted code.

Why? Copy/pasted code is impossible to spot at code review time (unless the reviewer knows thy code very well). At the same time, it is such an obvious offence that it is more often than not caused by haste and forgetfulness rather than an evil intent (or so I'd like to believe). If caught early, it is unlikely spark an argument and has higher chances to get fixed. Ergo, worth pointing out.

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
