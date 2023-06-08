/* eslint-disable import/first */
require('dotenv').config()

import amqp, { Channel, Connection } from 'amqplib'
import _ from 'lodash'


export class MessageBroker {
    private queues: any
    private connection!: Connection
    private channel!: Channel
    private static instance: MessageBroker
    private constructor() { }
    public uri: string

    public async init(uri: string) {
        this.uri = uri
        this.connection = await amqp.connect(uri)
        this.channel = await this.connection.createChannel()

        console.log(`Rabbit Connected: ${uri}`)

        return this
    }


    public static async getInstance(uri: string): Promise<MessageBroker> {
        if (!MessageBroker.instance) {
            const broker = new MessageBroker()
            MessageBroker.instance = await broker.init(uri)

        }
        return MessageBroker.instance
    }

    public async send(queue: string, msg: any) {
        if (!this.connection) {
            await this.init(this.uri)
        }
        await this.channel.assertQueue(queue, { durable: true })
        this.channel.sendToQueue(queue, msg)
    }

    public async subscribe(queue: string, handler: any) {
        if (!this.connection) {
            await this.init(this.uri)
        }

        console.log(`Queue Subscribed: ${queue}`);

        this.queues = { queue: handler, ...this.queues }

        if (this.queues[queue]) {
            const existingHandler = _.find(this.queues[queue], h => h === handler)

            if (existingHandler) {
                return () => this.unsubscribe(queue, existingHandler)
            }

            this.queues[queue].push(handler)
            return () => this.unsubscribe(queue, handler)
        }

        await this.channel.assertQueue(queue, { durable: true })
        this.queues[queue] = [handler]
        this.channel.consume(
            queue,
            async (msg: any) => {
                const ack = _.once(() => this.channel.ack(msg))
                this.queues[queue].forEach(h => h(msg, ack))
            }
        )

        return () => this.unsubscribe(queue, handler)
    }

    public async unsubscribe(queue: any, handler: any) {
        _.pull(this.queues[queue], handler)
    }
}
