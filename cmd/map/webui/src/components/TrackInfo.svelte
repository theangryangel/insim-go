<script>
import { state } from '../stores.js'
</script>
<style>
</style>
<div>
{$state.Event.Track.Name} ({$state.Event.Track.Code})

Weather:

{#if $state.Weather <= 1}
☀️
{/if}

{#if $state.Weather == 2}
🌤️
{/if}

{#if $state.Weather == 3}
☁️
{/if}

{#if $state.Wind == 1}
, Low Wind
{/if}

{#if $state.Wind == 2}
, Hard Wind
{/if}

</div>
