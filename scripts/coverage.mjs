#!/usr/bin/env node
// Run the Go tests with a coverprofile, then convert it to coverage/lcov.info
// for the fleet patch-coverage gate.

import { spawnSync } from 'node:child_process';
import { mkdirSync } from 'node:fs';
import { convert } from './coverprofile-to-lcov.mjs';

mkdirSync('coverage', { recursive: true });

const result = spawnSync('go', ['test', './...', '-coverprofile=coverage/cover.out', '-coverpkg=./...'], {
  stdio: 'inherit',
});
if (result.error) {
  process.stderr.write(`coverage: cannot run go: ${result.error.message}\n`);
  process.exit(1);
}
const status = result.status ?? 1;
if (status !== 0) process.exit(status);

convert();
