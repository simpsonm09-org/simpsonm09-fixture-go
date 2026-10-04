#!/usr/bin/env node
// Convert a Go coverprofile to lcov for the fleet patch-coverage gate. Go
// writes a coverprofile, not lcov, so this repository-owned converter reads
// coverage/cover.out and writes coverage/lcov.info.
//
// usage: node scripts/coverprofile-to-lcov.mjs

import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

const PROFILE = 'coverage/cover.out';
const OUTPUT = 'coverage/lcov.info';

// The coverprofile records each file by import path. Strip the module path from
// go.mod so lcov carries the repository-relative path the diff uses.
function modulePrefix() {
  if (!existsSync('go.mod')) return '';
  const match = /^module\s+(\S+)/m.exec(readFileSync('go.mod', 'utf8'));
  return match ? `${match[1]}/` : '';
}

// parseCoverprofile returns a map of file to a map of line number to the most
// hits any covering block recorded for that line.
export function parseCoverprofile(text, prefix = '') {
  const files = new Map();
  for (const rawLine of text.split('\n')) {
    const line = rawLine.trim();
    if (!line || line.startsWith('mode:')) continue;
    const match = /^(.+):(\d+)\.\d+,(\d+)\.\d+\s+\d+\s+(\d+)$/.exec(line);
    if (!match) continue;
    let file = match[1];
    if (prefix && file.startsWith(prefix)) file = file.slice(prefix.length);
    const start = Number(match[2]);
    const end = Number(match[3]);
    const hits = Number(match[4]);
    if (!files.has(file)) files.set(file, new Map());
    const lines = files.get(file);
    for (let number = start; number <= end; number += 1) {
      lines.set(number, Math.max(lines.get(number) ?? 0, hits));
    }
  }
  return files;
}

// toLcov renders the parsed coverage as an lcov report.
export function toLcov(files) {
  const out = [];
  for (const [file, lines] of files) {
    out.push(`SF:${file}`);
    const numbers = [...lines.keys()].sort((a, b) => a - b);
    for (const number of numbers) out.push(`DA:${number},${lines.get(number)}`);
    out.push('end_of_record');
  }
  return `${out.join('\n')}\n`;
}

// convert reads the coverprofile and writes coverage/lcov.info.
export function convert() {
  if (!existsSync(PROFILE)) {
    process.stderr.write(`coverprofile-to-lcov: no coverprofile at ${PROFILE}\n`);
    process.exit(1);
  }
  const files = parseCoverprofile(readFileSync(PROFILE, 'utf8'), modulePrefix());
  mkdirSync(dirname(resolve(OUTPUT)), { recursive: true });
  writeFileSync(OUTPUT, toLcov(files), 'utf8');
  process.stdout.write(`coverprofile-to-lcov: wrote ${OUTPUT} (${files.size} files)\n`);
}

if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  convert();
}
