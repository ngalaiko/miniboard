<script>
    import { createEventDispatcher } from 'svelte'
    import { onMount } from 'svelte';

    export let editable = false
    export let name
    export let title

    export let tips = []

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
		this.style.width = 2 + this.value.length + 'ch'
    }

    function onDelete() {
        dispatch('deleted', name)
    }

	onMount(() => {
		let labels = document.getElementsByClassName(name)
		Array.prototype.forEach.call(labels, (label) => {
			label.style.width = 0
			label.style.width = 2 + label.value.length + 'ch'
		})
	})

	let randId = 'label-' + Math.random()
</script>

<span class='container'>
    <input
		class='label {name}'
		disabled={!editable}
		value={title}
 		on:input={onKeyPress}
		on:keydown={onKeyDown}
		list={randId}
	/>
	<datalist id={randId}>
		{#each tips as tip}
			<option value={tip}></option>
		{/each}
    </datalist>
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
