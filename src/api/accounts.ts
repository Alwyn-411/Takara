import type { Accounts } from '../types/Accounts';
import { axiosInstance } from './client';
import type { CreateResponse, EditOrDeleteResponse, ListResponse } from './default';

export const createAccount = async (account: Partial<Accounts>): Promise<CreateResponse> => {
    const response = await axiosInstance.post<CreateResponse>('/v1/account/', account);
    return response.data;
};

export const getAccountWithAccountId = async (accountId: string): Promise<Accounts> => {
    const response = await axiosInstance.get<Accounts>(`/v1/account/${accountId}`);
    return response.data;
};

export const editAccountWithAccountId = async (params: { accountId: string; account: Partial<Accounts> }): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.put<EditOrDeleteResponse>(`/v1/account/${params.accountId}`, params.account);
    return response.data;
};

export const deleteAccountWithAccountId = async (accountId: string): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.delete<EditOrDeleteResponse>(`/v1/account/${accountId}`);
    return response.data;
};

export const listAccountsWithUserId = async (userId: string): Promise<ListResponse<Accounts>> => {
    const response = await axiosInstance.get<ListResponse<Accounts>>(`/v1/account/user/${userId}`);
    return response.data;
};
