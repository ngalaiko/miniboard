<script>
    import { createEventDispatcher } from 'svelte'
    import TimeAgo from '../../components/timeago/TimeAgo.svelte'

    const dispatch = createEventDispatcher()

    export let articles
    export let router

    export let article

    console.log('here', article)

    const onDeleted = async () => {
        await articles.delete(name)
        dispatch('deleted', name)
    }

    const onRead = async (isRead) => {
        let updated = await articles.update({
            name: name,
            is_read: isRead
        }, 'is_read')
        is_read = isRead
        dispatch('updated', updated)
    }

    const onStarred = async (isStarred) => {
        let updated = await articles.update({
            name: name,
            is_favorite: isStarred
        }, 'is_favorite')
        is_favorite = isStarred
        dispatch('updated', updated)
    }

    const onClick = async () => {
        onRead(true)
        router.route(`/${name}`)
    }
</script>

<div class='article' class:opacity={is_read}>
    <span class='title' on:click|preventDefault={onClick}>{title}</span>
    <ul class='article-info'>
        {#if site_name !== undefined}
        <li><a class='link padding' href={url}>{site_name}</a></li>
        {:else}
        <li><a class='link padding' href={url}>original</a></li>
        {/if}
        <li class='separator flex'><TimeAgo date={create_time}/></li>
        {#if is_favorite}
            <li class='separator'>
                <button on:click|preventDefault={() => onStarred(false)}><b>star</b></button>
            </li>
        {:else}
            <li class='separator'>
                <button on:click|preventDefault={() => onStarred(true)}>star</button>
            </li>
        {/if}
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
        display: flex;
        flex-direction: column;
    }
    
    .title {
        font-size: 1.1em;
        font-weight: 500;
        color: inherit;
        text-decoration: none;
        cursor: pointer;
        margin-bottom: 3px;
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
        text-decoration: none;
    }

    .link:hover {
        text-decoration: underline;
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
