import type { User, UserPref } from '../types/User';
import { axiosInstance } from './client';
import type { EditOrDeleteResponse } from './default';

export interface CreateUserRequest {
    userName: string;
    password?: string;
    altName?: string;

    email?: string;
    altEmail?: string;
}

export interface UpdateUserRequest extends Partial<CreateUserRequest> {
    oldPassword?: string;
}

export interface CreateUserResponse {
    id: string;
}

export const createUser = async (user: CreateUserRequest): Promise<CreateUserResponse> => {
    const response = await axiosInstance.post<CreateUserResponse>('/v1/user/', user);
    return response.data;
};

export const getUser = async (userId: string): Promise<User> => {
    const response = await axiosInstance.get<User>(`/v1/user/${userId}`);
    return response.data;
};

export const updateUser = async (user: UpdateUserRequest): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.put<EditOrDeleteResponse>('/v1/user/', user);
    return response.data;
};

export const deleteUser = async (): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.delete<EditOrDeleteResponse>('/v1/user/');
    return response.data;
};

export type CreateUserPref = Pick<UserPref, 'userId' | 'currency' | 'theme'>;

export const createUserPref = async (prefs: CreateUserPref): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.post<EditOrDeleteResponse>('/v1/user/pref', prefs);
    return response.data;
};

export const updateUserPref = async (prefs: Partial<CreateUserPref>): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.put<EditOrDeleteResponse>('/v1/user/pref', prefs);
    return response.data;
};

export const getUserPref = async (): Promise<UserPref> => {
    const response = await axiosInstance.get<UserPref>('/v1/user/pref');
    return response.data;
};

export const updateUserAvatar = async (): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.put<EditOrDeleteResponse>('/v1/user/avatar');
    return response.data;
};

export const uploadUserAvatar = async (file: File) => {
    const formData = new FormData();
    formData.append('avatar', file);

    await axiosInstance.put('/v1/user/avatar', formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};
