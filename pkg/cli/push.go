package cli

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/replicate/yolo/pkg/auth"
	"github.com/replicate/yolo/pkg/images"
	"github.com/spf13/cobra"
)

var (
	sToken    string
	sRegistry string
	baseRef   string
	dest      string
)

func newPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "push",
		Short:  "update an existing image",
		Hidden: false,
		RunE:   pushCommmand,
		Args:   cobra.MinimumNArgs(1),
	}

	cmd.Flags().StringVarP(&sToken, "token", "t", "", "replicate cog token")
	cmd.Flags().StringVarP(&sRegistry, "registry", "r", "r8.im", "registry host")
	cmd.Flags().StringVarP(&baseRef, "base", "b", "", "base image reference - include tag: r8.im/username/modelname@sha256:hexdigest")
	cmd.MarkFlagRequired("base")
	cmd.Flags().StringVarP(&dest, "dest", "d", "", "destination image reference: r8.im/username/modelname")
	cmd.MarkFlagRequired("dest")

	return cmd
}

func pushCommmand(cmd *cobra.Command, args []string) error {
	if sToken == "" {
		sToken = os.Getenv("COG_TOKEN")
	}

	u, err := auth.VerifyCogToken(sRegistry, sToken)
	if err != nil {
		fmt.Fprintln(os.Stderr, "authentication error, invalid token or registry host error")
		return err
	}
	auth := authn.FromConfig(authn.AuthConfig{Username: u, Password: sToken})

	tar, err := makeTar(args)
	if err != nil {
		return err
	}

	image_id, err := images.Affix(baseRef, dest, tar, auth)
	if err != nil {
		return err
	}

	fmt.Println(image_id)

	return nil
}

func makeTar(args []string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	for _, file := range args {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}

		// FIXME(ja): we shouldn't always just put things in /src stripping the path
		baseName := filepath.Base(file)
		dest := filepath.Join("/src", baseName)
		fmt.Println("adding:", dest)

		hdr := &tar.Header{
			Name: dest,
			Mode: int64(stat.Mode()),
			Size: stat.Size(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := io.Copy(tw, f); err != nil {
			return nil, err
		}
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
