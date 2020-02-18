export class ApiClient {
    get(url: string): Promise<any> {
        return this.send(url, 'GET')
    }

    post(url: string, data: any): Promise<any> {
        return this.send(url, 'POST', data)		
    }

    patch(url: string, data: any): Promise<any> {
        return this.send(url, 'PATCH', data)
    }

    delete(url: string): Promise<any> {
        return this.send(url, 'DELETE')
    }

    async send(url: string, method: string, body?: any): Promise<any> {
        const options: RequestInit = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
        }

        if (body) options.body = JSON.stringify(body)
        
        const resp: Response = await fetch(url, options)
        
        if (resp.status / 100 === 2) {
            return resp.json()
        }

        throw `${method} ${url} - status code: ${resp.status}`		
    }
}
