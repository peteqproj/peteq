export interface CommnadResponse {
    status: 'accepted' | 'rejected';
    reason?: string;
    type: string;
    id: string;
}