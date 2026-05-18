import { axiosInstance } from "./client";

export interface CreateUserRequest {
  userName: string;
  password?: string;
  altName?: string;

  email?: string;
  altEmail?: string;
}

export interface CreateUserResponse {
  id: string;
}

export const createUser = async (
  user: CreateUserRequest,
): Promise<CreateUserResponse> => {
  const response = await axiosInstance.post<CreateUserResponse>(
    "/v1/user/",
    user,
  );

  return response.data;
};
