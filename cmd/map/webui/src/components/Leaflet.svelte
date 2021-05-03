<script>
  import {state} from '../stores.js'
  import L from 'leaflet';
  let el;
  let markers = {};
  let currentTrack;
  let map

  $: createMap(el, $state.Track.Code);

  let scalex = 1
  let scaley = 1
  let transformx = 1
  let transformy = 1

  async function createMap(container, track) {
    if (map) {
      return;
    }
    if (!container) {
      return;
    }
    if (track == currentTrack) {
      return;
    }

    const res = await fetch(`/api/track/${track}`)
    const data = await res.json()

    currentTrack = track;
    scalex = data.Fit.ScaleX
    scaley = data.Fit.ScaleY
    transformx = data.Fit.TranslateX
    transformy = data.Fit.TranslateY

    map = L.map(
			container,
			{
				crs: L.CRS.Simple,
				preferCanvas: true
			}
		);

    var bounds = [[0,0], [1024,1024]];
    var image = L.imageOverlay(data.Image, bounds).addTo(map);
    map.fitBounds(bounds);
  }

	function resizeMap() {
	  if(map) { map.invalidateSize(); }
  }

  let createMarkers = (players) => {
    if (!map) {
      return;
    }

    for (const [plid, player] of Object.entries(players)) {
      if (!plid) {
        continue
      }
      console.log(markers)
      if (markers[plid]) {
        markers[plid].setLatLng([
          (player.Position.Y / 65536) * scaley + transformy,
          (player.Position.X / 65536) * scalex + transformx,
        ])
      } else {
        markers[plid] = L.marker([
          (player.Position.Y / 65536) * scaley + transformy,
          (player.Position.X / 65536) * scalex + transformx,
        ])

        markers[plid].addTo(map).bindPopup(plid)
      }
    }
  }

  $: createMarkers($state.Players);
</script>
<svelte:window on:resize={resizeMap} />

<style></style>

<link rel="stylesheet" href="https://unpkg.com/leaflet@1.6.0/dist/leaflet.css"
   integrity="sha512-xwE/Az9zrjBIphAcBb3F6JVqxf46+CDLwfLMHloNu6KEQCAWi6HcDUbeOfBIptF7tcCzusKFjFw2yuvEpDL9wQ=="
   crossorigin=""/>
<div
   bind:this={el}
   class="map" style="height:100%;width:100%" />
