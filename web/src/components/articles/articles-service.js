export class ArticlesService {
  constructor(api) {
    this.api  = api
    this.from = ''
  }

  // next returns a Promise with articles.
  next(pageSize) {
    // if there are no more articles, return en empty list. 
    if (this.from === undefined) {
      return new Promise((resolve, reject) => { resolve([]) })
    }
    return this.api.get(`/api/v1/${this.api.subject()}/articles?page_size=${pageSize}&page_token=${this.from}`)
      .then(resp => { 
        this.from = resp.next_page_token 
        return resp.articles
      })
  }

  // create returns a Promise with the article.
  add(url) {
    return this.api.post(`/api/v1/${this.api.subject()}/articles`, {
      url: url
    })
    .then(this.ifError)
  }

  delete(name) {
    return this.api.delete(`/api/v1/${name}`)
      .then(this.ifError)
  }

  updateLabels(article) {
    return this.api.patch(`/api/v1/${article.name}?update_mask=label_ids`, {
      label_ids: article.label_ids,
    }).then(this.ifError)
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
