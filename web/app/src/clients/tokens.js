import proto from './tokens_service_grpc_web_pb.js'

export class TokensClient {
    constructor(hostname) {
        this.client = new proto.TokensServicePromiseClient(hostname)
    }

    async exchangeCode(code) {
        const request = new proto.CreateTokenRequest()
        request.setCode(code)
        return await this.client.createToken(request)
    }
}
