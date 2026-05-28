export interface CreateResponse {
    id: string;
}

export interface EditOrDeleteResponse {
    message: 'success';
}

export interface ListResponse<T> {
    count: number;
    records: T[];
}
