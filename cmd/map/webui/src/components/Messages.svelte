<script>
  import { messages } from '../stores.js'

  import { UnorderedList, ListItem } from "carbon-components-svelte";

  import TimeAgo from './TimeAgo.svelte'
</script>
<style>
  h6 {
	padding-top: .4375rem;
	padding-bottom: .4375rem;
  }
</style>

<h6>Chat</h6>

<UnorderedList>
  {#each $messages as message}
    <ListItem>
      <TimeAgo date={message.SentAt} />: {message.Msg}
    </ListItem>
  {/each}
</UnorderedList>
