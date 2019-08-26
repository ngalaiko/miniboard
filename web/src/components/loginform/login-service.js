export class LoginService {
  constructor(api) {
    this.api = api
  }

  // login returns a Promise with authorization.
  login(username, password) {
    return this.getAuthorization(username, password)
      .then(this.ifNotFound(() => this.signup(username, password) ))
      .catch(this.ifError)
  }

  getAuthorization(username, password) {
    return this.api.postJSON(
     `/api/v1/authorizations`, {
        username: username,
        password: password,
        grant_type: "password"
    })
  }

  ifError(error) {
    if (error.code !== undefined) {
      throw error.message
    }
    console.error(error)
    throw "something went wrong, try again"
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
    return this.api.postJSON(
      `/api/v1/users`, {
        username: username,
        password: password 
       })
      .then(() => { return this.getAuthorization(username, password) })
  }
}
