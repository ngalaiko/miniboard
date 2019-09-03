<script>
    import { createEventDispatcher } from 'svelte'

    export let editable = false
    export let name
    export let title

    const dispatch = createEventDispatcher()

    // prevents linebreak
    function onKeyDown(e) {
        if (!e) {
            e = window.event;
        }
        var keyCode = e.which || e.keyCode,
            target = e.target || e.srcElement;

        if (keyCode !== 13) { // Enter
            return
        }

        if (e.preventDefault) {
            e.preventDefault();
        } else {
            e.returnValue = false;
        }

        dispatch('created', this.innerText)
        editable = false
    }

    function onDelete() {
        dispatch('deleted')
    }
</script>

<span class='container'>
    <div contenteditable={editable} class='label' on:keydown={onKeyDown}>{title}</div>
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
        display: inline-block;
        vertical-align: text-top;
        border-right: 1px solid;
        padding-right: 5px;
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
