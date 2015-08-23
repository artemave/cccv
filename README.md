cccv
====

cccv finds what parts of a diff were copy/pasted from elsewhere in the project.

Why? Because copy/pasted code is impossible to spot at code review time.

## Installation

You can get the binary:

    % curl -A "`uname -sp`" https://gox-server.herokuapp.com/artemave/cccv > cccv
    % chmod a+x ./cccv

Or, if you have [Golang](http://golang.org/doc/install) installed:

    % go get github.com/artemave/cccv

## Usage
```
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
