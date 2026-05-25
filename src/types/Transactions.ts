import type { Base, Timestamp } from './Common';

const TransactionDirection = { Debit: 'Debit', Credit: 'Credit' } as const;
export type TransactionType = (typeof TransactionDirection)[keyof typeof TransactionDirection];

export interface Transaction extends Base, Timestamp {
    accountId: string;
    transactionId: string;
    type: TransactionType;

    settledAmount: string;
    settledCurrency: string;
    accountAmount: string;
    accountCurrency: string;

    exchangeRate: string;
    merchant: string;
    category: string;
    description: string;
    TransactionAt: number;
}
