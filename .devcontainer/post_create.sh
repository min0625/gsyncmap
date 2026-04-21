#!/bin/bash

# note that bash will read from ~/.profile or ~/.bash_profile if the latter exists
# ergo, you may want to check to see which is defined on your system and only append to the existing file
echo 'eval "$(mise activate bash --shims)"' >>~/.bash_profile # this sets up non-interactive sessions
echo 'eval "$(mise activate bash)"' >>~/.bashrc               # this sets up interactive sessions

mise trust .

mise install

mise exec -- go mod download

mise exec -- go install -v golang.org/x/tools/gopls@latest
mise exec -- go install -v github.com/go-delve/delve/cmd/dlv@latest
