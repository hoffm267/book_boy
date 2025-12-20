#!/usr/bin/env bun
import { readFileSync, writeFileSync, readdirSync } from 'fs';
import { join } from 'path';

// Build the application
const result = await Bun.build({
  entrypoints: ['./src/main.jsx'],
  outdir: './dist',
  target: 'browser',
  minify: true,
  splitting: true,
  naming: {
    entry: '[name]-[hash].[ext]',
    chunk: '[name]-[hash].[ext]',
    asset: '[name]-[hash].[ext]'
  }
});

if (!result.success) {
  console.error('Build failed:', result.logs);
  process.exit(1);
}

// Find the generated files
const distFiles = readdirSync('./dist');
const jsFile = distFiles.find(f => f.startsWith('main-') && f.endsWith('.js'));
const cssFile = distFiles.find(f => f.startsWith('main-') && f.endsWith('.css'));

// Read the template HTML
const htmlTemplate = readFileSync('./index.html', 'utf-8');

// Generate the production HTML
const productionHtml = htmlTemplate
  .replace(
    '<script type="module" src="/src/main.jsx"></script>',
    `<link rel="stylesheet" href="/${cssFile}">\n  <script type="module" src="/${jsFile}"></script>`
  );

// Write the production HTML
writeFileSync('./dist/index.html', productionHtml);

console.log('Build complete!');
console.log('Generated files:');
console.log(`  - ${jsFile}`);
if (cssFile) console.log(`  - ${cssFile}`);
console.log('  - index.html');
