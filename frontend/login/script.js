import AuthorizationsService from '/users/services/authorizations.js'

const loginButton = document.getElementById('login-button')
const inputUsername = document.getElementById('username')
const inputPassword = document.getElementById('password')

const handleButtonClick = async (e) => {
    e.preventDefault()

    AuthorizationsService.create({
        username: inputUsername.value,
        password: inputPassword.value,
    }).then(() => {
        document.location = '/users'
    }).catch((error) => {
        alert(`error: ${error}`)
    })
}

loginButton.addEventListener('click', handleButtonClick)
