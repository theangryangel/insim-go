// svelte.config.js
const {
  optimizeCarbonImports,
} = require("carbon-components-svelte/preprocess");

module.exports = {
  preprocess: [optimizeCarbonImports()],
};
