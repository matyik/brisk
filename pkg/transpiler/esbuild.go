package transpiler

import (
	"fmt"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// Process bundles the file and its imports, stripping TS types automatically.
func Process(filePath string) (string, error) {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{filePath},
		Bundle:      true,
		Write:       false,
		Platform:    api.PlatformNeutral,
		Format:      api.FormatIIFE,
	})

	if len(result.Errors) > 0 {
		var errMsgs []string
		for _, err := range result.Errors {
			errMsgs = append(errMsgs, err.Text)
		}
		return "", fmt.Errorf("Bundling failed:\n%s", strings.Join(errMsgs, "\n"))
	}

	return string(result.OutputFiles[0].Contents), nil
}