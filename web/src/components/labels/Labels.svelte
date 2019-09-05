<script>
    import Label from '../label/Label.svelte'
    import { createEventDispatcher } from 'svelte'

    export let tips

    export let labels
    const dispatch = createEventDispatcher()

    export let labelIds
    let labelsList = []

    labelIds.forEach(async (labelName) => {
        let label = await labels.get(labelName)
        labelsList = labelsList.concat([label])
    })

    function onAdd() {
        labelsList = labelsList.concat([{
            title: 'enter name',
            editable: true,
        }])
    }

    function onDeleted(e) {
        if (e.detail != null) {
            dispatch('labeldeleted', e.detail)
        }
        this.$destroy();
    }

    async function onCreated(e) {
		let create = true
		// if the label already exists, don't add it again
		labelsList.forEach(label => {
			if (label.title == e.detail) {
				this.$destroy()
				create = false
			}
		})

		if (!create) {
			return
		}

        let label = await labels.create(e.detail)

        dispatch('labeladded', label.name)
    }
</script>

<span class='container'>
    {#each labelsList as label}
        <Label
            {...label}
            on:deleted={onDeleted}
            on:created={onCreated}
            tips={labels.titles}
        />
    {/each}
    <button class='button-add' on:click|preventDefault={onAdd}>➕</button>
</span>

<style>
    .container {
        display: inline-flex;
        align-items: center;
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
</style>
