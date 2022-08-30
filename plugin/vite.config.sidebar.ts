/// <reference types="vitest" />
/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />

import { defineConfig } from "vite";

import { plugin } from "./vite.config.template";

// https://vitejs.dev/config/
export default defineConfig(plugin("sidebar"));
