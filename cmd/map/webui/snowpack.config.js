/** @type {import("snowpack").SnowpackUserConfig } */
module.exports = {
  mount: {
    /* ... */
  },
  plugins: [
    "@snowpack/plugin-svelte",
    "@snowpack/plugin-postcss",
    // temporary workaround for https://github.com/snowpackjs/snowpack/issues/2916
    "@jadex/snowpack-plugin-tailwindcss-jit",
  ],
  routes: [
    /* Enable an SPA Fallback in development: */
    // {"match": "routes", "src": ".*", "dest": "/index.html"},
  ],
  optimize: {
    /* Example: Bundle your final build: */
    // "bundle": true,
  },
  packageOptions: {
    /* ... */
  },
  devOptions: {
    /* ... */
  },
  buildOptions: {
    clean: true,
    out: "../static",
    /* ... */
  },
};
