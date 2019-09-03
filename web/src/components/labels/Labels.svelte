<script context="module">
    import { LabelsService } from './labels-service.js'
    import Label from '../label/Label.svelte'

    export let api

    let labelsService = new LabelsService(api)
</script>

<script>
    export let labelIds

    let labels = []

    function addLabel() {
        labels = labels.concat([{
            title: 'new label'
        }])
    }
</script>

<span class='container'>
    {#each labels as label}
        <Label {...label} />
    {/each}
    <button class='button-add' on:click|preventDefault={addLabel}>➕</button>
</span>

<style>
    .container {
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
