<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import BookOpen from '../../../icons/BookOpen.svelte'
  import List from '../../../icons/List.svelte'
  import Star from '../../../icons/Star.svelte'

	const dispatch = createEventDispatcher()

  const clicked = (e: Event) => {
    const target = e.target as HTMLElement
    dispatch(target.id)
  }

  export let current: string = 'unread'
</script>

<div class="menu">
  <div
    id="unread"
    class="menu-element {current === 'unread' ? 'active' : ''}"
    on:click={clicked}
  >
    <BookOpen size={'1em'} />
    <div class='title'>Unread</div>
  </div>
  <div
    id="favorite"
    class="menu-element {current === 'favorite' ? 'active' : ''}"
    on:click={clicked}
  >
    <Star size={'1em'} />
    <div class='title'>Favorite</div>
  </div>
  <div
    id="all"
    class="menu-element {current === 'all' ? 'active' : ''}"
    on:click={clicked}
  >
    <List size={'1em'} />
    <div class='title'>All</div>
  </div>
</div>

<style>
  .menu {
    display: flex;
    flex-direction: column;
  }

  .menu-element {
    display: flex;
    padding: 5px;
    align-items: center;
  }

  .active {
    background: lightgray;
  }

  .title {
    padding-left: 5px;
  }
</style>
