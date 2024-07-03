/** @type {import('next').NextConfig} */
const nextConfig = {}

module.exports = {
    reactStrictMode: false,
    webpack: (config) => {
        config.resolve.fallback = { fs: false };
        return config;
    },
};

