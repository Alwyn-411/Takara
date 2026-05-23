import type { Base, Timestamp } from './Common';

export interface Accounts extends Base, Timestamp {
    accountId: string;
    type: string;
    name: string;
    accountNumber: string;
    description?: string;
    currency: string;
    interest: string;
    balance: string;
}

export const currencies = [
    {
        label: '₹ INR - Indian Rupee',
        value: 'INR',
        symbol: '₹',
    },
    {
        label: '$ USD - US Dollar',
        value: 'USD',
        symbol: '$',
    },
];
