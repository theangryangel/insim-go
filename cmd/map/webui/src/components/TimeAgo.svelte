<script>
	import { format } from 'timeago.js';
	export let date = new Date();
</script>

<time datetime={date} title={date}>{format(date)}</time>
