<script>
    import { createEventDispatcher } from 'svelte'
    import { onMount } from 'svelte';

    export let editable = false
    export let name
    export let title

    const dispatch = createEventDispatcher()

    function onKeyDown(e) {
        if (!e) {
            e = window.event;
        }
        var keyCode = e.which || e.keyCode,
            target = e.target || e.srcElement;

        if (keyCode !== 13) { // Enter
            return
        }

        dispatch('created', this.value)
        editable = false
    }

    function onKeyPress() {
        this.style.width = 0 
        this.style.width = 5 + this.scrollWidth + 'px'
    }

    function onDelete() {
        dispatch('deleted', name)
    }

	onMount(() => {
        let i = document.getElementById(name)
        i.style.width = 0
        i.style.width = 7 + i.scrollWidth + 'px'
	})
</script>

<span class='container'>
    <input id={name} class='label' disabled={!editable} value={title} on:input={onKeyPress} on:keydown={onKeyDown} />
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
        background: inherit;
        color: inherit;
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

    .labal:disabled {
    }
</style>
