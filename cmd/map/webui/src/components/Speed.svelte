<script>
export let speed;

function toMph(speed) {
   let mph = speed / 146.486067
   return Math.round(mph * 100) / 100
}
</script>
<style>
</style>
<span class="font-mono">
   {#if speed > 0}
      {toMph(speed)}<acronym title="Miles Per Hour">mph</acronym>
   {:else}
      -
   {/if}
</span>
