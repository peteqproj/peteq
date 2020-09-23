import axios from "axios";
import { CommnadResponse } from './commandResponse';


export interface Task {
    metadata: {
        id: string;
        name: string;
        description: string;
    },
    spec: any,
    status: {
        completed: boolean
    }
}

export interface TaskAPI {
    create(task: Task): Promise<Task>
    list(): Promise<Task[]>
    get(id: string): Promise<Task>
    remove(id: string): Promise<void>
    update(task: Task): Promise<Task>
    complete(task: string): Promise<void>
    reopen(task: string): Promise<void>
}


async function get(id: string): Promise<Task> {
    const res = await axios.get(`http://localhost:8080/api/task/${id}`)
    return res.data as Task
}

async function list(): Promise<Task[]> {
    const res = await axios.get('http://localhost:8080/api/task')
    return res.data as Task[]
}

async function remove(id: string): Promise<void> {
    await axios.post(`http://localhost:8080/api/task/delete`, { id })
}

async function create(task: Task): Promise<Task> {
    const res = await axios.post('http://localhost:8080/api/task/create', task);
    const cmdResponse = res.data as CommnadResponse
    if (cmdResponse.reason) {
        throw new Error(`Failed to create task: ${cmdResponse.reason}`)
    }
    return get(cmdResponse.id)
}

async function update(task: Task): Promise<Task> {
    const res = await axios.post(`http://localhost:8080/api/task/update`, task);
    const cmdResponse = res.data as CommnadResponse
    if (cmdResponse.reason) {
        throw new Error(`Failed to update task: ${cmdResponse.reason}`)
    }
    return get(cmdResponse.id)
}

async function complete(task: string): Promise<void> {
    await axios.post(`http://localhost:8080/api/task/complete`, {task});
}
async function reopen(task: string): Promise<void> {
    await axios.post(`http://localhost:8080/api/task/reopen`, {task});
}

export const API = {
    get,
    list,
    remove,
    create,
    update,
    complete,
    reopen,
};