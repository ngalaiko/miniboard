<script>
	import { createEventDispatcher, onMount } from 'svelte'
  import { UploadIcon } from 'svelte-feather-icons'

	const dispatch = createEventDispatcher()
	const close = () => dispatch('close')

	let modal
  let inputElement

  onMount(() => { inputElement.focus() })

  const onFile = e => {
    const files = e.target.files

    if (files.length === 0) return

    dispatch('file', files[0])
    close()
  }

	const handleKeydown = e => {
    switch (e.key) {
      case 'Enter':
        dispatch('link', inputElement.value)
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

<svelte:window on:keydown={handleKeydown}/>

<div class="modal-background" on:click={close}></div>

<div class="modal" role="dialog" aria-modal="true" bind:this={modal}>
  <input 
    bind:this={inputElement}
    class="input-url"
    placeholder="Link, RSS"
  />
  <input type="file" id="file" on:change={onFile} />
  <label for="file" class="input-file">
    <UploadIcon size="30" />
  </label>
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
    display: inline-flex;
		position: absolute;
		left: 50%;
		top: 50%;
		width: calc(100vw - 4em);
		max-width: 32em;
		max-height: calc(100vh - 4em);
		overflow: auto;
		transform: translate(-50%,-50%);
		padding: 1em;
		border-radius: 0.3em;
		background: white;
	}

  input[type="file"] {
    display: none;
  }

  .input-file {
    cursor: pointer;
  }

  .input-url {
    font: inherit;
    font-size: 1.5em;
    border: 0;
    width: 100%;
    padding: 0;
    margin: 0;
  }

  .input-url:focus {
    outline: none;
  }
</style>
