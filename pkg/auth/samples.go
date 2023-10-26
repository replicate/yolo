package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/authn"
)

func MakeSamples(image_id string,
	sampleDir string,
	session authn.Authenticator, baseApi string) error {

	// for each json in sampleDir, run the sample
	sampleFiles, err := os.ReadDir(sampleDir)
	if err != nil {
		return err
	}
	for _, sampleFile := range sampleFiles {
		if sampleFile.IsDir() {
			continue
		}
		if filepath.Ext(sampleFile.Name()) != ".json" {
			continue
		}
		samplePath := filepath.Join(sampleDir, sampleFile.Name())
		fmt.Println("running sample:", samplePath)
		err = RunSample(samplePath, image_id, session, baseApi)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunSample(samplePath string,
	image_id string,
	auth authn.Authenticator,
	baseApi string,
) error {

	fileData, err := os.ReadFile(samplePath)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/submit?r8=%s&filename=%s", baseApi, image_id, samplePath)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(fileData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	return nil

}
