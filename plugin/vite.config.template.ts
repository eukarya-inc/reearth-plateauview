/// <reference types="vitest" />
/// <reference types="vite/client" />

import react from "@vitejs/plugin-react";
import type { UserConfigExport, Plugin } from "vite";
import { viteSingleFile } from "vite-plugin-singlefile";

export const plugin = (name: string): UserConfigExport => ({
  build: {
    outDir: "dist/plugin",
    emptyOutDir: false,
    lib: {
      formats: ["iife"],
      // https://github.com/vitejs/vite/pull/7047
      entry: `src/${name}.ts`,
      name: `ReearthPluginPV_${name}`,
      fileName: () => `${name}.js`,
    },
    rollupOptions: {
      external: ["react", "react-dom", "antd"],
      output: {
        globals: {
          react: "React",
          "react-dom": "ReactDOM",
          antd: "antd",
        },
      },
    },
  },
});

export const web =
  (name: string): UserConfigExport =>
  () => ({
    plugins: [react(), viteSingleFile(), serverHeaders()],
    publicDir: false,
    globals: true,
    root: `./web/${name}`,
    build: {
      outDir: `../../dist/web/${name}`,
    },
    test: {
      globals: true,
      environment: "jsdom",
      setupFiles: "./web/test/setup.ts",
    },
    css: {
      preprocessorOptions: {
        less: {
          javascriptEnabled: true,
          modifyVars: {
            "@primary-color": "#00BEBE",
          },
        },
      },
    },
  });

const serverHeaders = (): Plugin => ({
  name: "server-headers",
  configureServer(server) {
    server.middlewares.use((_req, res, next) => {
      res.setHeader("Service-Worker-Allowed", "/");
      next();
    });
  },
});
