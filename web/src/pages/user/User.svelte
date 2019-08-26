<script>
    export let username

    export let api;

    api.get(`/api/v1/users/${username}`)
        .then(resp => console.log(resp))
        .catch(err => console.error(err))
</script>

<div>
    welcome, {username}
</div>

