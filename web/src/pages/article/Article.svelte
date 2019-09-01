<script>
    import { createEventDispatcher } from 'svelte'

    export let api

    export let name
    export let url
    export let title
    export let create_time
    export let iconURL

    const dispatch = createEventDispatcher()
    
    function onDelete() {
      api.delete(`/api/v1/${name}`)
        .then(() => dispatch('deleted', name))
    }
</script>

<div>
    {#if title != undefined}
        <a>{title}</a>
    {/if}
    <a>{url}</a>
    <div>{create_time}</div>
    <div>{iconURL}</div>
    <button on:click|preventDefault={onDelete}>x</button>
</div>
