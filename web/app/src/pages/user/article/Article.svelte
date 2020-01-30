<script>
    import { createEventDispatcher } from 'svelte'
    import { navigate } from "svelte-routing"
    import TimeAgo from './timeago/TimeAgo.svelte'

    const dispatch = createEventDispatcher()

    export let article

    const onDeleted = async () => dispatch('deleted', article.getName())

    const onRead = async (isRead) => {
        article.setIsRead(isRead)
        dispatch('updated', article)
    }

    const onStarred = async (isStarred) => {
        article.setIsFavorite(isStarred)
        dispatch('updated', article)
    }

    const onClick = async () => {
        onRead(true)
        navigate(`/${article.getName()}`)
    }
</script>

<div class='article' class:opacity={article.getIsRead()}>
    <span class='title' on:click|preventDefault={onClick}>{article.getTitle()}</span>
    <ul class='article-info'>
        {#if article.getSiteName() !== ''}
        <li><a class='link padding' href={article.getUrl()}>{article.getSiteName()}</a></li>
        {:else}
        <li><a class='link padding' href={article.getUrl()}>source</a></li>
        {/if}
        <li class='separator flex'><TimeAgo date={article.getCreateTime()}/></li>
        {#if article.getIsFavorite()}
            <li class='separator'>
                <button on:click|preventDefault={() => onStarred(false)}><b>star</b></button>
            </li>
        {:else}
            <li class='separator'>
                <button on:click|preventDefault={() => onStarred(true)}>star</button>
            </li>
        {/if}
        {#if article.getIsRead()}
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
