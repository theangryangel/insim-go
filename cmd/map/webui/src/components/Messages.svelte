<script>
import { messages } from '../stores.js'
</script>
<style>
</style>
<ul>
  {#each $messages as message}
    <li class="break-all">{message}</li>
  {/each}
</ul>
