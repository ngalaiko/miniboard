<script lang="ts">
  import { navigate } from "svelte-routing"

  // @ts-ignore
  import { CodesClient } from '../../clients/codes.ts'

  export let codesClient: CodesClient

  let email = ''

  const urlParams = new URLSearchParams(window.location.search)
  let error = urlParams.get('error')

  let showCode = false
  let code = ''

  const handleClick = async () => {
    if (email == '' ) {
        return
    }
    let resp = codesClient.sendCode(email)
    error = ''
    showCode = true
  }

  const handleCode = async () => navigate(`/codes/${code}`)
</script>

<svelte:head>
  <title>Miniboard</title>
</svelte:head>

<div id='login'>
  {#if error}
    <div class='alert'>{error}</div>
  {/if}
  <form>
    <input
      name='email'
      type='email'
      bind:value={email}
      placeholder='email'
    />
    <button on:click|preventDefault={handleClick} />
  </form>
  <form>
  {#if showCode}
    <form>
      <input
        name='code'
        type='text'
        bind:value={code}
        placeholder='code from the email'
      />
      <button on:click|preventDefault={handleCode} />
    </form>
  {/if}
</div>

<style>
  #login {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    margin: auto;
  }

  button {
    padding-left: 0; padding-right: 0;
    border-left-width: 0; border-right-width: 0;
    white-space: nowrap;
    overflow: hidden;
  }

  input {
    border: 1px solid;
    width: 100%;
    font-size: 1.1em;
    padding: 5px;
    padding-left: 7px;
    -webkit-appearance: none;
    border-radius: 0;
  }

  input:focus {
    outline: none;
    outline-width: 0;
  }

  .alert {
    color: #b94a48;
    background-color: #f2dede;
    border-color: #eed3d7;
    padding: 8px;
    border: 1px solid #fbeed5;
    overflow: auto;
    margin: 10px;
  }
</style>
