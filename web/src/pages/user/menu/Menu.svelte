<script lang="ts">
  import { createEventDispatcher } from 'svelte'

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
    unread
  </div>
  <div
    id="favorite"
    class="menu-element {current === 'favorite' ? 'active' : ''}"
    on:click={clicked}
  >
    favorite
  </div>
  <div
    id="all"
    class="menu-element {current === 'all' ? 'active' : ''}"
    on:click={clicked}
  >
    all
  </div>
</div>

<style>
  .menu {
    display: flex;
    flex-direction: column;
  }

  .menu-element {
    padding: 5px;
  }

  .active {
    background: gray;
  }
</style>
