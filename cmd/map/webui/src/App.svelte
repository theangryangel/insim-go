<!-- Map.svelte -->
<script>
  import { onMount } from 'svelte';

    let state = {
        Track: {},
        Players: {},
        Connections: {},
        Laps: 0,
      };

    let messages = [];

    let ws = null;

    const pushMsg = function(msg) {
        messages.unshift(msg);
        // needed for svetle to see the change. eww.
        messages = messages.slice(0, 10)
      }

    const dial = () => {
        ws = new WebSocket(`ws://${document.location.host}/subscribe`);

        ws.addEventListener("close", (ev) => {
            console.log(
                `WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`,
                true
              );
            if (ev.code !== 1001) {
                console.log("Reconnecting in 1s", true);
                setTimeout(dial, 1000);
              }
          });
        ws.addEventListener("open", (ev) => {
            console.info("websocket connected");
          });

        // This is where we handle messages received.
        ws.addEventListener("message", (ev) => {
            if (typeof ev.data !== "string") {
                console.error("unexpected message type", typeof ev.data);
                return;
              }

            let data = JSON.parse(ev.data);
            switch (data.Type) {
                case 'chat':
                  pushMsg(data.Payload);
                  break;
                case 'state':
                  console.log("new state")
                  state = data.Payload
                  break;
                case 'player-state':
                  let obj = {}
                  obj[data.Payload.Plid] = data.Payload.State
                  state.Players = Object.assign(state.Players, obj)
                  break;
                case 'player-left':
                  if (state.Players[data.Payload]) {
                      delete state.Players[data.Payload]
                      state = state
                    }
                  break;
                default:
                  break;
              }

          });
      }

    const initialState = async () => {
        return fetch("/api/state")
          .then((res) => {
              return res.json();
            })
          .then((res) => {
              state = res;
            });
      }

    onMount(async () => {
        initialState().then(() => { dial(); console.log(state) })
      });
</script>
<style>
  .Map {
    @apply flex flex-col w-full h-screen;
  }

  .Map table {
    @apply w-full text-left border-collapse;
  }

  .Map table td, .Map table th {
    @apply p-2;
  }

  .Map-header {
    @apply flex flex-grow bg-blue-200 flex-row overflow-y-auto;
  }

  .Map-left {
    @apply w-3/5 bg-gray-800 p-3 overflow-y-auto;
  }

  .Map-right {
    @apply flex-grow flex flex-col bg-white text-gray-900;
  }

  .Map-right .Map-map {
    @apply h-3/4 p-3;
  }

  .Map-right .Map-chat {
    @apply flex-grow overflow-y-auto bg-gray-200 p-3;
  }

  .Map-footer {
    @apply flex border-red-900 border-t p-3;
  }
</style>
<div class="Map">
  <div class="Map-header">
    <div class="Map-left">
      <table>
        <thead>
          <tr>
            <th></th>
            <th>#</th>
            <th>Driver<br><small>Car</small></th>
            <th>Lap</th>
            <th>Pit Stops</th>
            <th>Gap</th>
            <th>Best</th>
            <th>Total</th>
          </tr>
        </thead>
        <tbody class="">
          {#each Object.values(state.Players).filter((a) => { return a.RacePosition > 0; }).sort((a, b) => { return a.RacePosition - b.RacePosition; }) as player}
            <tr>
              <td>
                {#if player.RaceFinished}
                  🏁
                {/if}
              </td>
              <td>{player.RacePosition}
              </td>
              <td>{player.Playername}<br><small>{player.Vehicle}</small></td>
              <td>
                {#if player.RaceFinished}
                  🎉
                {:else}
                  {player.RaceLap} / {state.Laps}
                {/if}
                <br>
                <small>{JSON.stringify(player.Position)}</small>
                <br>
                <small>Pitlane: {player.PitLane}</small>
              </td>
              <td>{player.NumStops}</td>
              <td>{JSON.stringify(player.Gaps)}</td>
              <td>{player.BTime}</td>
              <td>{player.TTime}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
    <div class="Map-right">
      <div class="Map-map">
        map
      </div>
      <div class="Map-chat">
        <ul>
          {#each messages as message}
            <li>{message}</li>
          {/each}
        </ul>
      </div>
    </div>
  </div>
  <footer class="Map-footer">
    Track: {state.Track.Name} ({state.Track.Code}), Weather: {state.Weather}, Wind: {state.Wind}
  </footer>
</div>
