#!/usr/bin/env bash

rm -rf completions
mkdir completions
go run . completion bash > completions/shopwell-cli.bash
go run . completion zsh > completions/shopwell-cli.zsh
go run . completion fish > completions/shopwell-cli.fish