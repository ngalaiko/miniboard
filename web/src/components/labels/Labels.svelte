<script>
    import Label from '../label/Label.svelte'
    import { createEventDispatcher } from 'svelte'

    export let labelIds

    export let articleName

    export let labelsService
    const dispatch = createEventDispatcher()

    let labels = []
    labelIds.forEach(labelName => {
        labelsService.get(labelName)
            .then(label => {
                labels = labels.concat([label])
            })
    })

    function onAdd() {
        labels = labels.concat([{
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

    function onCreated(e) {
        labelsService.create(e.detail)
            .then(resp => { dispatch('labeladded', resp.name) })
    }
</script>

<span class='container'>
    {#each labels as label}
        <Label {...label} on:deleted={onDeleted} on:created={onCreated} />
    {/each}
    <button class='button-add' on:click|preventDefault={onAdd}>➕</button>
</span>

<style>
    .container {
        vertical-align: middle;
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
