import type { Transaction } from '../types/Transactions';
import { axiosInstance } from './client';
import type { CreateResponse, EditOrDeleteResponse, ListResponse } from './default';

export const createTransaction = async (transaction: Partial<Transaction>): Promise<CreateResponse> => {
    const response = await axiosInstance.post<CreateResponse>('/v1/transaction/', transaction);
    return response.data;
};

export const getTransactionByTransactionId = async (transactionId: string): Promise<Transaction> => {
    const response = await axiosInstance.get<Transaction>(`v1/transaction/${transactionId}`);
    return response.data;
};

export const editTransactionWithTransactionId = async (params: {
    transactionId: string;
    transaction: Partial<Transaction>;
}): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.put<EditOrDeleteResponse>(`v1/transaction/${params.transactionId}`, params.transaction);
    return response.data;
};

export const deleteTransactionWithTransactionId = async (transactionId: string): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.delete<EditOrDeleteResponse>(`/v1/transaction/${transactionId}`);
    return response.data;
};

export const listAccountsWithUserId = async (accountId: string): Promise<ListResponse<Transaction>> => {
    const response = await axiosInstance.get<ListResponse<Transaction>>(`/v1/transaction/account/${accountId}`);
    return response.data;
};
