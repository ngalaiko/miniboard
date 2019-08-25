export class LoginService {
  constructor(api) {
    this.api = api
  }

  login(username, password) {
    this.api.post(
      `/api/v1/users/${username}/authorizations`, {
        password: password 
      })
    .then(this.handleLogin)
  }

  handleLogin(response) {
    console.error(response.error)
  }

}
