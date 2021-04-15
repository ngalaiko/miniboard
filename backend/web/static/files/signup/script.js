const signupButton = document.getElementById('signup-button')
const inputUsername = document.getElementById('username')
const inputPassword = document.getElementById('password')

const handleButtonClick = async (e) => {
    e.preventDefault()

    const response = await fetch('/api/v1/users', {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
            username: inputUsername.value,
            password: inputPassword.value,
        }),
        headers: new Headers({
            "Content-Type": "application/json",
        })
    })

    const body = await response.json()
    if (response.status !== 200) {
        alert(`error: ${body.message}`)
        return
    }

    alert('You are now signed up')
    document.location = '/login'
}

signupButton.addEventListener('click', handleButtonClick)
