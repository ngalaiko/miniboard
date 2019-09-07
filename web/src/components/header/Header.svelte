<script>
    export let api
    export let router

    function onLogout() {
        api.logout()
        router.route('/')
    }
</script>

<div>
    <button on:click|preventDefault={onLogout}>logout</button>
</div>

<style>
  div {
    display: flex;
    flex-direction: row-reverse;
    margin-top: 15px;
    margin-bottom: 20px;
  }

  button {
      background: inherit;
      -webkit-appearance: none;
      -moz-appearance: none;
      font-size: 1.1em;
      cursor: pointer;
      padding: 3px 10px;
      border: 1px solid;
      border-radius: unset;
  }

  button:hover, button:focus {
      outline-width: 0;
      text-decoration: underline;
  }
</style>
