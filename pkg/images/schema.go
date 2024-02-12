package images

import (
	"bytes"
	_ "embed"
	"fmt"
	"os/exec"
)

//go:embed ast_openapi_schema.py
var ast_openapi_schema string

func GetSchema(predictorToParse string) (string, error) {
	cmd := exec.Command("python3", "-c", ast_openapi_schema, predictorToParse)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("running ast_openapi_schema.py: %w", err)
	}

	schema := string(bytes.TrimSpace(out.Bytes()))
	return schema, nil
}
