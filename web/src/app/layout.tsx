import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Bawo - Learn African Languages",
  description: "Master Nigerian languages like Yoruba, Igbo, and Hausa with interactive lessons, streaks, and personalized learning paths.",
  keywords: ["language learning", "African languages", "Yoruba", "Igbo", "Hausa", "Nigerian languages"],
  openGraph: {
    title: "Bawo - Learn African Languages",
    description: "Master Nigerian languages with interactive lessons",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
