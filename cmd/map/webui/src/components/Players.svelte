<script>
  import { players } from '../stores.js'
  import Duration from './Duration.svelte'
  import Colours from './Colours.svelte'
  import { DataTable, Tag } from "carbon-components-svelte";

  let headers = [
    { key: 'RacePosition', value: 'P' },
    { key: 'id', value: '#' },
    { key: 'Playername', value: 'Player' },
    { key: 'RaceLap', value: 'Lap' },
    { key: 'GapNext', value: 'Next' },
    { key: 'CurrentLapTiming.Split[0].Time', value: 'S1' },
    { key: 'CurrentLapTiming.Split[1].Time', value: 'S2' },
    { key: 'CurrentLapTiming.Split[2].Time', value: 'S3' },
    { key: 'BTime', value: 'Best Time' },
    { key: 'LTime', value: 'Last Time' },
    { key: 'TTime', value: 'Total Time' },

  ];

  const getRows = (players) => {
    let output = [];
    const sorted = Object.entries(players).filter((a) => { return a[1].RacePosition > 0; }).sort((a, b) => { return a[1].RacePosition - b[1].RacePosition; })

    for (const [key, value] of sorted) {
      output.push({ id: key, ...value});
    }
    return output;
  }

  $: rows = getRows($players);

</script>
<style>
</style>

<DataTable
  size="compact"
  zebra={true}
  headers={headers}
  rows={rows}>

  <span slot="cell" let:row let:cell>
    {#if ['BTime', 'TTime', 'LTime'].includes(cell.key)}
      <Duration duration={cell.value}/>
    {:else if cell.key == 'Playername'}
      <Colours string={cell.value}/> {#if row.RaceFinished}🏁{/if} {#if row.PitLane}<Tag type="teal" size="sm">Pitlane</Tag>{/if}
    {:else if cell.key.includes("Split")}
      <Duration duration={cell.value}/>
    {:else}{cell.value}{/if}
  </span>

</DataTable>
