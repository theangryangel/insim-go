<script>
  import {
    Header,
    HeaderUtilities,
    HeaderNav,
    HeaderNavItem,
    HeaderNavMenu,
    HeaderGlobalAction,
    SkipToContent,
    Content,
    Grid,
    Row,
    Column,
  } from "carbon-components-svelte";
  import ChoroplethMap20 from "carbon-icons-svelte/lib/ChoroplethMap20";
  import Chat20 from "carbon-icons-svelte/lib/Chat20";

  import { onMount } from 'svelte'

  import Messages from './components/Messages.svelte';
  import TrackInfo from './components/TrackInfo.svelte'
  import Players from './components/Players.svelte'
  import Map from './components/Map.svelte'
  import Leaflet from './components/Leaflet.svelte'
  import { state, messages } from './stores.js'

  let showMap = true;
  let showChat = true;

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
      let data = JSON.parse(ev.data)

      $messages.unshift(data)
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

<Header company="insim.go" platformName="LiveMap">
  <div slot="skip-to-content">
    <SkipToContent />
  </div>

  <HeaderUtilities>
  <HeaderGlobalAction aria-label="Map" icon={ChoroplethMap20} isActive={showMap} on:click={() => { showMap = !showMap }}/>
  <HeaderGlobalAction aria-label="Map" icon={Chat20} isActive={showChat} on:click={() => { showChat = !showChat }}/>
  </HeaderUtilities>
</Header>

<Content>
  <Grid>
    <Row>
      {#if showMap}
      <Column sm={16} md={6} lg={6}>
        <Map/>
      </Column>
      {/if}
      <Column sm={16} md={10} lg={10}>
        <Players/>
      </Column>
    </Row>

    {#if showChat}
    <Row>
      <Column>
        <Messages />
      </Column>
    </Row>
    {/if}
  </Grid>
</Content>
