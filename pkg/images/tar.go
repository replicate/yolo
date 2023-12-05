package images

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type LayerFile struct {
	Header *tar.Header
	Body   []byte
}

func MakeTar(files []LayerFile, layers []v1.Layer) (*bytes.Buffer, error) {
	added := make(map[string]struct{})

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	for _, file := range files {
		fmt.Fprintln(os.Stderr, "adding:", file.Header.Name)

		if err := tw.WriteHeader(file.Header); err != nil {
			return nil, err
		}
		if _, err := tw.Write(file.Body); err != nil {
			return nil, err
		}

		added[file.Header.Name] = struct{}{}
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

		rc.Close()
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
