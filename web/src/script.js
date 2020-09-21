const emailButton = document.getElementById('email-button')
const emailInput = document.getElementById('email-input')

const codeButton = document.getElementById('code-button')
const codeInput = document.getElementById('code-input')
const codeForm = document.getElementById('code-form')

const handleEmailClick = async (e) => {
    e.preventDefault()

    let response = await fetch("/api/v1/codes", {
        method: 'POST',
        body: JSON.stringify({
            email: emailInput.value
        })
    })

    if (response.status == 200) {
        codeForm.style.visibility = 'visible'
    }
}

const handleEmailInput = async (e) => {
    e.preventDefault()

    codeForm.style.visibility = 'hidden'
}

const handleCodeClick = async (e) => {
    e.preventDefault()

    document.location = `/codes?code=${codeInput.value}`
}

emailButton.addEventListener('click', handleEmailClick)
codeButton.addEventListener('click', handleCodeClick)

emailInput.addEventListener('input', handleEmailInput)
