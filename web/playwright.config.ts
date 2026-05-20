import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  baseURL: "http://localhost:3000",
  use: {
    headless: true,
    viewport: { width: 390, height: 844 },
  },
  webServer: {
    command: "pnpm dev",
    port: 3000,
    reuseExistingServer: true,
  },
});
