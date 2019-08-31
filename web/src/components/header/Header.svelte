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
  }

  button {
      display: inline-block;
      -webkit-appearance: none;
      -moz-appearance: none;
      font-size: 1.1em;
      cursor: pointer;
      padding: 3px 10px;
      border: 1px solid;
      border-radius: unset;
  }
  
  button:hover, .button:focus {
      border-color: #2f5bb7;
      outline-width: 0
  }
</style>
