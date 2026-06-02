package transpiler

// nodePolyfills maps Node.js package names to our internal Brisk wrappers
var nodePolyfills = map[string]string{
	"fs": `
		export const readFileSync = globalThis.fs.readFileSync;
		export const writeFileSync = globalThis.fs.writeFileSync;
		export default { readFileSync, writeFileSync };
	`,
	"crypto": `
		export const randomUUID = globalThis.crypto.randomUUID;
		export const getRandomValues = globalThis.crypto.getRandomValues;
		export default { randomUUID, getRandomValues };
	`,
	"process": `
		export default globalThis.process;
	`,
}