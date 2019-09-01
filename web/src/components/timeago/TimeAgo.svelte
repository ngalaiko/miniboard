<script>
  export let date

  function timeAgo() {
    let seconds = Math.floor((Date.now() - new Date(date)) / 1000)
    let minutes = Math.floor(seconds / 60)
    if (minutes == 0) {
      return `${seconds}s ago`
    }
    let hours = Math.floor(minutes / 60)
    if (hours == 0) {
      return `${minutes}m ago`
    }
    let days = Math.floor(hours / 24)
    if (days == 0) {
      return `${hours}h ago`
    }
    return `${days}d ago`
  }

  let ago = timeAgo()

  let intervalID = window.setInterval(() => {
    ago = timeAgo()
  }, 1000)
</script>

<div>
    {ago} 
</div>
