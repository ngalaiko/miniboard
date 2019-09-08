<script>
    import { createEventDispatcher } from 'svelte'

    const dispatch = createEventDispatcher()

    export let name
    export let title

    function onDelete() {
        dispatch('deleted', name)
    }
</script>

<span class='container'>
    <span class='label'>
		{title}
	</span>
    <button class='button-delete' on:click|preventDefault={onDelete}>X</button>
</span>

<style>
    .container {
        display: flex;
        border: 0px;
        border-radius: 10px;
        font-size: 0.8em;
        padding: 0 5px;
        border: 1px solid;
        margin: 0px;
        margin-left: 3px;
        margin-bottom: 3px;
    }

    .label {
        border: 0px;
        border-right: 1px solid;
        padding: 0px;
        padding-right: 5px;
        min-width: 20px;
        background: inherit;
        color: inherit;
        margin: auto;
        text-align: center;
    }

    .button-delete {
        background: inherit;
        border: 0px;
        border-radius: unset;
        padding: 0px;
        margin: 0px;
        cursor: pointer;
        margin: auto;
        text-align: center;
        padding-left: 3px;
    }

    .button-delete:hover, .button-delete:focus, .label:hover, .label:focus {
        outline-width: 0;
        outline: none;
    }
</style>
