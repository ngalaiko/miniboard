<script>
    import { createEventDispatcher } from 'svelte'
    import TimeAgo from '../../components/timeago/TimeAgo.svelte'

    const dispatch = createEventDispatcher()

    export let articles

    export let name
    export let url
    export let title
    export let create_time
    export let icon_url
    export let label_ids

    if (label_ids === undefined) {
        label_ids = []
    }
</script>

<div class='article'>
    <a href='/{name}' class='title'>{title}</a>
    <ul class='article-info'>
        <li><a class='link padding' href={url}>original</a></li>
        <li class='separator flex'><TimeAgo date={create_time}/></li>
        <li class='separator'><button on:click|preventDefault={() => dispatch('deleted', name)}>delete</button></li>
    </ul>
</div>

<style>
    .article {
        border: 1px solid;
        border-radius: unset;
        margin-bottom: 10px;
        padding: 5px;
        padding-left: 7px;
        padding-right: 7px;
    }
    
    .title {
        font-size: 1.1em;
        font-weight: 500;
        color: inherit;
        text-decoration: none;
    }

    .title:hover {
        text-decoration: underline;
    }

    .article-info {
        display: flex;
        flex-flow: row wrap;
        align-items: center;
        margin: 0px;
        padding: 0px;
        margin-top: 5px;
        font-size: 0.9em;
    }

    li {
        display: inline;
        white-space: nowrap;
    }

    .separator {
        padding-left: 0.3em;
    }

    .separator:before {
        content: "|";
        padding-right: 0.3em;
    }

    .flex {
        display: flex;
    }

    .link {
        color: inherit;
    }

    button {
        -webkit-appearance: none;
        -moz-appearance: none;
        font-size: inherit;
        border: 0px;
        padding: 0px;
        cursor: pointer;
        background: inherit;
    }
    
    button:hover, button:focus {
        outline-width: 0;
        text-decoration: underline;
    }
</style>
