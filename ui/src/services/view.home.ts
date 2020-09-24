import http from './http';
import { Task } from './tasks';
import { Project } from './project';

export interface HomeViewModel {
    lists: HomeViewList[];
}

interface HomeViewList {
    metadata: {
        id: string;
        name: string;
    },
    tasks: HomeViewTask[];
}

interface HomeViewTask extends Task {
    project?: Project;
}

export interface HomeViewAPI {
    get(): Promise<HomeViewModel>
}

async function get(): Promise<HomeViewModel> {
    const res = await http.get('/view/home')
    return res.data as HomeViewModel
}

export const API = {
    get,
};
