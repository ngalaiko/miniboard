<script lang="ts">
  import BookOpen from '../../../icons/BookOpen.svelte'
  import List from '../../../icons/List.svelte'
  import Star from '../../../icons/Star.svelte'
  import { Link } from 'svelte-routing'

  const getProps = (props: any): any => {
    if (props.isCurrent) {
      return { class: 'menu-element menu-active' }
    }
    return { class: 'menu-element' }
  }
</script>

<div class="menu">
  <Link to="unread" getProps={getProps}>
      <BookOpen size={'1em'} />
      <div class='title'>Unread</div>
  </Link>
  <Link to="favorite" getProps={getProps}>
      <Star size={'1em'} />
      <div class='title'>Favorite</div>
  </Link>
  <Link to="all" getProps={getProps}>
      <List size={'1em'} />
      <div class='title'>All</div>
  </Link>
</div>

<style>
  .menu {
    display: flex;
    flex-direction: column;
  }

  :global(.menu-element) {
    display: flex;
    padding: 5px;
    align-items: center;
    color: inherit;
    text-decoration: none;
  }

  :global(.menu-active) {
    background: gainsboro;
  }

  .title {
    padding-left: 5px;
  }
</style>
