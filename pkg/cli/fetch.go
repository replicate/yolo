package cli

import (
	"fmt"
	"os"

	"github.com/replicate/yolo/pkg/images"
	"github.com/spf13/cobra"
)

func newFetchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "fetch",
		Short:  "fetch an existing image",
		Hidden: false,
		RunE:   fetchCommmand,
		Args:   cobra.ExactArgs(1),
	}

	cmd.Flags().StringVarP(&sToken, "token", "t", "", "replicate api token")
	cmd.Flags().StringVarP(&sRegistry, "registry", "r", "r8.im", "registry host")
	cmd.Flags().StringVarP(&baseRef, "base", "b", "", "base image reference.  examples: owner/model or r8.im/owner/model@sha256:hexdigest")
	cmd.MarkFlagRequired("base")

	return cmd
}

func fetchCommmand(cmd *cobra.Command, args []string) error {
	dest := args[0]

	session := authenticate()
	if session == nil {
		fmt.Fprintln(os.Stderr, "authentication error, invalid token or registry host error")
		return nil
	}

	baseRef = ensureRegistry(baseRef)
	return images.Extract(baseRef, dest, session)
}
