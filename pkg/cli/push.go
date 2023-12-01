package cli

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/replicate/yolo/pkg/auth"
	"github.com/replicate/yolo/pkg/images"
	"github.com/spf13/cobra"
)

var (
	sToken        string
	sBaseApi      string
	baseRef       string
	dest          string
	ast           string
	commit        string
	sampleDir     string
	relativePaths bool
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
	cmd.Flags().BoolVarP(&relativePaths, "relative-paths", "p", false, "preserve relative paths from where yolo is run instead of placing all files under /src")
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
	session := authenticate()
	if session == nil {
		fmt.Fprintln(os.Stderr, "authentication error, invalid token or registry host error")
		return nil
	}

	baseRef = images.EnsureRegistry(baseRef)
	dest = images.EnsureRegistry(dest)

	base, err := crane.Pull(baseRef, crane.WithAuth(session))
	if err != nil {
		return fmt.Errorf("pulling %w", err)
	}

	// FIXME(ja): I think there should be a method images.UpdateYolo
	// that takes the new files and updates the yolo if it exists
	// -- Affix is too low level for this?

	yoloLayers, err := images.GetSourceLayers(base, false, true)
	if err != nil {
		return err
	}

	var files []images.LayerFile
	for _, path := range args {
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var dest string
		if relativePaths {
			dest = filepath.Join("src", path)
		} else {
			baseName := filepath.Base(path)
			dest = filepath.Join("src", baseName)
		}
		file := images.LayerFile{
			Header: &tar.Header{
				Name: dest,
				Mode: 0644,
				Size: int64(len(body)),
			},
			Body: body,
		}

		files = append(files, file)
	}

	tar, err := images.MakeTar(files, yoloLayers)
	if err != nil {
		return err
	}

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
