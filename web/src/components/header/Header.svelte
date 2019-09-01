<script>
    export let api

    function onLogout() {
        api.logout()
        location.href = '/'
    }
</script>

<div>
    <button on:click|preventDefault={onLogout}>Logout</button>
</div>

<style>
  div {
    display: flex;
    flex-direction: row-reverse;
    margin-bottom: 15px;
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
  
  button:hover, .button:focus {
      outline-width: 0;
      text-decoration: underline;
  }
</style>
