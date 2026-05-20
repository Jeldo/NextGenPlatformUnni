import type { Config } from "tailwindcss";
import { heroui } from "@heroui/react";

const config: Config = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx,mdx}",
    "./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#F66336",
        },
        black: "#131517",
        gray: {
          description: "#697683",
        },
      },
    },
  },
  darkMode: "class",
  plugins: [heroui()],
};

export default config;
