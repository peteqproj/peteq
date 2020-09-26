import http from '../http';
import { Task } from '../tasks';

export interface BacklogTask extends Task {
    list: BacklogList;
    project: BacklogProject;
}

export interface BacklogList {
    name: string;
    id: string;
}

export interface BacklogProject {
    name: string;
    id: string;
}

export interface BacklogViewModel {
    tasks: BacklogTask[];
    lists: BacklogList[];
    projects: BacklogProject[];
}

export interface BacklogViewAPI {
    get(): Promise<BacklogViewModel>
}

async function get(): Promise<BacklogViewModel> {
    const res = await http.get('/view/backlog')
    return res.data as BacklogViewModel
}

export const API = {
    get,
};
