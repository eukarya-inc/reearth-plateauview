/// <reference types="vitest" />
/// <reference types="vite/client" />

import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import importToCDN, { autoComplete } from "vite-plugin-cdn-import";
import dts from "vite-plugin-dts";

export default defineConfig({
  plugins: [
    react(),
    importToCDN({
      modules: [autoComplete("antd")],
    }),
    dts({
      insertTypesEntry: true,
    }),
  ],
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: "./web/test/setup.ts",
  },
});
