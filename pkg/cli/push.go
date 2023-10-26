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
	sBaseApi  string
	baseRef   string
	dest      string
	ast       string
	commit    string
	sampleDir string
)

func newPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "push",
		Short:  "update an existing image",
		Hidden: false,
		RunE:   pushCommmand,
		Args:   cobra.MinimumNArgs(1),
	}

	cmd.Flags().StringVarP(&sToken, "token", "t", "", "replicate api token")
	cmd.Flags().StringVarP(&sRegistry, "registry", "r", "r8.im", "registry host")
	cmd.Flags().StringVarP(&baseRef, "base", "b", "", "base image reference.  examples: owner/model or r8.im/owner/model@sha256:hexdigest")
	cmd.MarkFlagRequired("base")
	cmd.Flags().StringVarP(&dest, "dest", "d", "", "destination image. examples: owner/model or r8.im/owner/model")
	cmd.MarkFlagRequired("dest")
	cmd.Flags().StringVarP(&ast, "ast", "a", "", "optional file to parse AST to update openapi schema")
	cmd.Flags().StringVarP(&commit, "commit", "c", "", "optional commit hash to update git commit")
	cmd.Flags().StringVarP(&sampleDir, "sample-dir", "s", "", "optional directory to run samples")
	cmd.Flags().StringVarP(&sBaseApi, "test-api", "u", "http://localhost:4000", "experiment endpoint")

	return cmd
}

func pushCommmand(cmd *cobra.Command, args []string) error {
	if sToken == "" {
		sToken = os.Getenv("REPLICATE_API_TOKEN")
	}

	if sToken == "" {
		sToken = os.Getenv("COG_TOKEN")
	}

	u, err := auth.VerifyCogToken(sRegistry, sToken)
	if err != nil {
		fmt.Fprintln(os.Stderr, "authentication error, invalid token or registry host error")
		return err
	}
	session := authn.FromConfig(authn.AuthConfig{Username: u, Password: sToken})

	tar, err := makeTar(args)
	if err != nil {
		return err
	}

	baseRef = ensureRegistry(baseRef)
	dest = ensureRegistry(dest)

	image_id, err := images.Affix(baseRef, dest, tar, ast, commit, session)
	if err != nil {
		return err
	}

	fmt.Println(image_id)

	if sampleDir != "" {
		fmt.Println("running samples")
		err = auth.MakeSamples(image_id, sampleDir, sToken, sBaseApi)
		if err != nil {
			return err
		}
	}

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
		dest := filepath.Join("src", baseName)
		fmt.Fprintln(os.Stderr, "adding:", dest)

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
