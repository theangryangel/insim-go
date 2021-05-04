<script>
  import {
    Header,
    HeaderUtilities,
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
  import Players from './components/Players.svelte'
  import Map from './components/Map.svelte'
  import { SSE } from './sse.js'

  let showMap = true;
  let showChat = true;

  let es = new SSE();

  onMount(async () => {
    es.dial()
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
  <Grid noGutter={true}>
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
