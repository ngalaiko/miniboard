<script>
  import Menu from './menu/Menu.svelte'
  import List from './list/List.svelte'
  import Reader from './reader/Reader.svelte'
</script>

<div id="user">
  <div id="menu" class="column">
    <Menu />
  </div>
  <div id="list" class="column">
    <List />
  </div>
  <div id="reader" class="column">
    <Reader />
  </div>
</div>

<style>
  #user {
    display: flex;
    width: 100%;
  }

  #list {
    flex-basis: 20%;
  }

  #menu {
    flex-basis: 20%;
  }

  #reader {
    flex-basis: 60%;
  }

  .column {
    border-left: 1px solid;
    border-right: 1px solid;
  }
</style>
