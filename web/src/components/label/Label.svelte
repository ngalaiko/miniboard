<script>
    import { createEventDispatcher } from 'svelte'

    export let editable = false
    export let name
    export let title

    const dispatch = createEventDispatcher()

    function onKeyPress() {
        this.style.width = 0 
        this.style.width = 5 + this.scrollWidth + 'px'
    }

    function onDelete() {
        dispatch('deleted', name)
    }
</script>

<span class='container'>
    <input class='label' value={title} on:input={onKeyPress}/>
    <button class='button-delete' on:click|preventDefault={onDelete}>x</button>
</span>

<style>
    .container {
        display: inline-block;
        vertical-align: text-top;
        border: 0px;
        border-radius: 10px;
        font-size: 0.8em;
        cursor: text;
        padding: 0 5px;
        border: 1px solid;
        margin: 0px;
        margin-left: 3px;
    }

    .label {
        border: 0px;
        display: inline-block;
        vertical-align: text-top;
        border-right: 1px solid;
        padding-right: 5px;
        min-width: 20px;
        width: 30px;
    }

    .button-delete {
        display: inline-block;
        vertical-align: text-top;
        background: inherit;
        border: 0px;
        border-radius: unset;
        padding: 0px;
        margin: 0px;
        cursor: pointer;
    }

    .button-delete:hover, .button-delete:focus, .label:hover, .label:focus {
        outline-width: 0;
    }
</style>
