import { defineConfig } from "vite";
import { devtools } from "@tanstack/devtools-vite";
import viteReact from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

import { tanstackRouter } from "@tanstack/router-plugin/vite";
import { fileURLToPath, URL } from "node:url";
import { resolve } from "node:path";
import { cp, mkdir, rm } from "node:fs/promises";

const frontendDistDir = "dist";
const serverStaticDir = fileURLToPath(
  new URL("../leszmonitor-server/src/static", import.meta.url),
);

const copyServerStatic = () => ({
  name: "copy-server-static",
  async closeBundle() {
    const sourceDir = resolve(process.cwd(), frontendDistDir);

    await rm(serverStaticDir, { recursive: true, force: true });
    await mkdir(serverStaticDir, { recursive: true });
    await cp(sourceDir, serverStaticDir, { recursive: true });
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
    copyServerStatic(),
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
