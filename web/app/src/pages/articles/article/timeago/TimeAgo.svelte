<script>
    export let date

    const timeAgo = () =>  {
        let seconds = ~~((Date.now() - new Date(0).setUTCSeconds(date)) / 1000)
        let minutes = ~~(seconds / 60)
        if (minutes == 0) {
            return `${seconds}s ago`
        }
        let hours = ~~(minutes / 60)
        if (hours == 0) {
            return `${minutes}m ago`
        }
        let days = ~~(hours / 24)
        if (days == 0) {
            return `${hours}h ago`
        }
        let weeks = ~~(days / 7)
        if (weeks == 0) {
            return `${days}d ago`
        }
        return `${weeks}w ago`
    }

    let ago = timeAgo()

    let intervalID = window.setInterval(() => {
        ago = timeAgo()
    }, 1000)
</script>

<div>
    {ago} 
</div>
