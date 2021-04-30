<script>
  import { state } from '../stores.js'
  import Duration from './Duration.svelte'
  import Colours from './Colours.svelte'
  import Speed from './Speed.svelte'
</script>
<style>
  th, td {
    @apply p-2;
  }
</style>
<table class="w-full text-left border-collapse">
  <thead>
    <tr>
      <th></th>
      <th>P</th>
      <th>Driver<br><small>Car</small></th>
      <th>Lap</th>
      <th class="text-center">Gap</th>
      <th class="text-right">Best</th>
      <th class="text-right">Total</th>
    </tr>
  </thead>
  <tbody class="">
    {#each Object.values($state.Players).filter((a) => { return a.RacePosition > 0; }).sort((a, b) => { return a.RacePosition - b.RacePosition; }) as player}
      <tr>
        <td>
          {#if player.RaceFinished}
            🏁
          {/if}
        </td>
        <td>{player.RacePosition}</td>
        <td>
          <Colours string={player.Playername} />
          <br>
          <small>{player.Vehicle}</small>
        </td>
        <td>
          {player.RaceLap} / {$state.Laps}
          <br/>
          <small><Speed speed={player.Position.Speed}/></small>
          <br/>
          <small>{player.NumStops} Stops</small>
          {#if player.PitLane}
            <br/><small>In Pit Lane</small>
          {/if}
        </td>
        <td class="text-center">
          {#if player.Gaps.Next}
            {player.Gaps.Next}
          {/if}
        </td>
        <td class="text-right"><Duration duration={player.BTime}/></td>
        <td class="text-right"><Duration duration={player.TTime}/></td>
      </tr>
    {/each}
  </tbody>
</table>
