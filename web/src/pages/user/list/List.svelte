<script lang="ts">
  import SvelteInfiniteScroll from 'svelte-infinite-scroll'
  // @ts-ignore
  import { ArticlesClient, Articles, Article } from '../../../clients/articles.ts'
  // @ts-ignore
  import { SourcesClient, } from '../../../clients/sources.ts'
  import ArticleView from './article/Article.svelte'
  import { createEventDispatcher , onMount, onDestroy } from 'svelte'
  import Search from '../../../icons/Search.svelte'
  import Add from '../../../icons/Add.svelte'
	import Modal from './Modal.svelte'
  import Selector from './Selector.svelte'
  // @ts-ignore
  import { Category, Categories } from './Category.ts'

	let showModal = false

	const dispatch = createEventDispatcher()

  export let username: string = ''
  export let articlesClient: ArticlesClient
  export let sourcesClient: SourcesClient

  let pageToken = ''
  export let category: Category = Categories.All
  let hasMore = true
  let articlesList: Article[] = []

  let searchQuery: string|null = ''
  let selectedArticleName = ''

  $: history.pushState(null, "", `?${searchQuery ? 'searchQuery='+searchQuery : ''}#${selectedArticleName}`)
  $: selectedArticleName = location.hash.slice(1)
  $: searchQuery = new URLSearchParams(location.search).get('searchQuery')

  const listParams = category.listParams
  const loadMore = async () => {
    const articles = await articlesClient.list(
      username, 25, pageToken, listParams.withTitle(searchQuery))

    articlesList = [
      ...articlesList,
      ...articles.articles,
    ]

    pageToken = articles.nextPageToken
    hasMore = pageToken !== ''
  }

  const onSelected = (article: Article) => {
    selectedArticleName = article.name
    dispatch('selected', article.name)
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

  const onAdd = async (url: string) => {
    await sourcesClient.create(username, url)
    refresh()
  }

  onMount(loadMore)
  onDestroy(() => dispatch('selected', null))
</script>

<div class="list">
  <div class="list-header">
      <button class="button-search" on:click={() => searchQuery = searchQuery == null ? '' : null}>
      <Search size="1em" />
    </button>
    <input
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
        on:select={(e) => dispatch('select_category', e.detail)} 
      />
    </div>
  </div>
  <ul class="list-ul">
    {#each articlesList as article}
      <li class="list-li {article.name === selectedArticleName ? 'selected' : ''}">
        <ArticleView
          article={article}
          on:click={(e) => onSelected(article)}
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
    <Add />
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
    justify-content: space-between;
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
    padding: 0;
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
