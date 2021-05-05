<script>
  import {
    Row,
    Column,
  } from "carbon-components-svelte";

  import Players from '../components/Players.svelte'
  import Map from '../components/Map.svelte'
</script>


<Row>
<Column sm={16} md={6} lg={6}>
<Map/>
</Column>
<Column sm={16} md={10} lg={10}>
<Players/>
</Column>
</Row>
