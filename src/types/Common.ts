const DataStatusFlag = {
    Active: 1,
    InActive: 0,
} as const;

export type DataStatus = (typeof DataStatusFlag)[keyof typeof DataStatusFlag];

export interface Base {
    userId: string;
    active: DataStatus;
}

export interface Timestamp {
    createdAt: number;
    updatedAt: number;
}
