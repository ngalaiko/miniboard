import proto from './users_service_grpc_web_pb.js'

export class UsersClient {
    constructor(hostname) {
        this.client = new proto.UsersServicePromiseClient(hostname)
    }

    async me() {
        const request = new proto.GetMeRequest()
        return await this.client.getMe(request)
    }
}
