<script lang="ts">
  import SvelteInfiniteScroll from 'svelte-infinite-scroll'
  // @ts-ignore
  import { ArticlesClient, Articles, Article } from '../../../clients/articles.ts'
  // @ts-ignore
  import { SourcesClient, } from '../../../clients/sources.ts'
  import ArticleView from './article/Article.svelte'
  import { createEventDispatcher, onMount } from 'svelte'
  import { PlusIcon, SearchIcon } from 'svelte-feather-icons'
	import Modal from './Modal.svelte'
  import Selector from './Selector.svelte'
  // @ts-ignore
  import { Category, Categories } from './Category.ts'

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
  const onAdd = async (url: string) => {
    await sourcesClient.create(username, url)
    refresh()
  }

  const onCalegorySelect = async (c: Category) => {
    category = c
    refresh()
  }

  onMount(loadMore)
</script>

<div class="list">
  <div class="list-header">
    <button class="button-search" on:click={() => {
      searchQuery = searchQuery == null ? '' : null
      console.log(inputElement)
      inputElement.focus()
    }}>
      <SearchIcon size="20" />
    </button>
    <input
      bind:this={inputElement}
      class="search-input {searchQuery != null ? '' : 'hidden'}"
      placeholder="search"
      bind:value={searchQuery}
      on:change={onInput}
      on:input={onInput}
      on:cut={onInput}
      on:copy={onInput}
      on:paste={onInput}
    />
    <div class="select {searchQuery == null ? '' : 'hidden'}">
      <Selector 
        selectedValue={category}
        on:select={(e) => onCalegorySelect(Categories[e.detail])} 
      />
    </div>
  </div>
  <ul class="list-ul">
    {#each articlesList as article}
      <li class="list-li {article.name === selectedArticleName ? 'selected' : ''}">
        <ArticleView
          article={article}
          on:click={(e) => dispatch('select', article.name)}
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
  <Modal 
    on:close={() => showModal = false} 
    on:add={(e) => onAdd(e.detail)}
  />
{/if}

<style>
  .list {
    display: flex;
    flex-direction: column;
    max-height: 100%;
    min-height: 100%;
  }

  .button-search {
    background: inherit;
    padding: 0;
    border: 0;
    cursor: pointer;
    font: inherit;
  }

  .button-search:focus {
    outline: none;
  }

  .list-header {
    display: flex;
    align-items: center;
    padding: 5px;
    border-bottom: 1px solid;
  }

  .search-input {
    font: inherit;
    border: 0;
    background: inherit;
    width: 100%;
    padding: 4px;
  }

  .search-input:focus {
    outline: none;
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
    border-bottom: 1px solid;
    padding-right: 5px;
  }

  .selected {
    background: gainsboro;
  }

  .button-add {
    display: flex;
    align-items: center;
    justify-content: center;
    font: inherit;
    border: 0;
    background: none;
    border-top: 1px solid;
    padding: 10px;
    cursor: pointer;
    margin-top: auto;
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

  .select {
    width: 100%;
  }
</style>
