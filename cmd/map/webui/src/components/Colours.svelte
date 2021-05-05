<script>
  export let string

  const colours = {
    0: 'text-gray-300', // this should be black really
    1: 'text-red-500',
    2: 'text-green-200',
    3: 'text-yellow-300',
    4: 'text-blue-500',
    5: 'text-purple-500',
    6: 'text-blue-300',
    7: 'text-white',
    8: '', // default colour
    9: '' // reset
  }

  function sanitize(string) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#x27;',
        "/": '&#x2F;',
    };
    const reg = /[&<>"'/]/ig;
    return string.replace(reg, (match)=>(map[match]));
  }

  // Courtesy of https://apidocs.tc-gaming.co.uk/guides/converting-lfs-colours
  // with minor modifications
  function colourise(str) {
    var parts = str.split(/(\^[0-9])/g)

    if (parts.length < 2) {
      return str;
    }

    var res = "";
    parts.slice(1).forEach(function(el, i, arr) {
      (i % 2 === 0) ? arr[i] = el.slice(1) : res += '<span class="' + colours[arr[i-1]] + '">' + el + '</span>';
    });

    return res
  }
</script>
<style>
</style>
{@html colourise(sanitize(string))}
