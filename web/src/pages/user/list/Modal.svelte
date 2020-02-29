<script>
	import { createEventDispatcher, onMount } from 'svelte'

	const dispatch = createEventDispatcher()
	const close = () => dispatch('close')

	let modal
  let inputElement

  onMount(() => { inputElement.focus() })

	const handle_keydown = e => {
    switch (e.key) {
      case 'Enter':
        dispatch('add', inputElement.value)
        close()
        break
      case 'Escape':
        close()
        break
      case 'Tab':
        // trap focus
        const nodes = modal.querySelectorAll('*')
        const tabbable = Array.from(nodes).filter(n => n.tabIndex >= 0)

        let index = tabbable.indexOf(document.activeElement)
        if (index === -1 && e.shiftKey) index = 0

        index += tabbable.length + (e.shiftKey ? -1 : 1)
        index %= tabbable.length

        tabbable[index].focus()
        e.preventDefault()
        break
		}
	}
</script>

<svelte:window on:keydown={handle_keydown}/>

<div class="modal-background" on:click={close}></div>

<div class="modal" role="dialog" aria-modal="true" bind:this={modal}>
  <input 
    bind:this={inputElement}
    class="input-url"
    placeholder="Link, RSS"
  />
</div>

<style>
	.modal-background {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0,0,0,0.3);
	}

	.modal {
		position: absolute;
		left: 50%;
		top: 50%;
		width: calc(100vw - 4em);
		max-width: 32em;
		max-height: calc(100vh - 4em);
		overflow: auto;
		transform: translate(-50%,-50%);
		padding: 1em;
		border-radius: 0.2em;
		background: white;
	}

  .input-url {
    font: inherit;
    font-size: 1.5em;
    border: 0;
    width: 100%;
    padding: 0;
  }

  .input-url:focus {
    outline: none;
  }
</style>
