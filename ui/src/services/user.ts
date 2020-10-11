import http from './http'
import { CommnadResponse } from './commandResponse';

export interface LoginResponse {
    token: string;
}

export interface UserAPI {
    login(email: string, password: string): Promise<LoginResponse>;
    updateDefaultsWithAPIToken(token: string): void;
}


async function login(email: string, password: string): Promise<LoginResponse> {
    const res = await http.post(`/c/user/login`, { email, password })
    const response = res.data as CommnadResponse;
    if (response.status !== "accepted") {
        throw new Error(response.reason || 'Failed to login');
    }
    http.defaults.headers['authorization'] = (response.data as LoginResponse).token;
    return response.data as LoginResponse
}

function updateDefaultsWithAPIToken(token: string) {
    http.defaults.headers['authorization'] = token;
}

export const API = {
    login,
    updateDefaultsWithAPIToken
};