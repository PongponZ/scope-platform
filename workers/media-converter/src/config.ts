/* eslint-disable import/first */
require('dotenv').config()

import { MediaStorageCredential } from './lib/MediaStorage/index';


export interface IConfig {
    RABBITMQ_URI: string,
    RABBITMQ_CONVERT_QUEUE_NAME: string,
    RABBITMQ_CONVERT_ERROR_QUEUE_NAME: string,
    RABBITMQ_CONVERT_COMPLETE_QUEUE_NAME: string,
    MEDIA_STORAGE_CREDENTAIL: MediaStorageCredential,
    MEDIA_BUCKET: string,
    REDIS_HOST:string,
    REDIS_PORT:string
}

export default function GetConfig(): IConfig {
    /** Rabbit Config */
    const host = process.env.RABBITMQ_HOST || null
    const port = process.env.RABBITMQ_PORT || null
    const user = process.env.RABBITMQ_USERNAME || null
    const password = process.env.RABBITMQ_PASSWORD || null
    const convertQueueName = process.env.RABBITMQ_MEDIA_CONVERT_QUEUE_NAME || null
    const errConvertQueueName = process.env.RABBITMQ_MEDIA_CONVERT_ERROR_QUEUE_NAME || null
    const completeQueueName = process.env.RABBITMQ_MEDIA_CONVERT_COMPLETE_QUEUE_NAME || null
    const rabbitmqURI = `amqp://${user}:${password}@${host}:${port}`

    /** Minio Config */
    const minioEndpoint = process.env.MINIO_ENDPOINT || null
    const minioPort = process.env.MINIO_PORT || null
    const minioAccessKey = process.env.MINIO_ACCESS_KEY || null
    const minioSecretKey = process.env.MINIO_SECERT || null
    const minioBucket = process.env.MINIO_BUCKET || null

    const redisHost = process.env.REDIS_HOST || null
    const redisPort = process.env.REDIS_PORT || null

    if (
        !redisHost||
        !redisPort ||
        !host ||
        !port ||
        !user ||
        !password ||
        !convertQueueName ||
        !minioEndpoint ||
        !minioAccessKey ||
        !minioSecretKey ||
        !minioBucket ||
        !minioPort ||
        !errConvertQueueName ||
        !completeQueueName) {

        throw console.error("environment missing");
    }

    return {
        RABBITMQ_URI: rabbitmqURI,
        RABBITMQ_CONVERT_QUEUE_NAME: convertQueueName,
        MEDIA_STORAGE_CREDENTAIL: {
            accessKey: minioAccessKey,
            endPoint: minioEndpoint,
            secretKey: minioSecretKey,
            port: parseInt(minioPort)
        },
        RABBITMQ_CONVERT_ERROR_QUEUE_NAME: errConvertQueueName,
        RABBITMQ_CONVERT_COMPLETE_QUEUE_NAME: completeQueueName,
        MEDIA_BUCKET: minioBucket,
        REDIS_HOST:redisHost,
        REDIS_PORT:redisPort
    }
}
