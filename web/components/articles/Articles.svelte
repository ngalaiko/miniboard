<script context="module">
    import { writable } from 'svelte/store'

    import { 
        add as addAll,
        remove as deleteAll,
        update as updateAll,
    } from './all/All.svelte'

    import { 
        add as addUnread,
        remove as deleteUnread,
        update as updateUnread,
    } from './unread/Unread.svelte'

    export const add = async (url, addFunc) => {
        let mock = {
          'url': url,
          'title': url,
          'create_time': Date.now(),
          'name': Math.random()
        }

        addAll(mock)
        addUnread(mock)

        let article = await addFunc(url)

        deleteAll(mock.name)
        deleteUnread(mock.name)

        addAll(article)
        addUnread(article)
    }

    let paneStore = writable('unread')
    export const show = (pane) => paneStore.set(pane)
</script>

<script>
    import { onDestroy } from 'svelte'

    import All from './all/All.svelte'
    import Unread from './Unread/Unread.svelte'

    export let api
    export let articles
    export let router
    export let pane

    onDestroy(paneStore.subscribe(updated => pane = updated))

    const onDeleted = async (name) => {
        deleteAll(name)
        deleteUnread(name)
    }
    const onUpdated = async (updated) => {
        updateAll(updated)
        updateUnread(updated)
    }
</script>

<div>
    {#if pane === 'unread'}
        <Unread
            articles={articles}
            router={router}
            on:updated={(e) => onUpdated(e.detail)}
            on:deleted={(e) => onDeleted(e.detail)}
        />
    {/if}
    {#if pane === 'all'}
        <All
            articles={articles}
            router={router}
            on:updated={(e) => onUpdated(e.detail)}
            on:deleted={(e) => onDeleted(e.detail)}
        />
    {/if}
</div>
