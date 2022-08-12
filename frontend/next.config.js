const backendURL = process.env.SPOUT_BACKEND_URL || "http://127.0.0.1:3280/";

module.exports = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${backendURL}api/:path*`,
      },
    ];
  },
};
