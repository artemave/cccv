cccv
====

cccv finds what parts of a diff were copy/pasted from elsewhere in the project.

Why? Because copy/pasted code is impossible to spot at code review time.

## Usage

[Golang](http://golang.org/doc/install) is required to install this package.

```
% go get github.com/artemave/cccv
% git checkout pr1
% git diff master | cccv
```

For fine tuning, drop `.cccv.yml` into the root of your project. Example:
```
exclude-lines:
  - "require\\(['\"].+['\"]\\);$" # will exclude `var x = require('y');`
                                  # this regexp gets parsed by YAML first, hence extra escaping for \ and "
  - console.log

exclude-files:
  - "README.*" # regexp, NOT a glob

min-hunk-size: 3 # mininum number of consecutive duplicate
                 # lines to consider relevant; default 2

min-line-length: 15 # minimum line length (bar leading/trailing tabs and spaces)
                    # to be considered relevant; default 10
```

## Limitations

Relies on _default_ git diff output format.
