<script>
    import Articles from '../articles/Articles.svelte'

    export let username
    export let api

    let user = null

    api.get(`/api/v1/users/${username}`)
        .then(resp => { user = resp })
        .catch(err => console.error(err))
</script>

<div>
    {#if (user != null)}
    welcome, {user.name}
    <Articles api={api} user={user} />
    {/if}
</div>
