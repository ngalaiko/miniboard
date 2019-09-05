export class ArticlesService {
  constructor(api, user) {
    this.api  = api
    this.user = user
    this.from = ''
  }

  // next returns a Promise with articles.
  next(pageSize) {
    // if there are no more articles, return en empty list. 
    if (this.from === undefined) {
      return new Promise((resolve, reject) => { resolve([]) })
    }
    return this.api.get(`/api/v1/${this.user.name}/articles?page_size=${pageSize}&page_token=${this.from}`)
      .then(resp => { 
        this.from = resp.next_page_token 
        return resp.articles
      })
  }

  // create returns a Promise with the article.
  add(url) {
    return this.api.post(`/api/v1/${this.user.name}/articles`, {
      url: url
    })
    .then(this.ifError)
  }

  delete(name) {
    return this.api.delete(`/api/v1/${name}`)
      .then(this.ifError)
  }

  ifError(resp) {
    if (resp.error === undefined) {
      return resp
    }
    if (resp.code !== undefined) {
      throw resp.message
    }
    console.error(resp)
    throw "something went wrong, try again"
  }
}
