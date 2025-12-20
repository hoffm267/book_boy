#!/usr/bin/env bun
import { watch } from 'fs';
import { readFileSync } from 'fs';
import { resolve, extname } from 'path';

const PORT = process.env.PORT || 5173;
const clients = new Set();

// Watch for file changes
const watcher = watch('./src', { recursive: true }, () => {
  // Notify all connected clients to reload
  clients.forEach(client => {
    try {
      client.send('reload');
    } catch (e) {
      clients.delete(client);
    }
  });
});

// Serve the application
const server = Bun.serve({
  port: PORT,
  async fetch(req, server) {
    const url = new URL(req.url);

    // WebSocket for hot reload
    if (url.pathname === '/__hot') {
      if (server.upgrade(req)) {
        return;
      }
      return new Response('Upgrade failed', { status: 500 });
    }

    // Handle file requests
    let filePath = url.pathname === '/' ? '/index.html' : url.pathname;

    try {
      // Try to read the file
      const fullPath = resolve('.', filePath.slice(1));
      const file = Bun.file(fullPath);

      if (await file.exists()) {
        let content = await file.text();

        // Inject hot reload script into HTML
        if (filePath.endsWith('.html')) {
          content = content.replace(
            '</body>',
            `  <script>
    const ws = new WebSocket('ws://localhost:${PORT}/__hot');
    ws.onmessage = () => location.reload();
  </script>
</body>`
          );
          return new Response(content, {
            headers: { 'Content-Type': 'text/html' }
          });
        }

        // Handle JSX/TSX files
        if (filePath.endsWith('.jsx') || filePath.endsWith('.tsx')) {
          const transpiled = await Bun.build({
            entrypoints: [fullPath],
            target: 'browser',
          });

          if (transpiled.success && transpiled.outputs[0]) {
            return new Response(await transpiled.outputs[0].text(), {
              headers: { 'Content-Type': 'application/javascript' }
            });
          }
        }

        return new Response(file);
      }

      return new Response('Not found', { status: 404 });
    } catch (error) {
      console.error('Error serving file:', error);
      return new Response('Error: ' + error.message, { status: 500 });
    }
  },
  websocket: {
    open(ws) {
      clients.add(ws);
    },
    message(ws, message) {},
    close(ws) {
      clients.delete(ws);
    }
  }
});

console.log(`Dev server running at http://localhost:${PORT}`);
console.log('Watching for changes...');
