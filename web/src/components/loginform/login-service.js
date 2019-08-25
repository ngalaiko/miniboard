export class LoginService {
  constructor(api) {
    this.api = api
  }

  // login returns a Promise with authorization.
  login(username, password) {
    return this.getAuthorization(username, password)
      .then(this.ifNotFound(() => { return this.signup(username, password) }))
      .then(this.handleLoggedIn)
      .catch(this.ifError)
  }

  getAuthorization(username, password) {
    return this.api.post(
     `/api/v1/users/${username}/authorizations`, {
      password: password 
    })
  }

  handleLoggedIn(data) {
    console.log("logged in", data)
    // todo: create and return authorization
  }

  ifError(error) {
    if (error.message !== undefined) {
      throw Error(`api error: ${error.message}`)
      // todo: show error message
    }
    throw error
  }

  ifNotFound(then) {
    return function(response) {
      if (response.error === undefined) {
        return response
      }
      if (response.code === 5) { // Not found
        return then()
      } 
      throw response
    }
  }

  signup(username, password) {
    return this.api.post(
      `/api/v1/users`, {
        username: username,
        password: password 
       })
      .then(() => { return this.getAuthorization(username, password) })
  }
}
