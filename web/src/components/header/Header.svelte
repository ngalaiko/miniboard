<script>
    export let api

    function onLogout() {
        api.logout()
        location.href = '/'
    }
</script>

<div>
    <div on:click={onLogout}>Logout</div>
</div>
