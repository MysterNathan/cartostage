import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',
  basePath: '/cartostage',
  assetPrefix: '/cartostage',
  images: {
    domains: [
      'tile.openstreetmap.org',
      'a.tile.openstreetmap.org',
      'b.tile.openstreetmap.org',
      'c.tile.openstreetmap.org'
    ],
  },
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL
  },
  eslint: {
    // Ignore ESLint errors during build
    ignoreDuringBuilds: true,
  },
  typescript: {
    // Ignore TypeScript errors during build
    ignoreBuildErrors: true,
  },};

export default nextConfig;
