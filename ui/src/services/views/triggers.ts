import http from '../http';

export interface TriggersViewModel {
    triggers: TriggerView[];
}

export interface TriggerView {
    name: string;
    description: string;
    type: string;
    spec: any;
}

export interface TriggersViewAPI {
    get(): Promise<TriggersViewModel>
}

async function get(): Promise<TriggersViewModel> {
    const res = await http.get('/q/triggers')
    return res.data as TriggersViewModel
}

export const API = {
    get,
};
