<script>
import { messages } from '../stores.js'
  import TimeAgo from './TimeAgo.svelte'
</script>
<style>
</style>
<ul>
  {#each $messages as message}
  <li class="break-all">
    <span class="text-gray-500"><TimeAgo date={message.SentAt} /></span>: {message.Msg}
  </li>
  {/each}
</ul>
