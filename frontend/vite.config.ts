import { defineConfig } from "vite";
import { devtools } from "@tanstack/devtools-vite";
import viteReact from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

import { tanstackRouter } from "@tanstack/router-plugin/vite";
import { fileURLToPath, URL } from "node:url";
import { resolve } from "node:path";
import { cp, mkdir, rm } from "node:fs/promises";

const frontendDistDir = "dist";
const backendStaticDir = fileURLToPath(
  new URL("../backend/src/static", import.meta.url),
);

const copyBackendStatic = () => ({
  name: "copy-backend-static",
  async closeBundle() {
    const sourceDir = resolve(process.cwd(), frontendDistDir);

    await rm(backendStaticDir, { recursive: true, force: true });
    await mkdir(backendStaticDir, { recursive: true });
    await cp(sourceDir, backendStaticDir, { recursive: true });
  },
});

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    devtools(),
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
    viteReact(),
    tailwindcss(),
    copyBackendStatic(),
  ],
  build: {
    outDir: frontendDistDir,
  },
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
});
