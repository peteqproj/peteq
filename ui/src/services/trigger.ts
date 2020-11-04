import http from './http';
import { CommnadResponse } from './commandResponse';

export interface Trigger {

}

export interface CreateProjectRequestBody {
    id: string;
}

export interface ProjectAPI {
    create(req: CreateProjectRequestBody): Promise<void>
}



async function create(req: CreateProjectRequestBody): Promise<void> {
    const res = await http.post('/c/project/create', req);
    const cmdResponse = res.data as CommnadResponse
    if (cmdResponse.reason) {
        throw new Error(`Failed to create project: ${cmdResponse.reason}`)
    }
    return Promise.resolve();
}


export const API = {
    create,
}