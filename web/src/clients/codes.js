import proto from './proto/codes_service_grpc_web_pb.js'

export class Code {
    constructor(protoCode) {
        this.proto = protoCode
    }
}

export class CodesClient {
    constructor(hostname) {
        this.client = new proto.CodesServicePromiseClient(hostname)
    }

    async sendCode(email) {
        const request = new proto.CreateCodeRequest()
        request.setEmail(email)
        return new Code(await this.client.createCode(request))
    }
}
