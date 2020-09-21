import axios from "axios";
 
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
    const res = await axios.get('http://localhost:8080/api/list')
    return res.data as List[]
}

async function moveTasks(source: string, destination: string, tasks: string[]): Promise<void> {
    await axios.post('http://localhost:8080/api/list/moveTasks', {
        source,
        destination,
        tasks,
    })
}



export const API = {
    list,
    moveTasks,
};