import type { Metadata } from "next";
import { HeroUIProvider } from "@heroui/react";
import { QueryProvider } from "@/lib/QueryProvider";
import "./globals.css";

export const metadata: Metadata = {
  title: "시술 캘린더",
  description: "시술 관리 캘린더",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ko">
      <body className="text-black">
        <QueryProvider>
          <HeroUIProvider>{children}</HeroUIProvider>
        </QueryProvider>
      </body>
    </html>
  );
}
