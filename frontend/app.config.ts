import { defineConfig } from "@solidjs/start/config";

const frontendUrl = process.env.FRONTEND_URL ?? "http://localhost:3000";
const hmrClientPort = Number(process.env.VITE_HMR_CLIENT_PORT ?? (new URL(frontendUrl).port || 3000));

export default defineConfig({
  ssr: false,
  server: {
    preset: "static",
    prerender: {
      crawlLinks: true,
    },
  },
  vite: {
    server: {
      host: "0.0.0.0",
      allowedHosts: true,
      origin: frontendUrl,
      hmr: {
        protocol: "ws",
        // host: hmrHost, // Removed to avoid EADDRNOTAVAIL in Docker. Client will use window.location.hostname.
        clientPort: hmrClientPort,
      },
    },
    plugins: [],
  },
});
