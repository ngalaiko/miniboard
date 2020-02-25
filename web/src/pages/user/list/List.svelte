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

  const loadMore = async () => {
    const articles = await articlesClient.list(username, 25, pageToken, listParams)

    articlesList = [
      ...articlesList,
      ...articles.articles,
    ]

    pageToken = articles.nextPageToken
    hasMore = pageToken !== ''
  }

  let selectedArticleName = ''
  $: selectedArticleName = location.hash.slice(1)

  const onSelected = (article: Article) => {
    selectedArticleName = article.name
    history.pushState(null, "", `#${article.name}`)
    dispatch('selected', article.name)
  }

  onMount(loadMore)
  onDestroy(() => dispatch('selected', null))
</script>

<ul class="list">
  {#each articlesList as article}
    <li class="list-element {article.name === selectedArticleName ? 'selected' : ''}">
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

<style>
  .list {
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    max-height: 100%;
    overflow-y: scroll;
  }

  .list-element {
    border-bottom: 1px solid;
    padding-right: 5px;
  }

  .selected {
    background: gainsboro;
  }
</style>
