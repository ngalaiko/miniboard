(async () => {
    const HTMLTemplate = document.createElement('template')
    HTMLTemplate.innerHTML = `
    <style>
        #toast-container {
            visibility: hidden;

            display: flex;
        }

        #toast-message {
            display: flex;
            align-items: center;

            min-width: 100px;
            max-width: 100%;

            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
        }

        #toast-button {
            display: flex;
            align-items: center;

            background: inherit;
            font-size: inherit;

            cursor: pointer;

            margin: 0;
            padding: 0;
            border: 0;
        }

        #toast-button:active {
            color: inherit;
        }

        #toast-container.show {
            visibility: visible;
        }
    </style>
    <span id="toast-container">
        <span id="toast-message"></span>
        <button id="toast-button" type="button">
            <svg viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
        </button>
    </span>
    `

    class Toast extends HTMLElement {
        constructor() { 
            super()

            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)
        }

        connectedCallback() {
            this.show()

            this.shadowRoot.querySelector('#toast-button').addEventListener('click', () => this.closeNow())
        }

        static get observedAttributes() {
            return ['message']
        }

        attributeChangedCallback(attribute, oldValue, newValue) {
            switch (attribute) {
            case 'message':
                this.message = newValue
                break
            }
        }

        show() {
            this.shadowRoot.querySelector('#toast-container').classList.add('show')
        }
        
        closeNow() {
            this.shadowRoot.querySelector('#toast-container').classList.remove('show')
            this.remove()
        }

        close() {
            setTimeout(() => this.closeNow(), 3000)
        }

        set message(value) {
            if (!this.shadowRoot) return

            this.shadowRoot.querySelector('#toast-message').innerText = value
        }
    }

    customElements.define('x-toast', Toast)
})()
