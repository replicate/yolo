package images

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func MakeTar(args []string, relative bool, layers []v1.Layer) (*bytes.Buffer, error) {
	added := make(map[string]struct{})

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

		var dest string
		if relative {
			dest = filepath.Join("src", file)
		} else {
			baseName := filepath.Base(file)
			dest = filepath.Join("src", baseName)
		}

		fmt.Fprintln(os.Stderr, "adding:", dest)

		hdr := &tar.Header{
			Name: dest,
			Mode: int64(stat.Mode()),
			Size: stat.Size(),
		}
		added[hdr.Name] = struct{}{}

		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := io.Copy(tw, f); err != nil {
			return nil, err
		}
	}

	// we need to add all the layers in reverse order, so that files from
	// the most recent layer is included (not skipped)
	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]
		rc, err := layer.Uncompressed()
		if err != nil {
			return nil, err
		}

		tr := tar.NewReader(rc)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			if _, ok := added[header.Name]; ok {
				fmt.Fprintln(os.Stderr, "skipping:", header.Name)
				continue
			}

			fmt.Fprintln(os.Stderr, "including prior:", header.Name)

			if err := tw.WriteHeader(header); err != nil {
				return nil, err
			}
			if _, err := io.Copy(tw, tr); err != nil {
				return nil, err
			}

			added[header.Name] = struct{}{}
		}
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
