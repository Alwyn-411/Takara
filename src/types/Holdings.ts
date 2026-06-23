import type { ListResponse } from '../api/default';

export const HoldingType = { Asset: 'Asset', Liability: 'Liability' } as const;
export type HoldingDataType = (typeof HoldingType)[keyof typeof HoldingType];

type Timestamp = number;

export interface Holding {
    holdingId: string;
    userId: string;
    kind: string;
    type: HoldingDataType;
    name: string;
    description?: string;
    currency: string;
    openedAt: Timestamp;
    closedAt?: Timestamp | null;
}

export interface HoldingValuation {
    valuationId: string;
    holdingId: string;
    value: string;
    quantity?: string | null;
    unitPrice?: string | null;
    observedAt: Timestamp;
    note?: string | null;
}

export interface HoldingWithValue extends Holding {
    currentValue?: string | null;
    valuedAt?: Timestamp | null;
}

export type HoldingsListResponse = ListResponse<HoldingWithValue>;
export type ValuationsListResponse = ListResponse<HoldingValuation>;

export const assetKindOptions = [
    { label: 'Cash', value: 'cash' },
    { label: 'Bonds', value: 'bonds' },
    { label: 'Mutual Funds / ETFs', value: 'funds' },
    { label: 'Retirement Account', value: 'retirement' },
    { label: 'Real Estate', value: 'real_estate' },
    { label: 'Cryptocurrency', value: 'crypto' },
    { label: 'Gold', value: 'gold' },
    { label: 'Silver', value: 'silver' },
    { label: 'Business Equity', value: 'business_equity' },
    { label: 'Other', value: 'other' },
];

export const liabilityKindOptions = [
    { label: 'Mortgage', value: 'mortgage' },
    { label: 'Auto Loan', value: 'auto_loan' },
    { label: 'Credit Card', value: 'credit_card' },
    { label: 'Student Loan', value: 'student_loan' },
    { label: 'Personal Loan', value: 'personal_loan' },
    { label: 'Line of Credit', value: 'line_of_credit' },
    { label: 'Business Loan', value: 'business_loan' },
    { label: 'Medical Debt', value: 'medical_debt' },
    { label: 'Taxes Payable', value: 'taxes_payable' },
    { label: 'Other', value: 'other' },
];

export const kindOptionsByType: Record<HoldingDataType, { label: string; value: string }[]> = {
    [HoldingType.Asset]: assetKindOptions,
    [HoldingType.Liability]: liabilityKindOptions,
};
