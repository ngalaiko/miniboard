<script>
    import { createEventDispatcher } from 'svelte'

    export let api

    export let name
    export let url
    export let title

    const dispatch = createEventDispatcher()
    
    function onDelete() {
      api.delete(`/api/v1/${name}`)
        .then(() => dispatch('deleted', name))
    }
</script>

<div>
    {#if title != undefined}
        <a>{title}</a>
    {:else}
        <a>{url}</a>
    {/if}
    <button on:click|preventDefault={onDelete}>x</button>
</div>
