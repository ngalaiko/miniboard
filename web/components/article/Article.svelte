<script>
    import { createEventDispatcher } from 'svelte'
    import TimeAgo from '../../components/timeago/TimeAgo.svelte'

    const dispatch = createEventDispatcher()

    export let articles
    export let router

    export let name
    export let url
    export let title
    export let create_time
    export let icon_url
    export let label_ids
    export let is_read

    if (label_ids === undefined) {
        label_ids = []
    }

    const onDeleted = async () => {
        await articles.delete(name)
        dispatch('deleted', name)
    }

    const onRead = async (isRead) => {
        articles.update({
            name: name,
            is_read: isRead
        }, 'is_read')
        is_read = isRead
    }

    const onClick = async () => {
        onRead(true)
        router.route(`/${name}`)
    }
</script>

<div class='article' class:opacity={is_read}>
    <a class='title' on:click|preventDefault={onClick}>{title}</a>
    <ul class='article-info'>
        <li><a class='link padding' href={url}>original</a></li>
        <li class='separator flex'><TimeAgo date={create_time}/></li>
        {#if is_read}
            <li class='separator'>
                <button on:click|preventDefault={() => onRead(false)}>unread</button>
            </li>
        {:else}
            <li class='separator'>
                <button on:click|preventDefault={() => onRead(true)}>read</button>
            </li>
        {/if}
        <li class='separator'>
            <button on:click|preventDefault={onDeleted}>delete</button>
        </li>
    </ul>
</div>

<style>
    .opacity {
        opacity: 0.5;
    }

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
    
    button:focus {
        outline-width: 0;
    }

    button:hover {
        outline-width: 0;
        text-decoration: underline;
    }
</style>
