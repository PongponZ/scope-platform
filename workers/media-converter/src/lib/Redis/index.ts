import { createClient, RedisClientType } from "redis"

export class Redis {
    private static instance: Redis
    public client: RedisClientType
    private constructor() { }

    
    public async init(host: string, port: string) {
        this.client = createClient({
            socket:{
                host:host,
                port:parseInt(port) 
            }
        })
        await this.client.connect()
        await this.client.ping()

        console.log(`redis connected: ${host}:${port}`)

        return this
    }

    public static async getInstance(host: string, port: string): Promise<Redis> {
        if (!Redis.instance) {
            const broker = new Redis()
            Redis.instance = await broker.init(host, port)

        }
        return Redis.instance
    }
}