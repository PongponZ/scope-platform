export interface IMessagePayload {
    id: string,
    filename: string,
    output: string,
    outputPath:string,
    path: string,
    type: string
}

export function GetMessagePayload(payload: any): IMessagePayload {
    return JSON.parse(payload?.content.toString() || payload)
}
