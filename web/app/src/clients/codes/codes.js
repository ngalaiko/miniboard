import proto from './codes_service_grpc_web_pb.js'

export class CodesClient {
    constructor(hostname) {
        this.client = new proto.CodesServicePromiseClient(hostname)
    }

    async sendCode(email) {
        const request = new proto.CreateCodeRequest()
        request.setEmail(email)
        return await this.client.createCode(request)
    }
}
