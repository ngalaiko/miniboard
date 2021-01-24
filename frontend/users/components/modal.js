(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #background {
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0,0,0,0.3);
        }

        #content {
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

        #input-url {
            font: inherit;
            font-size: 1.5em;
            border: 0;
            width: 100%;
        }

        #input-url:focus {
            outline: none;
        }

        #input-icon {
            cursor: pointer;
        }
    </style>
    <div id="background"></div>
    <div id="content">
        <input id="input-url" placeholder="RSS" />
    </div>
    `

    class AddModal extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        connectedCallback() {
            this.keydownEventHandler = _keydownEventHandler(this)

            window.addEventListener('keydown', this.keydownEventHandler)
            this.shadowRoot.querySelector("#background").addEventListener('click', () => _close(this))
        }

        disconnectedCallback() {
            window.removeEventListener('keydown', this.keydownEventHandler)
        }
    }

    const _keydownEventHandler = (self) => {
        return (e) => {
            switch (e.key) {
            case 'Enter':
                const url = self.shadowRoot.querySelector("#input-url").value
                if (url !== '') self.dispatchEvent(new CustomEvent('FeedAdded', {
                    detail: {
                        url: url,
                    },
                    bubbles: true,
                }))
                _close(self)
                break
            case 'Escape':
                _close(self)
                break
            }
        }
    }

    const _close = (self) => {
        let event = new CustomEvent('Closed', {
            bubbles: true,
        })
        self.dispatchEvent(event)
    }

    customElements.define('add-modal', AddModal)
})()
