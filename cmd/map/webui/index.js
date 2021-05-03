import "carbon-components-svelte/css/all.css";

import App from "./src/App.svelte";

let app = new App({
  target: document.body,
});

export default app;
