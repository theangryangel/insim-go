import { readable, writable, derived } from "svelte/store";

export const state = writable({
  Event: {
    Track: {},
  },
  Players: {},
  Connections: {},
});

export const messages = writable([]);
