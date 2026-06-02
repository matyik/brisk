// the folllowing imports may throw errors in editor, but they will work when run with Brisk
// Brisk is designed to support these imports for ease of migration but they are not needed or recommended
import fs from 'node:fs';
import { randomUUID } from 'crypto';

console.info('Running legacy Node code on Brisk!');

const id = randomUUID();
console.log(`Generated Trace ID: ${id}`);

fs.writeFileSync('node-test.txt', `Trace: ${id}`);
const output = fs.readFileSync('node-test.txt');

console.log(`Successfully read from file: ${output}`);
