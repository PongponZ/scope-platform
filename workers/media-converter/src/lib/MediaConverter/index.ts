import ffmpeg, { FfmpegCommand } from 'fluent-ffmpeg';
const ffmpegInstaller = require('@ffmpeg-installer/ffmpeg');

/**
 * on 'progress': callback
 *  The progress event is emitted every time ffmpeg reports progress information. 
 *  It is emitted with an object argument with the following keys:
 *  "frames": total processed frame count
 *  "currentFps": framerate at which FFmpeg is currently processing
 *  "currentKbps": throughput at which FFmpeg is currently processing
 *  "targetSize": current size of the target file in kilobytes
 *  "timemark": the timestamp of the current frame in seconds
 *  "percent": an estimation of the progress percentage
 */

export class MediaConverter {
    private converter: FfmpegCommand
    private mediaFolder: string
    private convertedFolder: string

    constructor(mediaFolder: string, convertedFolder: string) {
        this.init(mediaFolder, convertedFolder)
    }

    private init(mediaFolder: string, convertedFolder: string) {
        this.mediaFolder = mediaFolder
        this.convertedFolder = convertedFolder

        const defaultTimeout = (60000 * 15) // 15 minute

        this.converter = ffmpeg({ timeout: defaultTimeout }).setFfmpegPath(ffmpegInstaller.path)
        return this
    }



    public convertVideoToSegment(filename, outputName, startCallback, codecDataCallback, processCallback: any, endCallback: any, errCallback: any) {
        this.converter
            .input(`${this.mediaFolder}/${filename}`)
            .addOption([
                '-profile:v baseline',
                '-level 3.0',
                '-start_number 0',
                '-hls_time 10',
                '-hls_list_size 0',
                '-f hls'
            ])
            .output(`${this.convertedFolder}/${outputName}.m3u8`)  //output extension .m3u8
            .on('error', errCallback)
            .on("start", startCallback)
            .on("codecData", codecDataCallback)
            .on("progress", processCallback)
            .on("end", endCallback)
            .run()
    }

}