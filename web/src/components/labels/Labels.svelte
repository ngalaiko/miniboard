<script>
    import { createEventDispatcher } from 'svelte'
    import Label from '../label/Label.svelte'

    export let labels
    const dispatch = createEventDispatcher()

    export let labelIds
    let labelsList = []

    let randId = Math.floor(Math.random() * 10000)

    labelIds.forEach(async (labelName) => {
        let label = await labels.get(labelName)
        labelsList = labelsList.concat([label])
    })

    function onAdd() {
        let input = document.querySelector(`.input-${randId}`)
        input.hidden = false
        input.focus()
    }

    function onDeleted(e) {
        if (e.detail != null) {
            dispatch('labeldeleted', e.detail)
        }
        this.$destroy();
    }

    async function onKeyDown(e) {
        let keyCode = e.which || e.keyCode,
            target = e.target || e.srcElement;

        if (keyCode != 13) { // enter
            return
        }

        let labelTitle = target.value 

        target.hidden = true
        target.value = ''

		let create = true
		// if the label already exists, don't add it again
		labelsList.forEach(label => {
			if (label.title == labelTitle) {
				this.$destroy()
				create = false
			}
		})

		if (!create) {
			return
		}

        let label = await labels.create(labelTitle)
        labelsList = labelsList.concat([label])

        dispatch('labeladded', label.name)
    }

    function onKeyPress() {
        this.style.width = 0 
		this.style.width = 2 + this.value.length + 'ch'
    }
</script>

<span class='container'>
    {#each labelsList as label}
        <Label
            {...label}
            on:deleted={onDeleted}
        />
    {/each}
    <input
        class='input-new-label input-{randId}'
        hidden=true
 		on:input={onKeyPress}
		on:keydown={onKeyDown}
    />
    <button class='button-add' on:click|preventDefault={onAdd}>➕</button>
</span>

<style>
    .container {
        display: inline-flex;
        align-items: center;
        flex-wrap: wrap;
    }

    .button-add {
        vertical-align: text-top;
        background: inherit;
        border: 0px;
        border-radius: unset;
        padding: 0px;
        margin: 0px;
        cursor: pointer;
    }

    button-add:hover, .button-add:focus {
        outline-width: 0;
    }

    .input-new-label {
        margin-left: 5px;
        width: 1ch;
        border: 0;
    }

    .input-new-label:focus {
        outline: none;
        outline-width: 0;
    }
</style>
