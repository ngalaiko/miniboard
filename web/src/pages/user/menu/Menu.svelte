<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte'

	const dispatch = createEventDispatcher()

  const clicked = (e: Event) => {
    const target = e.target as HTMLElement
    const parent = target.parentElement as HTMLElement
    parent.childNodes.forEach((child: ChildNode) => {
      const e = child as HTMLElement
      if (e.classList) e.classList.remove('menu-selected')
    })
    target.classList.add('menu-selected')
    dispatch(target.id)
  }
</script>

<div class="menu">
  <div id="unread" class="menu-element menu-selected" on:click={clicked}>
    unread
  </div>
  <div id="favorite" class="menu-element" on:click={clicked}>
    favorite
  </div>
  <div id="all" class="menu-element" on:click={clicked}>
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

  .menu-selected {
    background: gray;
  }
</style>
