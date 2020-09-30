import http from './http'
 
export interface List {
    metadata: {
        id: string;
        name: string;
    },
    tasks: string[];
}

export interface ListAPI {
    list(): Promise<List[]>
    moveTasks(source: string, destination: string, tasks: string[]): Promise<void>
}



async function list(): Promise<List[]> {
    const res = await http.get('/api/list')
    return res.data as List[]
}

async function moveTasks(source: string, destination: string, tasks: string[]): Promise<void> {
    await http.post('/api/list/moveTasks', {
        source,
        destination,
        tasks,
    })
}



export const API = {
    list,
    moveTasks,
};