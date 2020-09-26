import http from '../http';
import { Project } from '../project';
import { Task } from '../tasks';

export interface ProjectView extends Project {
    objects: Task[];
}

export interface ProjectsViewModel {
    projects: ProjectView[];
}

export interface ProjectsViewAPI {
    get(): Promise<ProjectsViewModel>
}

async function get(): Promise<ProjectsViewModel> {
    const res = await http.get('/view/projects')
    return res.data as ProjectsViewModel
}

export const API = {
    get,
};
