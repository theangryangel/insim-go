<script>
  import { state } from '../stores.js'
  import Duration from './Duration.svelte'
  import Colours from './Colours.svelte'
  import { DataTable, Tag } from "carbon-components-svelte";

  let headers = [
    { key: 'RacePosition', value: 'P' },
    { key: 'Playername', value: 'Player' },
    { key: 'RaceLap', value: 'Lap' },
    { key: 'Gaps.Next', value: 'Gap' },
    { key: 'S1', value: 'Split 1' },
    { key: 'S2', value: 'Split 1' },
    { key: 'S3', value: 'Split 1' },
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

  $: rows = getRows($state.Players);

</script>
<style>
</style>

<DataTable
  size="short"
  expandable
  headers={headers}
rows={rows}>

  <div slot="expanded-row" let:row>
    <pre>
      {JSON.stringify(row.Position, null, 2)}
    </pre>
  </div>

  <span slot="cell" let:row let:cell>
    {#if ['BTime', 'TTime', 'LTime'].includes(cell.key)}
      <Duration duration={cell.value}/>
    {:else if cell.key == 'Playername'}
      <Colours string={cell.value}/> {#if row.RaceFinished}🏁{/if} {#if row.PitLane}<Tag>Pitlane</Tag>{/if}
    {:else if cell.key == 'RaceLap'}
      {cell.value} / {$state.Laps}
    {:else}{cell.value}{/if}
  </span>

</DataTable>
