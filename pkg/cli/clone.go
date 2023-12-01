package cli

import (
	"fmt"
	"os"

	"github.com/replicate/yolo/pkg/images"
	"github.com/spf13/cobra"
)

func newCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "clone",
		Short:  "clone an existing image",
		Hidden: false,
		RunE:   cloneCommmand,
		Args:   cobra.ExactArgs(0),
	}

	cmd.Flags().StringVarP(&sToken, "token", "t", "", "replicate api token")
	cmd.Flags().StringVarP(&baseRef, "base", "b", "", "base image reference.  examples: owner/model or r8.im/owner/model@sha256:hexdigest")
	cmd.Flags().StringVarP(&dest, "dest", "d", "", "destination image. examples: owner/model or r8.im/owner/model")
	cmd.MarkFlagRequired("base")
	cmd.MarkFlagRequired("dest")

	return cmd
}

func cloneCommmand(cmd *cobra.Command, args []string) error {
	session := authenticate()
	if session == nil {
		fmt.Fprintln(os.Stderr, "authentication error, invalid token or registry host error")
		return nil
	}

	baseRef = images.EnsureRegistry(baseRef)
	dest = images.EnsureRegistry(dest)

	image_id, err := images.Clone(baseRef, dest, session)
	if err != nil {
		return err
	}
	fmt.Println(image_id)

	return nil
}
