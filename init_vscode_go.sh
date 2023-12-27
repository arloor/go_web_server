#参考 https://github.com/devcontainers/features/blob/main/src/go/install.sh#L54

# Install Go tools that are isImportant && !replacedByGopls based on
# https://github.com/golang/vscode-go/blob/v0.38.0/src/goToolsInformation.ts
GO_TOOLS="\
    golang.org/x/tools/gopls@latest \
    honnef.co/go/tools/cmd/staticcheck@latest \
    golang.org/x/lint/golint@latest \
    github.com/mgechev/revive@latest \
    github.com/go-delve/delve/cmd/dlv@latest \
    github.com/fatih/gomodifytags@latest \
    github.com/haya14busa/goplay/cmd/goplay@latest \
    github.com/cweill/gotests/gotests@latest \ 
    github.com/josharian/impl@latest"
(echo "${GO_TOOLS}" | xargs -n 1 go install -v )2>&1 | tee -a ./init_go.log

echo "Installing golangci-lint latest..."
curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b "$HOME/go/bin"