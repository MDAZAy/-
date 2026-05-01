if (-not $env:CONFIG_PATH) {
  $env:CONFIG_PATH = "config/config.yaml"
}
go run ./cmd/agent
