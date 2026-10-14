import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Output standalone is used in Docker container builds
  output: process.env.NEXT_OUTPUT_STANDALONE === "true" ? "standalone" : undefined,
  reactStrictMode: true,
};

export default nextConfig;
