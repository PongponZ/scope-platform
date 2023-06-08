/* eslint-disable import/first */
require('dotenv').config()

import * as Minio from 'minio';

export interface MediaStorageCredential {
    endPoint: string,
    accessKey: string,
    secretKey: string,
    port: number
}


export class MediaStorage {
    private static instance: MediaStorage
    private bucket: string
    private storageClient: Minio.Client
    private constructor() { }

    private async init(credential: MediaStorageCredential) {
        this.storageClient = new Minio.Client({
            endPoint: credential.endPoint,
            port: credential.port,
            accessKey: credential.accessKey,
            secretKey: credential.secretKey,
            useSSL: false
        })

        console.log(`Minio Connected: ${credential.endPoint}:${credential.port}`);

        return this
    }

    private async initBucket(buckget: string) {
        this.storageClient.bucketExists(buckget, (err, exist) => {
            if (err) {
                console.log(err);
            }

            if (!exist) {
                this.storageClient.makeBucket(buckget, "", (err) => {
                    console.log(err)
                })
            }
        })
    }

    public static async getInstance(credential: MediaStorageCredential, buckget: string): Promise<MediaStorage> {
        if (!MediaStorage.instance) {
            const storage = new MediaStorage()
            storage.bucket = buckget
            MediaStorage.instance = await storage.init(credential)
            await storage.initBucket(buckget)
        }
        return MediaStorage.instance
    }


    public async getMedia(filname: string, path: string, saveFolder: string): Promise<boolean> {
        return new Promise((reslove, reject) => {
            this.storageClient.fGetObject(this.bucket, path, `${saveFolder}/${filname}`, function (err) {
                if (err) {
                    reject(err)
                }
                reslove(true)
            })
        })
    }

    public async upload(path: string, uploadPath: string): Promise<boolean> {
        return new Promise((reslove, reject) => {
            this.storageClient.fPutObject(this.bucket, uploadPath, path, {}, function (err, etag) {
                if (err) {
                    console.log(err)
                    return reject(false)
                }
                console.log(`${path} : uploaded successfully.`)
                reslove(true)
            });
        })
    }
}