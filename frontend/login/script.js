const loginButton = document.getElementById('login-button')
const inputUsername = document.getElementById('username')
const inputPassword = document.getElementById('password')

const handleButtonClick = (e) => {
    e.preventDefault()

    console.log(inputUsername.value, inputPassword.value)
}

loginButton.addEventListener('click', handleButtonClick)
