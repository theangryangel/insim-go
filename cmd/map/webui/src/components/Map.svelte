<script>
import { state } from '../stores.js'
import * as d3 from 'd3';
let el;
let tr;
let map;
let players = {}

state.subscribe(value => {
   console.log(value)
   if (value.Track.Code != tr) {
      d3.xml(`/map/tr`)
      .then(data => {
         map = d3.select(el).node().append(data.documentElement)
      })
   }


   if (map) {
      Object.entries($state.Players).forEach((v, i) => {
         console.log(v, i)

         if (!players[i]) {
            players[i] = map.append("svg:circle")
         }

         players[i].attr("cx", 512)
         players[i].attr("cy", 512)
      })

   }

})
</script>
<style>
</style>
<div bind:this={el} class="chart"></div>
