import { readable, writable, derived } from "svelte/store";

export const state = writable({
  Track: {},
  Players: {},
  Connections: {},
  Laps: 0,
});

export const messages = writable([]);
