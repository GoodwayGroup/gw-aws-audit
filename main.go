package main

// Import the package
import (
	"github.com/GoodwayGroup/gw-aws-audit/command"
	"github.com/yitsushi/go-commander"
)

func registerCommands(registry *commander.CommandRegistry) {
	// Register available commands
	registry.Register(command.NewVersion)
	registry.Register(command.NewListStoppedHostsCommand)
	registry.Register(command.NewListDetachedVolumesCommand)
	registry.Register(command.NewAddCostTagCommand)
	registry.Register(command.NewClearBucketObjectsCommand)
	registry.Register(command.NewBatchDeleteObjectsCommand)
}

// Main Section
func main() {
	registry := commander.NewCommandRegistry()

	registerCommands(registry)

	registry.Execute()
}