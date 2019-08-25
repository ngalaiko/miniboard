export class Api {
  constructor() {
    this.authentication = null
  }

  post(url, data) {
    return fetch(url, {
        method: "POST",
        body: JSON.stringify(data),
    })
    .then(response => response.json())
  }

  authenticate(auth) {
    this.authentication = auth
  }
}
