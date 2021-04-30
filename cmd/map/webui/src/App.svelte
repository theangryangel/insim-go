<script>
  import { onMount } from 'svelte'

    import TrackInfo from './components/TrackInfo.svelte'
    import Messages from './components/Messages.svelte'
    import Players from './components/Players.svelte'
    import { state, messages } from './stores.js'

    let es = null;

    const dial = () => {
        es = new EventSource(`/events`);
        es.onerror = () => {
            console.log(
                `EventSource error`,
              );
            console.log("Reconnecting in 1s", true);
            setTimeout(dial, 1000);
          }
        es.onopen = () => {
            console.info("eventsource connected");
          }

        es.addEventListener("chat", (ev) => {
          $messages.unshift(ev.data)
          $messages = $messages.slice(0, 10)
        })

        es.addEventListener("state", (ev) => {
            let data = JSON.parse(ev.data)
            $state = data
          })

        es.addEventListener("player-state", (ev) => {
            let data = JSON.parse(ev.data)
            $state.Players[data.Plid] = data.State
          })

        es.addEventListener("player-left", (ev) => {
            let data = JSON.parse(ev.data)
            delete $state.Players[data.Plid]
          })
      }

    const initialState = async () => {
        return fetch("/api/state")
          .then((res) => {
              return res.json();
            })
          .then((res) => {
              $state = res;
            });
      }

    onMount(async () => {
        initialState().then(() => { dial(); console.log($state) })
      });
</script>
<style>
  th, td {
    @apply p-2;
  }
</style>
<div class="flex flex-col w-full h-screen">
  <div class="flex flex-grow bg-blue-200 flex-row overflow-y-auto">
    <div class="w-3/5 bg-gray-800 p-3 overflow-y-auto">
      <Players />
    </div>
    <div class="w-2/5 flex flex-col bg-white text-gray-900">
      <div class="h-3/4 overflow-auto">
      </div>
      <div class="h-1/4 overflow-y-auto bg-gray-200 p-3 break-all">
        <Messages />
      </div>
    </div>
  </div>
  <footer class="Map-footer flex border-red-900 border-t p-3">
    <TrackInfo />
  </footer>
</div>
