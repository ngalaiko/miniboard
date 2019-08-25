export class Api {
  post(url, data) {
    return fetch(url, {
        method: "POST",
        body: JSON.stringify(data),
    })
    .then(function(res){ return res.json(); })
  }
}
