#!/bin/bash

set -euo pipefail

echo 'eval "$(mise activate bash --shims)"' >>~/.bash_profile
echo 'eval "$(mise activate bash)"' >>~/.bashrc

mise trust .

mise install

mise exec -- go mod download

mise exec -- go install -v golang.org/x/tools/gopls@latest
mise exec -- go install -v github.com/go-delve/delve/cmd/dlv@latest

if [[ -f ".devcontainer/post_create.local.sh" ]]; then
    source ".devcontainer/post_create.local.sh"
fi
