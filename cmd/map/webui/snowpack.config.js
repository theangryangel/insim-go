const production = process.env["NODE_ENV"] == "production";

/** @type {import("snowpack").SnowpackUserConfig } */
module.exports = {
  mount: {
    /* ... */
  },
  plugins: [
    "@snowpack/plugin-svelte",
    "@snowpack/plugin-postcss",
    //"@snowpack/plugin-sass",
    // temporary workaround for https://github.com/snowpackjs/snowpack/issues/2916
    //"@jadex/snowpack-plugin-tailwindcss-jit",
  ],
  routes: [
    /* Enable an SPA Fallback in development: */
    // {"match": "routes", "src": ".*", "dest": "/index.html"},
  ],
  optimize: {
    bundle: production,
    minify: production,
    target: "es2018",
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
