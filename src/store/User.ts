import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '../types/User';

type UserState = Partial<Omit<User, 'active' | 'password'>>;

interface UserStore extends UserState {
    token: string | undefined;
    setToken: (token: string) => void;
    updateUser: (data: Partial<UserState>) => void;
    clearUser: () => void;
}

const initialState = {
    userId: undefined,
    userName: undefined,
    altName: undefined,
    email: undefined,
    altEmail: undefined,
    token: undefined,
};

export const useUserStore = create<UserStore>()(
    persist(
        (set) => ({
            ...initialState,
            setToken: (token) => set({ token }),
            updateUser: (data) => set(data),
            clearUser: () => set(initialState),
        }),
        {
            name: 'user-storage',
        },
    ),
);
