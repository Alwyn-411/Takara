import type { Accounts } from "../types/Accounts";
import { axiosInstance } from "./client";

export interface CreateAccountResponse {
  id: string;
}

export const createAccount = async (
  account: Partial<Accounts>,
): Promise<CreateAccountResponse> => {
  const response = await axiosInstance.post<CreateAccountResponse>(
    "/v1/account/",
    account,
  );

  return response.data;
};

export const getAccountWithAccountId = async (accountId: string): Promise<Accounts> => {
  const response = await axiosInstance.get<Accounts>(`/v1/account/${accountId}`);
  return response.data;
};

export const listAccountsWithUserId = async (userId: string): Promise<Accounts[]> => {
  const response = await axiosInstance.get<Accounts[]>(`/v1/account/${userId}`);
  return response.data;
};
