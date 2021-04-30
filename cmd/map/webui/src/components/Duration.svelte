<script>
export let duration;

function nsToTime(duration) {
   duration = duration / 1000000

    var milliseconds = parseInt((duration%1000)/100)
    , seconds = parseInt((duration/1000)%60)
    , minutes = parseInt((duration/(1000*60))%60)
    , hours = parseInt((duration/(1000*60*60))%24);

   let output = minutes + "m" + seconds + "." + milliseconds + "s";
   if (hours > 0) {
      output = hours + 'h' + output
   }
   return output
}
</script>
<style>
</style>
<span>
   {#if duration > 0}
      {nsToTime(duration)}
   {:else}
      -
   {/if}
</span>
