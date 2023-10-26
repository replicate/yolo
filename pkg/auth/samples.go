package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func MakeSamples(image_id string,
	sampleDir string,
	token string,
	baseApi string) error {

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
		err = RunSample(samplePath, image_id, token, baseApi)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunSample(samplePath string,
	image_id string,
	token string,
	baseApi string,
) error {

	fileData, err := os.ReadFile(samplePath)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/submit?r8=%s&filename=%s", baseApi, image_id, samplePath)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(fileData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	return nil
}
