export default function labels(api) {
    let $ = {}

    let titleToLabel = {}
    let nameToLabel = {}
    $.titles = []

    $.create = async (title) => {
        let known = titleToLabel[title]
        if (known !== undefined) {
            return new Promise(resolve => resolve(known))
        }
        let label = await api.post(`/api/v1/${api.subject()}/labels`, {
            title: title,
        })
        saveLabel(label)
        return label
    }

    $.get = async (labelName) => {
        let known = nameToLabel[labelName]
        if (known !== undefined) {
            return known
        }
        let label = await api.get(`/api/v1/${labelName}`)
        saveLabel(label)
        return label
    }

    let list = async () => {
        // todo: list all in smaller batches
        let resp = await api.get(`/api/v1/${api.subject()}/labels?page_size=100`)

        resp.labels.forEach(label => {
            saveLabel(label)
        })
    }
    list()

    let saveLabel = (label) => {
        if (titleToLabel[label.title] !== undefined) {
            return
        }

        titleToLabel[label.title] = label
        nameToLabel[label.name] = label
        $.titles.push(label.title)
    }

    return $
}
