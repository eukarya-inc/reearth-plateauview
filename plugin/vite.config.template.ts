/// <reference types="vitest" />
/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />

import react from "@vitejs/plugin-react";
// import { visualizer } from "rollup-plugin-visualizer";
import type { UserConfigExport, Plugin } from "vite";
import importToCDN, { autoComplete } from "vite-plugin-cdn-import";
import { viteExternalsPlugin } from "vite-plugin-externals";
import { viteSingleFile } from "vite-plugin-singlefile";
import svgr from "vite-plugin-svgr";

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
  },
});

export const web =
  (name: string): UserConfigExport =>
  ({ mode }) => ({
    plugins: [
      react(),
      svgr(),
      serverHeaders(),
      viteSingleFile(),

      mode === "production" &&
        (importToCDN /* workaround */ as any as { default: typeof importToCDN }).default({
          modules: [autoComplete("react"), autoComplete("react-dom")],
        }),
      mode === "production" &&
        viteExternalsPlugin({
          react: "React",
          "react-dom": "ReactDOM",
        }),
    ],
    publicDir: false,
    root: `./web/${name}`,
    build: {
      outDir: `../../dist/web/${name}`,
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
    test: {
      globals: true,
      environment: "jsdom",
      setupFiles: "./web/test/setup.ts",
    },
  });

function serverHeaders(): Plugin {
  return {
    name: "server-headers",
    configureServer(server) {
      server.middlewares.use((_req, res, next) => {
        res.setHeader("Service-Worker-Allowed", "/");
        next();
      });
    },
  };
}
