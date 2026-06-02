# Brisk.js

Hyper-lightweight JavaScript/TypeScript runtime powered by [goja](https://github.com/dop251/goja).

---

## Installation

Brisk is cross-compiled and globally distributed. Install it on macOS or Linux with a single command via our official Homebrew Tap:

```bash
brew tap matyik/brisk
brew install matyik/brisk/brisk
```

**For Windows (PowerShell):**

```powershell
iwr "https://raw.githubusercontent.com/matyik/brisk/main/install.ps1" -useb | iex
```

To verify the installation:

```bash
brisk --version
```

---

## Quick Start

Create your `server.ts` file:

```typescript
console.log('Booting Brisk Edge API...');

const PORT = 3000;

Brisk.serve(PORT, (req) => {
  console.info(`[${req.method}] ${req.url}`);

  // Secure the endpoint via header propagation
  if (req.headers['Authorization'] !== 'Bearer top-secret') {
    return { status: 401, body: 'Unauthorized' };
  }

  // Parse queries natively using the web standard
  const searchParams = new URLSearchParams(req.query);
  const userId = searchParams.get('id') || 'guest';

  return {
    status: 200,
    headers: {
      'Content-Type': 'application/json',
      'X-Powered-By': 'Brisk Engine',
    },
    body: JSON.stringify({
      message: 'Hello from the Edge!',
      user: userId,
      traceId: crypto.randomUUID(),
    }),
  };
});
```

Run your TypeScript file instantly without any pre-compilation steps:

```bash
brisk server.ts
```

---

## TypeScript & Editor Setup

Brisk implements a tailored subset of modern Web and Node interfaces. To prevent VS Code from pulling in conflicting browser definitions (`lib.dom.d.ts`), configure your project as a dedicated serverless environment.

### 1. Structure Your Types

Create a folder named `types` in your project root and place the `brisk.d.ts` declaration file there.

### 2. Configure `tsconfig.json`

Explicitly restrict the environment to modern ECMAScript specifications and target your custom declarations:

```json
{
  "compilerOptions": {
    "target": "ESNext",
    "module": "ESNext",
    "moduleResolution": "node",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,

    /* Enforce Brisk types and strip global Browser DOM types */
    "lib": ["ESNext"],
    "typeRoots": ["./types"]
  },
  "include": ["**/*.ts", "types/**/*.d.ts"],
  "exclude": ["node_modules"]
}
```

---

## Serverless & Edge Deployments

Brisk's tiny memory footprint and rapid startup make it a perfect candidate for modern serverless layers and micro-containers.

### AWS Lambda Layer Deployment

Because Brisk compiles into a single, self-contained statically linked binary with absolute zero external system dependencies, you can package it seamlessly into a Lambda Custom Runtime (`provided.al2023`).

1. Zip the compiled `brisk` binary into a file named `bootstrap`.
2. Upload it as an AWS Lambda Layer or deploy it directly as the execution runtime.
3. Handle invocations directly using your `Brisk.serve` listener mapped to port `8080`.

### High-Density Edge Containers

To deploy Brisk across decentralized container platforms (like Fly.io or AWS Fargate Edge), use a minimalist scratch-based `Dockerfile` like the one provided.

---

## Supported APIs

### `Brisk`

- `Brisk.serve(port: number, callback: (req: BriskRequest) => BriskResponse | string | void): void`

### `fetch` (WinterCG Subset)

- `fetch(url: string, options?: BriskFetchOptions): Promise<BriskFetchResponse>`

### `crypto` (WinterCG Subset)

- `crypto.randomUUID(): string`
- `crypto.getRandomValues(array: TypedArray): TypedArray` _(Max 65,536 bytes parity restriction)_

### `fs` (Sandboxed)

- `fs.readFileSync(path: string): string`
- `fs.writeFileSync(path: string, data: string): void`

### `process`

- `process.env`
- `process.argv`
- `process.cwd()`
- `process.exit(code?)`
