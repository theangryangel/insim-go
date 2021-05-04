import { players, event, messages } from "./stores";

export class SSE {
  constructor() {
    this.es = null;
  }

  dial() {
    this.es = new EventSource(`/events`);
    this.es.onerror = () => {
      console.log(`EventSource error`);
      console.log("Reconnecting in 1s", true);
      setTimeout(() => {
        this.dial();
      }, 1000);
    };

    this.es.onopen = () => {
      console.info("eventsource connected, fetching initial state");

      fetch("/api/state")
        .then((res) => {
          return res.json();
        })
        .then((data) => {
          players.set(data.Players);
          event.set(data.Event);
        });
    };

    this.es.addEventListener("chat", (ev) => {
      const data = JSON.parse(ev.data);
      messages.update((n) => {
        n.unshift(data);
        return n.slice(0, 10);
      });
    });

    this.es.addEventListener("state", (ev) => {
      let data = JSON.parse(ev.data);
      players.set(data.Players);
      event.set(data.Event);
    });

    this.es.addEventListener("player-state", (ev) => {
      let data = JSON.parse(ev.data);

      players.update((n) => {
        return { ...n, [data.Plid]: data.State };
      });
    });

    this.es.addEventListener("player-left", (ev) => {
      let data = JSON.parse(ev.data);

      let p = players;
      delete p[data.Plid];
      players.set(p);
    });
  }
}
