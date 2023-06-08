import { Worker } from "./controller/worker"

(async () => {
    const wroker = new Worker()

    wroker.run()
    
})().catch(err => console.log(err))

