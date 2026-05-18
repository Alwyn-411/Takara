import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { User } from "../types/User";

type UserState = Partial<Omit<User, "active" | "password">>;

interface UserStore extends UserState {
  updateUser: (data: Partial<UserState>) => void;
  clearUser: () => void;
}

const initialState: UserState = {
  userId: undefined,
  userName: undefined,
  altName: undefined,
  email: undefined,
  altEmail: undefined,
};

export const useUserStore = create<UserStore>()(
  persist(
    (set) => ({
      ...initialState,

      updateUser: (data) => set(data),

      clearUser: () => set(initialState),
    }),
    {
      name: "user-storage",
    },
  ),
);
