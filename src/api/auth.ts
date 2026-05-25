import { axiosInstance } from "./client";

export interface AuthUserRequest {
  userName: string;
  password: string;
  remember: boolean;
}

export interface AuthUserResponse {
  id: string;
}

export const AuthUser = async (
  props: AuthUserRequest,
): Promise<AuthUserResponse> => {
  const response = await axiosInstance.get<AuthUserResponse>("/v1/auth", {
    params: {
      userName: props.userName,
      password: props.password,
    },
  });

  return response.data;
};
