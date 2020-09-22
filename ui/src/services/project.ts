import http from './http';
import { CommnadResponse } from './commandResponse';

export interface Project {
    metadata: {
        name: string;
        description: string;
        id: string;
    },
    tasks: string[];
}

export interface ProjectAPI {
    list(): Promise<Project[]>
    get(id: string): Promise<Project>
    create(project: Project): Promise<Project>
    addTasks(project: string, tasks: string[]): Promise<void>
}


async function list(): Promise<Project[]> {
    const res = await http.get('/project')
    return res.data as Project[]
}

async function get(id: string): Promise<Project> {
    const res = await http.get(`/project/${id}`)
    return res.data as Project
}

async function create(project: Project): Promise<Project> {
    const res = await http.post('/project/create', project);
    const cmdResponse = res.data as CommnadResponse
    if (cmdResponse.reason) {
        throw new Error(`Failed to create project: ${cmdResponse.reason}`)
    }
    return get(cmdResponse.id)
}

async function addTasks(project: string, tasks: string[]): Promise<void> {
    const res = await http.post('/project/addTasks', { project, tasks });
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