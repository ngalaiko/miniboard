<script>
    import Articles from '../articles/Articles.svelte'

    export let user
    export let api
</script>

<div>
    welcome, {user.name}
    <Articles api={api} user={user} />
</div>
