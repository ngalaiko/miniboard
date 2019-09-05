export class LabelsService {
  constructor(api) {
    this.api = api
    this.titleToLabel = {}
    this.nameToLabel = {}
    this.titles = []

    // todo: list all in smaller batches
    this.api.get(`/api/v1/${this.api.subject()}/labels?page_size=100`)
      .then(resp => resp.labels.forEach(label => {
        this.saveLabel(label)
      }))
  }

  create(title) {
    let known = this.titleToLabel[title]
    if (known !== undefined) {
      return new Promise(resolve => resolve(known))
    }
    return this.api.post(`/api/v1/${this.api.subject()}/labels`, {
      title: title,
    }).then(label => {
      this.saveLabel(label)
      return label
    })
  }

  get(labelName) {
    let known = this.nameToLabel[labelName]
    if (known !== undefined) {
      return new Promise(resolve => resolve(known))
    }
    return this.api.get(`/api/v1/${labelName}`)
      .then(label => {
        this.saveLabel(label)
        return label
      })
  }

  saveLabel(label) {
    if (this.titleToLabel[label.title] !== undefined) {
      return
    }

    this.titleToLabel[label.title] = label
    this.nameToLabel[label.name] = label
    this.titles.push(label.title)
  }
}
