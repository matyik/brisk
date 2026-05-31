package transpiler

import (
	"fmt"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func Process(filePath, code string) (string, error) {
	if !strings.HasSuffix(filePath, ".ts") {
		return code, nil
	}

	result := api.Transform(code, api.TransformOptions{
		Loader: api.LoaderTS,
	})

	if len(result.Errors) > 0 {
		var errMsgs []string
		for _, err := range result.Errors {
			errMsgs = append(errMsgs, err.Text)
		}
		return "", fmt.Errorf("TypeScript compilation failed:\n%s", strings.Join(errMsgs, "\n"))
	}

	return string(result.Code), nil
}