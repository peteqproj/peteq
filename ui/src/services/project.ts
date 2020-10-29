import http from './http';
import { CommnadResponse } from './commandResponse';

export interface Project {
    metadata: {
        name: string;
        description: string;
        id: string;
        color: string;
        imageUrl: string;
    },
    tasks: string[];
}

export interface CreateProjectRequestBody {
    name:        string;
    description: string;
    color:       string;
    imageUrl:    string;
}

export interface ProjectAPI {
    list(): Promise<Project[]>
    get(id: string): Promise<Project>
    create(req: CreateProjectRequestBody): Promise<Project>
    addTasks(project: string, tasks: string[]): Promise<void>
}


async function list(): Promise<Project[]> {
    const res = await http.get('/api/project')
    return res.data as Project[]
}

async function get(id: string): Promise<Project> {
    const res = await http.get(`/api/project/${id}`)
    return res.data as Project
}

async function create(req: CreateProjectRequestBody): Promise<Project> {
    const res = await http.post('/c/project/create', req);
    const cmdResponse = res.data as CommnadResponse
    if (cmdResponse.reason) {
        throw new Error(`Failed to create project: ${cmdResponse.reason}`)
    }
    return get(cmdResponse.id)
}

async function addTasks(project: string, tasks: string[]): Promise<void> {
    const res = await http.post('/c/project/addTasks', { project, tasks });
    const cmdResponse = res.data as CommnadResponse
    if (cmdResponse.reason) {
        throw new Error(`Failed to add tasks to project: ${cmdResponse.reason}`)
    }
}

export const API = {
    list,
    get,
    create,
    addTasks,
}