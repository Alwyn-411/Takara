import type { Category, Merchant, Tag, Transaction } from '../types/Transactions';
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

export const listTransactionsByAccountId = async (accountId: string): Promise<ListResponse<Transaction>> => {
    const response = await axiosInstance.get<ListResponse<Transaction>>(`/v1/transaction/account/${accountId}`);
    return response.data;
};

export const listMerchants = async (q: string): Promise<ListResponse<Merchant>> => {
    const response = await axiosInstance.get<ListResponse<Merchant>>('/v1/merchants/list', {
        params: {
            q: q,
        },
    });
    return response.data;
};

export const listCategories = async (q: string): Promise<ListResponse<Category>> => {
    const response = await axiosInstance.get<ListResponse<Category>>('/v1/category/list', {
        params: {
            q: q,
        },
    });
    return response.data;
};

export const listTag = async (q: string): Promise<ListResponse<Tag>> => {
    const response = await axiosInstance.get<ListResponse<Tag>>('/v1/tag/list', {
        params: {
            q: q,
        },
    });
    return response.data;
};
