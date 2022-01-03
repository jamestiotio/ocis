package command

import (
	"github.com/owncloud/ocis/graph-explorer/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GraphExplorerCommand is the entrypoint for the graph-explorer command.
func GraphExplorerCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "graph-explorer",
		Usage:    "Start graph-explorer server",
		Category: "Extensions",
		Before: func(ctx *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Subcommands: command.GetCommands(cfg.GraphExplorer),
	}
}

func init() {
	register.AddCommand(GraphExplorerCommand)
}
