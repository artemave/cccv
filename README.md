cccv
====

cccv finds what changes in a diff were copy/pasted from elsewhere in the project.

Why? Copy/pasted code is impossible to spot at code review time (unless the reviewer knows thy code very well). At the same time, it is such an obvious offence that it is more often than not caused by haste and forgetfulness rather than an evil intent (or so I'd like to believe). If caught early, it is unlikely spark an argument and has higher chances to get fixed. Ergo, worth pointing out.

## Usage

[Golang](http://golang.org/doc/install) is required to install the package.

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

# mininum number of consequtive duplicate lines to consider relevant; default 2
min-hunk-size: 3

# minimum line length (bar leading/trailing tabs and spaces) to be considered relevant; default 10
min-line-length: 15
```

## Limitations

Relies on _default_ git diff output format.
