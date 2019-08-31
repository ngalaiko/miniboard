<script>
    import { createEventDispatcher } from 'svelte'

    export let name
    export let url
    export let api

    const dispatch = createEventDispatcher()
    
    function onDelete() {
      api.delete(`/api/v1/${name}`)
        .then(() => dispatch('deleted', name))
    }
</script>

<div>
    <a>{url}</a>
    <button on:click|preventDefault={onDelete}>x</button>
</div>
