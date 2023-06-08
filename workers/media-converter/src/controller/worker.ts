import { MediaConverter } from './../lib/MediaConverter/index';
import { MediaStorage } from './../lib/MediaStorage/index';
import GetConfig, { IConfig } from '../config';
import { MessageBroker } from './../lib/MessageBroker/index';
import { GetMessagePayload, IMessagePayload } from '../enity/message';


import fs from 'fs-extra'
import { Redis } from '../lib/Redis';

export class Worker {
    private mediaFolder: string = "./temp"
    private convertedFolder: string = "./coverted"
    private broker: MessageBroker
    private mediaStorage: MediaStorage
    private converter: MediaConverter
    private redis: Redis
    private config: IConfig
    private totalTime: number = 0
    private payload: IMessagePayload = null
    private ack: any
    constructor() { }

    private async init() {
        await this.initFolder()
        this.config = GetConfig()
        this.broker = await MessageBroker.getInstance(this.config.RABBITMQ_URI)
        this.mediaStorage = await MediaStorage.getInstance(this.config.MEDIA_STORAGE_CREDENTAIL, this.config.MEDIA_BUCKET)
        this.redis = await Redis.getInstance(this.config.REDIS_HOST, this.config.REDIS_PORT)
    }

    private async initFolder() {
        if (!fs.existsSync(this.mediaFolder)) {
            fs.mkdirSync(this.mediaFolder);
        }

        if (!fs.existsSync(this.convertedFolder)) {
            fs.mkdirSync(this.convertedFolder);
        }
    }

    public async run() {
        await this.init()
        this.broker.subscribe(this.config.RABBITMQ_CONVERT_QUEUE_NAME, this.hanlder.bind(this))
    }

    private async hanlder(message, ack) {
        try {
            this.converter = new MediaConverter(this.mediaFolder, this.convertedFolder)

            this.ack = ack;
            this.payload = GetMessagePayload(message)

            const loaded = await this.mediaStorage.getMedia(this.payload.filename, this.payload.path, this.mediaFolder)
            
            if (loaded) {
                this.convertVideo()
            }
        } catch (error) {
            this.broker.send(this.config.RABBITMQ_CONVERT_ERROR_QUEUE_NAME, Buffer.from(JSON.stringify({ error: error })) )
            ack()
        }
    }

    private async  convertVideo() {
        this.converter.convertVideoToSegment(this.payload.filename,
            this.payload.output,
            this.onStart.bind(this),
            this.onCodecData.bind(this),
            this.onProcessing.bind(this),
            this.onEnd.bind(this),
            this.onError.bind(this))
    }

    private onStart() {
        console.log("payload:", this.payload);
    }

    private onCodecData(data) {
        this.totalTime = parseInt(data.duration.replace(/:/g, ''))
    }

    private onProcessing(progress) {
        const time = parseInt(progress.timemark.replace(/:/g, ''))
        const percent = (time / this.totalTime) * 100
        console.log(`processing: ${percent.toFixed(2)}%`,);
        this.statusUpdate("converting", percent.toFixed(2))
    }

    private async onEnd() {
        let retry = 3
     
        const fileList = this.getConvertedFileList()
        for (let i = 0; i < fileList.length; i++) {
            const uploadPath = `${this.convertedFolder}/${fileList[i]}`
            const storePath = `${this.payload.outputPath}/${fileList[i]}`
            const uploaded = await this.mediaStorage.upload(uploadPath, storePath)

            if (!uploaded && retry > 0) {
                console.log(`retry to upload: ${uploadPath}`);
                --i
                --retry
            }

            this.statusUpdate("uploading",  ((i/fileList.length)*100).toString())
        }

        await this.clean()

        this.sendToComplete()
        this.statusUpdate("done", "100")
       
        this.totalTime = 0
        this.ack()
    }

    private onError(err) {
        throw console.error(err);
    }

    private statusUpdate(status: string, percent: string) {
        const processStatus = {
            status: status,
            process: percent
        }
        this.redis.client.setEx(`convert_${this.payload.id}`, 300,JSON.stringify(processStatus))
    }

    private sendToComplete() {
        const completeMessage = {
            id:this.payload.id,
            timestamp:Date.now()
        }
        
        this.broker.send(this.config.RABBITMQ_CONVERT_COMPLETE_QUEUE_NAME, Buffer.from(JSON.stringify(completeMessage)))
    }

    private async clean() {
        try {
            await fs.emptyDir(this.convertedFolder)
            await fs.emptyDir(this.mediaFolder)
            console.log(`file cleaned: ${this.payload.id}\n`);
        } catch (err) {
            console.log(err);
        }
    }

    private getConvertedFileList(): string[] {
        return fs.readdirSync(this.convertedFolder);
    }

}