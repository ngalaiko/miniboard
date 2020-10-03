class FeedService {
    cache = []

    async get(id) {
        if (this.cache[id] !== undefined) return this.cache[id]

        const response = await fetch(`/api/v1/feeds/${id}`)
        if (response.status !== 200) {
            throw `failed to fetch feed: ${(await response.json()).message}`
        }

        const body = await response.json()

        this.cache[body.id] = body

        return body
    }

    async list(size, pageToken) {
        if (pageToken === undefined) pageToken = ''

        const response = await fetch(`/api/v1/feeds?page_size=${size}&page_token=${pageToken}`)
        if (response.status !== 200) {
            throw `failed to fetch feeds: ${(await response.json()).message}`
        }

        const body = await response.json()

        body.feeds.forEach(feed => this.cache[feed.id] = feed)

        return body
    }
}

export default new FeedService()
