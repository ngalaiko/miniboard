<script context='module'>
    import { writable } from 'svelte/store'
    const selectedPaneStore = writable(2)
</script>

<script>
    import { quintOut } from 'svelte/easing'
    import { onDestroy } from 'svelte'

    import BookOpen from '../../../icons/BookOpen.svelte'
    import List from '../../../icons/List.svelte'
    import Star from '../../../icons/Star.svelte'

    let selectedPane
    const unsubscribe = selectedPaneStore.subscribe(value => selectedPane = value)  
    onDestroy(() => unsubscribe())
</script>

<span class='container'>
    {#if selectedPane == 1}
        <button class='selected border border-left'><Star size=20 /></button>
    {:else}
        <button class='border border-left' on:click={() => selectedPaneStore.set(1)}><Star size=20 /></button>
    {/if}

    {#if selectedPane == 2}
        <button class='selected border border-middle'><BookOpen size=20 /></button>
    {:else}
        <button class='border border-middle' on:click={() => selectedPaneStore.set(2)}><BookOpen size=20 /></button>
    {/if}

    {#if selectedPane == 3}
        <button class='selected border border-right'><List size=20 /></button>
    {:else}
        <button class='border border-right' on:click={() => selectedPaneStore.set(3)}><List size=20 /></button>
    {/if}
</span>

<style>
    .container {
        display: flex;
        flex-direction: row;
        justify-content: center;
    }

    .selected {
        background: black;
        color: white;
    }

    .border {
        padding-top: 0.2em;
        padding-bottom: 0.2em;
        padding-left: 0.5em;
        padding-right: 0.5em;
    }

    .border-left {
        border-top: 1px solid black;
        border-bottom: 1px solid black;
        border-left: 1px solid black;
    }

    .border-middle {
        border: 1px solid black;
    }

    .border-right {
        border-top: 1px solid black;
        border-bottom: 1px solid black;
        border-right: 1px solid black;
    }


    button {
        margin: 0px;
        padding: 0px;
        border: 0px;
        background: inherit;
        -webkit-appearance: none;
        -moz-appearance: none;
        cursor: pointer;
    }

    button:hover, .button:focus {
        outline-width: 0;
    }

    button:hover, button:focus {
        outline-width: 0;
    }
</style>
