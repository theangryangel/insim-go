import { readable, writable, derived } from "svelte/store";

export const event = writable({
  Track: {},
});
export const players = writable({});
export const messages = writable([]);
