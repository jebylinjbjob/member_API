set -e

just
just build
just default
just install-tools
just test
just test-integration
just test-all
just vet
