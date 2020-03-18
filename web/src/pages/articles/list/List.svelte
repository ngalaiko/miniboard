<script lang="ts">
  import SvelteInfiniteScroll from 'svelte-infinite-scroll'
  // @ts-ignore
  import { ArticlesClient, Articles, Article } from '../../../clients/articles.ts'
  // @ts-ignore
  import { SourcesClient, } from '../../../clients/sources.ts'
  import ArticleView from './article/Article.svelte'
  import { createEventDispatcher, onMount } from 'svelte'
  import { PlusIcon } from 'svelte-feather-icons'
	import AddModal from './AddModal.svelte'
  import Selector from './Selector.svelte'
  // @ts-ignore
  import { Category, Categories } from './Category.ts'
  import Search from './Search.svelte'

	let showModal = false

	const dispatch = createEventDispatcher()

  export let username: string = ''
  export let articlesClient: ArticlesClient
  export let sourcesClient: SourcesClient
  
  const categoryParam = new URLSearchParams(location.search).get('category')
  let category: Category = categoryParam ? Categories[categoryParam] : Categories['unread']

  let pageToken = ''
  let hasMore = true
  let articlesList: Article[] = []

  let searchQuery: string|null = new URLSearchParams(location.search).get('title')
  export let selectedArticleName = ''

  $: history.pushState(null, "", `?category=${category.value}`
    + `${searchQuery ? '&title='+searchQuery : ''}`)

  const loadMore = async () => {
    const articles = await articlesClient.list(
      username, 25, pageToken, category.listParams.withTitle(searchQuery))

    articlesList = [
      ...articlesList,
      ...articles.articles,
    ]

    pageToken = articles.nextPageToken
    hasMore = pageToken !== ''
  }

  const refresh = () => {
    pageToken = ''
    articlesList = []
    loadMore()
  }

  let typingTimerId: number | undefined
  const onInput = () => {
    clearTimeout(typingTimerId)
    typingTimerId = setTimeout(refresh, 300)
  }

  let inputElement: HTMLElement
  const onAddLink = async (url: string) => {
    await sourcesClient.addLink(username, url)
    refresh()
  }

  const onAddFile = async (file: File) => {
    const content = await new Promise<string>((resolve, reject) => {
      const reader = new FileReader()
      reader.readAsBinaryString(file)
      reader.onload = e => {
        if (e.target === null) return
        resolve(e.target.result as string)
      }
      reader.onerror = () => reject(new Error('failed to read file'))
    })

    await sourcesClient.addFile(username, btoa(content))

    refresh()
  }

  const onCalegorySelect = async (c: Category) => {
    category = c
    refresh()
  }

  onMount(loadMore)
</script>

<div class="list">
  <div>
    <Selector
      selectedValue={category}
      on:select={(e) => onCalegorySelect(Categories[e.detail])}
    />
    <Search 
      bind:value={searchQuery}
      on:change={onInput}
      on:input={onInput}
      on:cut={onInput}
      on:copy={onInput}
      on:paste={onInput}
    />
  </div>
  <ul class="list-ul">
    {#each articlesList as article}
      <li class="list-li {article.name === selectedArticleName ? 'selected' : ''}">
        <ArticleView
          article={article}
          on:click={(e) => {
            if (!article.isRead) {
              article.isRead = true
              articlesClient.update(article.name, { isRead: true })
            }
            dispatch('select', article.name)
          }}
        />
      </li>
    {/each}
    <SvelteInfiniteScroll
      threshold={100}
      hasMore={hasMore}
      on:loadMore={loadMore}
    />
  </ul>
  <button class="button-add" on:click={() => showModal = true}>
    <PlusIcon size="24" />
    <div>Add</div>
  </button>
</div>

{#if showModal}
  <AddModal
    on:close={() => showModal = false} 
    on:link={(e) => onAddLink(e.detail)}
    on:file={(e) => onAddFile(e.detail)}
  />
{/if}

<style>
  .list {
    display: flex;
    flex-direction: column;
    max-height: 100%;
    min-height: 100%;
  }

  .list-ul {
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    max-height: 100%;
    overflow-y: scroll;
  }

  .list-li {
    border-top: 1px solid;
    padding-right: 0.4em;
  }

  .button-add {
    display: flex;
    align-items: center;
    justify-content: center;
    font: inherit;
    border: 0;
    background: none;
    border-top: 1px solid;
    cursor: pointer;
    margin-top: auto;
    margin-right: 0;
    margin-left: 0;
    margin-bottom: 0;
    min-height: 3em;
  }

  .button-add:focus {
    outline: none;
  }

  .button-add:hover {
    background: gainsboro;
  }

  .hidden {
    display: none;
  }

  .selected {
    background: gainsboro;
  }
</style>
