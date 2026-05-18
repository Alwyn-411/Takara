export type ActiveFlag = 0 | 1;

export interface Base {
  userId: string;
  active: ActiveFlag;
}

export interface Timestamp {
  createdAt: number;
  updatedAt: number;
}
