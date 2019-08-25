export class Api {
  post(url, data) {
    console.log("sending", url)
    return fetch(url, {
        method: "POST",
        body: JSON.stringify(data),
    })
    .then(response => response.json())
  }
}
