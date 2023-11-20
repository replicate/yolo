package cli

import (
	"github.com/replicate/yolo/pkg/serve"
	"github.com/spf13/cobra"
)

var tailscale string
var cog string
var port string

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "serve",
		Short:  "serve an existing image",
		Hidden: false,
		RunE:   serveCommmand,
		Args:   cobra.ExactArgs(0),
	}

	cmd.Flags().StringVarP(&tailscale, "tailscale", "t", "", "tailscale hostname")
	cmd.Flags().StringVarP(&cog, "cog", "c", "5000", "cog listen port")
	cmd.Flags().StringVarP(&port, "port", "p", "8080", "proxy listen port")

	return cmd
}

func serveCommmand(cmd *cobra.Command, args []string) error {
	return serve.Serve(port, cog, tailscale)
}
