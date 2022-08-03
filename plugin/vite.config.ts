/// <reference types="vitest" />
/// <reference types="vite/client" />

import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import importToCDN, { autoComplete } from "vite-plugin-cdn-import";

export default defineConfig({
  plugins: [
    react(),
    importToCDN({
      modules: [autoComplete("antd")],
    }),
  ],
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: "./web/test/setup.ts",
  },
});
