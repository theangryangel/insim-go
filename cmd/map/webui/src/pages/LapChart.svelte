<script>
  import { event, players } from '../stores.js'

  import {
    Row,
    Column,
    DataTable
  } from "carbon-components-svelte";

  import Colours from '../components/Colours.svelte'

  const getHeaders = (event) => {
    let headers = [
      { key: 'Playername', value: 'Player' },
    ]

    for (let i = 0; i <= event.Laps; i++) {
      headers.push({ key: i.toString(), value: 'Lap ' + i })
    }

    return headers
  }

  const getRows = (players, laps) => {
    let output = []

    const sorted = Object.entries(players).filter((a) => { return a[1].RacePosition > 0; }).sort((a, b) => { return a[1].RacePosition - b[1].RacePosition; })

    for (const [key, player] of sorted) {
      let poutput = {
        id: key,
        Playername: player.Playername,
        Vehicle: player.Vehicle,
        "0": player.RaceStartPosition || - '-',
      }

      for (let i = 1; i<= laps; i++) {
        poutput[i.toString()] = player.LapTimings[i]?.Split?.[3]?.RacePosition || '-'
      }

      output.push(poutput);
    }
    return output;
  }

  $: headers = getHeaders($event)
  $: rows = getRows($players, $event.Laps);
</script>
<style>
</style>



<Row>
  <Column>
    <DataTable
      size="compact"
      zebra={true}
      headers={headers}
      rows={rows}>

      <span slot="cell" let:row let:cell>
        {#if cell.key == 'Playername'}
          <Colours string={cell.value}/> {#if row.RaceFinished}🏁{/if}
          <br/><small>{row.Vehicle}</small>
        {:else}{cell.value}{/if}
      </span>

    </DataTable>
  </Column>
</Row>
