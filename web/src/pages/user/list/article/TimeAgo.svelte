<script lang="ts">
  import { onDestroy } from 'svelte'

  export let date = new Date()

  let ago
  let frame: number

  onDestroy(() => {
    cancelAnimationFrame(frame)
  })

  const timeAgo = (seconds: number) =>  {
    let minutes = ~~(seconds / 60)
    if (minutes == 0) return `${seconds}s`

    let hours = ~~(minutes / 60)
    if (hours == 0) return `${minutes}m`

    let days = ~~(hours / 24)
    if (days == 0) return `${hours}h`

    let weeks = ~~(days / 7)
    if (weeks == 0) return `${days}d`

    let years = ~~(days / 365)
    if (years == 0) return `${weeks}w`

    return `${years}y`
  }

  (function update() {
    frame = requestAnimationFrame(update)

    const now = performance.timing.navigationStart + performance.now()
    const secondsSince = ~~((now - date.getTime()) / 1000 )

    ago = timeAgo(secondsSince)
  }())
</script>

<div class="time">
  {ago}
</div>
