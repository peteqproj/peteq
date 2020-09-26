import http from '../http';
import { Task } from '../tasks';

export interface ProjectView {
    metadata: {
        name: string;
        description: string;
        id: string;
        color: string;
        imageUrl: string;
    },
    tasks: Task[];
}

export interface ProjectViewAPI {
    get(id: string): Promise<ProjectView>
}

async function get(id: string): Promise<ProjectView> {
    const res = await http.get(`/view/projects/${id}`)
    return res.data as ProjectView
}

export const API = {
    get,
};
