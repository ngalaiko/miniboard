import proto from './sources_service_grpc_web_pb.js'

export class SourcesClient {
    constructor(hostname) {
        this.client = new proto.SourcesServicePromiseClient(hostname)
    }

    async createSource(url) {
        const source = new proto.Source()
        source.setUrl(url)

        const request = new proto.CreateSourceRequest()
        request.setSource(source)
        return await this.client.createSource(request)
    }
}
