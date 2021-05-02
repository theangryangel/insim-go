<script>
import { state } from '../stores.js'

let code
let outer = []
let track = []
let startfinish = []

let scalex = 1
let scaley = 1
let transformx = 1
let transformy = 1

state.subscribe((value) => {

   if (value.Track.Code != code) {
      code = value.Track.Code
      fetch(`/api/track/${value.Track.Code}`)
         .then((res) => {
            return res.json()
         })
         .then((data) => {
            track = []
            for (var i = 0; i < data.RoadX.length && i < data.RoadY.length; i++) {
               track.push(`${data.RoadX[i].toFixed(2)},${data.RoadY[i].toFixed(2)}`)
            }

            outer = []
            for (var i = 0; i < data.OuterX.length && i < data.OuterY.length; i++) {
               outer.push(`${data.OuterX[i].toFixed(2)},${data.OuterY[i].toFixed(2)}`)
            }

            startfinish = [
               data.FinishX[0], data.FinishY[0],
               data.FinishX[1], data.FinishY[1],
            ]

            scalex = data.ScaleX
            scaley = data.ScaleY
            transformx = data.TranslateX
            transformy = data.TranslateY
         })
   }
})

</script>
<style>
.player {
   transition: all 700ms ease-in-out;
}
</style>
<svg xmlns="http://www.w3.org/2000/svg" width=1024 height=1024>
{#if outer.length > 0}
      <polygon points={outer.join(" ")} style="stroke: #059669; stroke-width:2px; fill: #059669; fill-rule: evenodd"/>
   {/if}

{#if track.length > 0}
      <polygon points={track.join(" ")} style="stroke: #1F2937; stroke-width:2px; fill: #1F2937; fill-rule: evenodd"/>
   {/if}

   {#if startfinish.length > 0}
      <line x1={startfinish[0]} y1={startfinish[1]} x2={startfinish[2]} y2={startfinish[3]} style="stroke: white; stroke-width: 1px;"/>

      <text x={startfinish[0]-15} y={startfinish[1]} style=" font-size: 15px;" text-anchor="end">🏁</text>
   {/if}

   {#each Object.entries($state.Players) as [plid, player]}
      <g class="player" transform={`translate(${(player.Position.X / 65536) * scalex + transformx}, ${(-player.Position.Y / 65536) * scaley + transformy})`}>
         <circle r="6" style="stroke: #EF4444; stroke-width: 1px;"
            width=5 height=5 fill="#111827">
         </circle>
         <text text-anchor="middle" dy="0.35em" style="font-size: 8px; stroke-width: 1px; fill: white; stroke: white">
            {plid}
         </text>
      </g>
   {/each}
</svg>
