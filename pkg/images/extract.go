package images

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func Extract(baseRef string, dest string, session authn.Authenticator) error {
	var err error

	if _, err = os.Stat(dest); !os.IsNotExist(err) {
		return fmt.Errorf("destination %s already exists", dest)
	}

	var base v1.Image

	fmt.Fprintln(os.Stderr, "fetching metadata for", baseRef)

	base, err = crane.Pull(baseRef, crane.WithAuth(session))
	if err != nil {
		return fmt.Errorf("pulling %w", err)
	}

	src, err := GetSourceLayers(base, true, true)
	if err != nil {
		return err
	}

	for _, layer := range src {
		rc, err := layer.Uncompressed()
		if err != nil {
			return err
		}

		tr := tar.NewReader(rc)
		err = extractTarFile(tr, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractTarFile(tarReader *tar.Reader, destDir string) error {
	startTime := time.Now()
	var _fileSize int64

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)
		targetDir := filepath.Dir(target)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			fmt.Fprintln(os.Stderr, target)
			targetFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(targetFile, tarReader); err != nil {
				targetFile.Close()
				return err
			}
			_fileSize += header.Size
			targetFile.Close()
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported file type for %s, typeflag %s", header.Name, string(header.Typeflag))
		}
	}

	elapsed := time.Since(startTime).Seconds()
	size := humanize.Bytes(uint64(_fileSize))
	throughput := humanize.Bytes(uint64(float64(_fileSize) / elapsed))
	fmt.Fprintf(os.Stderr, "Extracted %s in %.3fs (%s/s)\n", size, elapsed, throughput)

	return nil
}
