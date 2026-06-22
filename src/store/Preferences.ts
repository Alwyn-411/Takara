import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { UserPref } from '../types/User';

type UserPrefState = Pick<UserPref, 'currency' | 'theme'>;

interface UserPrefStore extends UserPrefState {
    updateUserPref: (data: Partial<UserPrefState>) => void;
    clearUserPref: () => void;
}

export const userPrefInitialState: UserPrefState = {
    currency: 'USD',
    theme: 'light',
};

export const useUserPrefStore = create<UserPrefStore>()(
    persist(
        (set) => ({
            ...userPrefInitialState,
            updateUserPref: (data: Partial<UserPrefState>) => set(data),
            clearUserPref: () => set(userPrefInitialState),
        }),
        {
            name: 'user-pref-storage',
        },
    ),
);
