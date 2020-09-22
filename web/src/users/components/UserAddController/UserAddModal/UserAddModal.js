(async () => {
    const res = await fetch('/users/components/UserAddController/UserAddModal/UserAddModal.html')
    const textTemplate = await res.text()

    const HTMLTemplate = new DOMParser().parseFromString(textTemplate, 'text/html')
                            .querySelector('template')

    class UserAddModal extends HTMLElement {
        constructor() { 
             super()
        }

        connectedCallback() {
            const shadowRoot = this.attachShadow({ mode: 'open' })

            const instance = HTMLTemplate.content.cloneNode(true)
            shadowRoot.appendChild(instance)

            this.keydownEventHandler = _keydownEventHandler(this)

            window.addEventListener('keydown', this.keydownEventHandler)
            shadowRoot.querySelector(".modal__input-file").addEventListener('change', _fileInputEventHandler(this))
            shadowRoot.querySelector(".modal__background").addEventListener('click', () => _close(this))
        }

        disconnectedCallback() {
            window.removeEventListener('keydown', this.keydownEventHandler)
        }
    }

    const _keydownEventHandler = (self) => {
        return (e) => {
            switch (e.key) {
            case 'Enter':
                const url = self.shadowRoot.querySelector(".modal__input-url").value
                if (url !== '') self.dispatchEvent(new CustomEvent('UrlAdded', {
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

    const _fileInputEventHandler = (self) => {
        return (e) => {
            const files = e.target.files

            if (files.length === 0) return

            self.dispatchEvent(new CustomEvent('FileAdded', {
                detail: {
                    file: files[0],
                },
                bubbles: true,
            }))

            _close(self)
        }
    }

    const _close = (self) => {
        let event = new CustomEvent('Closed', {
            bubbles: true,
        })
        self.dispatchEvent(event)
    }

    customElements.define('user-add-modal', UserAddModal)
})()
