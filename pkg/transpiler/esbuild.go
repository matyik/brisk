package transpiler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// Process bundles the file and its imports, stripping TS types automatically.
func Process(filePath string) (string, error) {
	// We use Build instead of Transform so it crawls your imports
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{filePath},
		Bundle:      true,                // THIS IS THE MAGIC BULLET
		Write:       false,               // Keep the output in memory, don't write a new file to disk
		Platform:    api.PlatformNeutral, // Don't assume Node.js or Browser
		Format:      api.FormatIIFE,      // Wrap everything nicely so variables don't leak
	})

	if len(result.Errors) > 0 {
		var errMsgs []string
		for _, err := range result.Errors {
			errMsgs = append(errMsgs, err.Text)
		}
		return "", errors.New(fmt.Sprintf("Bundling failed:\n%s", strings.Join(errMsgs, "\n")))
	}

	// esbuild returns an array of output files. Since we have one entrypoint, we grab index 0.
	return string(result.OutputFiles[0].Contents), nil
}