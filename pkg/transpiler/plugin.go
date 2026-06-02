package transpiler

import (
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// NodePolyfillPlugin intercepts Node imports and injects Brisk's native globals
func NodePolyfillPlugin() api.Plugin {
	return api.Plugin{
		Name: "brisk-node-polyfills",
		Setup: func(build api.PluginBuild) {
			
			build.OnResolve(api.OnResolveOptions{Filter: `^(node:)?(fs|crypto|process)$`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					return api.OnResolveResult{
						Path:      args.Path,
						Namespace: "brisk-polyfill", // Tag it so our loader catches it
					}, nil
				})

			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "brisk-polyfill"},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					moduleName := strings.TrimPrefix(args.Path, "node:")
					
					contents := nodePolyfills[moduleName]
					
					return api.OnLoadResult{
						Contents: &contents,
						Loader:   api.LoaderJS,
					}, nil
				})
		},
	}
}