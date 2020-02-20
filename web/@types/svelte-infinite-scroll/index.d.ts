declare module 'svelte-infinite-scroll' {
  export default interface SvelteInfiniteScroll {
      threshold?: number
      elementScroll?: Node
      hasMore?: boolean
      loadMore?: () => void
  }
}
