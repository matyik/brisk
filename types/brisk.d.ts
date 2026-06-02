/**
 * BRISK ENGINE STANDARD LIBRARY
 * Ambient Type Declarations
 * * Include this file in your tsconfig.json "typeRoots" or reference it
 * directly using /// <reference path="brisk.d.ts" />
 */

// ============================================================================
// 1. GLOBAL NAMESPACES & CLASSES
// ============================================================================

/**
 * The Brisk-specific HTTP Request object passed to `Brisk.serve` callbacks.
 */
interface BriskRequest {
  /** The HTTP Method (e.g., "GET", "POST") */
  method: string;
  /** The incoming path (e.g., "/api/users") */
  url: string;
  /** The raw stringified body of the request */
  body: string;
}

/**
 * The expected response format for `Brisk.serve`.
 */
interface BriskResponse {
  /** HTTP Status Code (Defaults to 200) */
  status?: number;
  /** The response payload */
  body?: string;
}

/**
 * The global Brisk Engine API.
 */
declare namespace Brisk {
  /**
   * Starts a high-performance, single-threaded Event Loop backed by Go's `net/http`.
   * @param port The port to bind the server to.
   * @param callback The function executed for every incoming request.
   */
  function serve(
    port: number,
    callback: (req: BriskRequest) => BriskResponse | string | void,
  ): void;
}

// ============================================================================
// 2. FETCH API (WINTER-CG COMPLIANT SUBSET)
// ============================================================================

interface BriskFetchOptions {
  /** HTTP Method (GET, POST, PUT, DELETE, etc.) */
  method?: string;
  /** Request headers as a simple key-value map */
  headers?: Record<string, string>;
  /** Request body payload */
  body?: string;
}

interface BriskFetchResponse {
  /** True if status is between 200 and 299 */
  readonly ok: boolean;
  /** The HTTP status code */
  readonly status: number;
  /** Reads the response body and returns it as a string */
  text(): Promise<string>;
  /** Reads the response body, parses it natively in Go, and returns a JS Object */
  json<T = any>(): Promise<T>;
}

/**
 * Performs an asynchronous network request natively handled by Go's net/http.
 * @param url The endpoint to fetch.
 * @param options HTTP options (method, headers, body).
 */
declare function fetch(
  url: string,
  options?: BriskFetchOptions,
): Promise<BriskFetchResponse>;

// ============================================================================
// 3. URL API (WINTER-CG COMPLIANT)
// ============================================================================

declare class URLSearchParams {
  /** Retrieves the first value associated with the given search parameter. */
  get(name: string): string | null;
  /** Sets the value associated with a search parameter. Overwrites existing values. */
  set(name: string, value: string): void;
  /** Appends a specified key/value pair as a new search parameter. */
  append(name: string, value: string): void;
  /** Deletes the given search parameter and all its associated values. */
  delete(name: string): void;
}

declare class URL {
  constructor(url: string);
  /** The full serialized URL */
  readonly href: string;
  /** The protocol scheme (e.g., "https:") */
  readonly protocol: string;
  /** The domain name */
  readonly hostname: string;
  /** The path section of the URL */
  readonly pathname: string;
  /** The serialized query parameters (e.g., "?q=brisk") */
  readonly search: string;
  /** An interface to interactively mutate the search query string */
  readonly searchParams: URLSearchParams;
  /** Returns the full serialized URL */
  toString(): string;
}

// ============================================================================
// 4. PROCESS API (NODE COMPLIANT SUBSET)
// ============================================================================

declare namespace process {
  /** * An object containing the user environment variables bridged from the host OS.
   */
  const env: Record<string, string | undefined>;

  /** * An array containing the command-line arguments passed when the Brisk process was launched.
   */
  const argv: string[];

  /** * Returns the current working directory of the Brisk process.
   */
  function cwd(): string;

  /** * Instructs Brisk to terminate the process synchronously with an exit status.
   * @param code The exit code (Defaults to 0 for success).
   */
  function exit(code?: number): never;
}

// ============================================================================
// 5. FILE SYSTEM API (BRISK SANDBOXED)
// ============================================================================

declare namespace fs {
  /**
   * Reads the entire contents of a file synchronously.
   * @security **Sandboxed**: Will throw an Error if attempting to read outside the execution root.
   * @param path The relative or absolute path to the file.
   */
  function readFileSync(path: string): string;

  /**
   * Writes data to a file synchronously, replacing the file if it already exists.
   * @security **Sandboxed**: Will throw an Error if attempting to write outside the execution root.
   * @param path The relative or absolute path to the file.
   * @param data The string payload to write.
   */
  function writeFileSync(path: string, data: string): void;
}

// ============================================================================
// 6. CRYPTO API (WINTER-CG COMPLIANT)
// ============================================================================

interface ArrayBufferView {
  /**
   * The length in bytes of the array.
   */
  byteLength: number;
}

declare namespace crypto {
  /**
   * Generates a cryptographically secure Version 4 UUID using Go's native `crypto/rand`.
   */
  function randomUUID(): string;

  /**
   * Fills the provided TypedArray with cryptographically secure random values.
   * @param array The array to fill (e.g., Uint8Array). Cannot exceed 65,536 bytes.
   * @returns The exact same array instance passed in, now filled with entropy.
   */
  function getRandomValues<T extends ArrayBufferView>(array: T): T;
}

// ============================================================================
// 7. CONSOLE API (POSIX STREAM AWARE)
// ============================================================================

declare namespace console {
  /** Prints to standard output (stdout). */
  function log(...data: any[]): void;
  /** Prints to standard output (stdout). */
  function info(...data: any[]): void;
  /** Prints a yellow warning to standard error (stderr). */
  function warn(...data: any[]): void;
  /** Prints a red error to standard error (stderr). */
  function error(...data: any[]): void;
}
