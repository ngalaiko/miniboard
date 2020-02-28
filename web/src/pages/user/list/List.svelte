<script lang="ts">
  import SvelteInfiniteScroll from 'svelte-infinite-scroll'
  // @ts-ignore
  import { ArticlesClient, Articles, Article, ListParams } from '../../../clients/articles.ts'
  import ArticleView from './article/Article.svelte'
  import { createEventDispatcher , onMount, onDestroy } from 'svelte'

	const dispatch = createEventDispatcher()

  export let username: string = ''
  export let articlesClient: ArticlesClient
  export let listParams: ListParams = new ListParams(false, false)

  let pageToken = ''
  let hasMore = true
  let articlesList: Article[] = []

  let query: string|null = ''
  let selectedArticleName = ''

  $: history.pushState(null, "", `?query=${query ? query : ''}#${selectedArticleName}`)
  $: selectedArticleName = location.hash.slice(1)
  $: query = new URLSearchParams(location.search).get('query')

  const loadMore = async () => {
    const articles = await articlesClient.list(
      username, 25, pageToken, listParams.withTitle(query))

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

  onMount(loadMore)
  onDestroy(() => dispatch('selected', null))
</script>


<div class="list">
  <input
    class="search-input"
    placeholder="filter"
    bind:value={query}
    on:change={onInput}
    on:input={onInput}
    on:cut={onInput}
    on:copy={onInput}
    on:paste={onInput}
  />
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
</div>

<style>
  .list {
    display: flex;
    flex-direction: column;
    max-height: 100%;
  }

  .search-input {
    font: inherit;
    border: 0;
    padding: 5px;
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
</style>
