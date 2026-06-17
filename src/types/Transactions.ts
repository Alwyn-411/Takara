import type { Base, Timestamp } from './Common';

export const TransactionType = { Debit: 'Debit', Credit: 'Credit' } as const;
export type TransactionDataType = (typeof TransactionType)[keyof typeof TransactionType];

export const HoldingType = { Asset: 'Asset', Liability: 'Liability' } as const;
export type HoldingDataType = (typeof HoldingType)[keyof typeof HoldingType];

export interface Transaction extends Base, Timestamp {
    accountId: string;
    transactionId: string;
    type: TransactionDataType;

    settledAmount: string;
    settledCurrency: string;
    accountAmount: string;
    accountCurrency: string;

    exchangeRate: string;
    merchantName: string;
    categoryName: string;
    tags: string[];
    description: string;
    transactionAt: number;
}

export interface Merchant extends Base, Timestamp {
    merchantId: string;
    merchantName: string;
}

export interface Category extends Base, Timestamp {
    categoryId: string;
    categoryName: string;
}

export interface Tag extends Base, Timestamp {
    tagId: string;
    tagName: string;
}
