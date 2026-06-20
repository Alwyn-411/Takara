import type { Holding, HoldingsListResponse, HoldingValuation, ValuationsListResponse } from '../types/Holdings';
import type { CreateResponse, EditOrDeleteResponse } from './default';

import { axiosInstance } from './client';

export type CreateHoldingRequest = Pick<Holding, 'kind' | 'type' | 'name' | 'currency'> &
    Partial<Pick<Holding, 'description' | 'openedAt'>> &
    Pick<HoldingValuation, 'value'> &
    Partial<Pick<HoldingValuation, 'quantity' | 'unitPrice' | 'observedAt' | 'note'>>;

export type UpdateHoldingRequest = Partial<Pick<Holding, 'type' | 'name' | 'description' | 'closedAt'>>;

export type RecordValuationRequest = Pick<HoldingValuation, 'value'> &
    Partial<Pick<HoldingValuation, 'quantity' | 'unitPrice' | 'observedAt' | 'note'>>;

export const createHolding = async (data: CreateHoldingRequest): Promise<CreateResponse> => {
    const response = await axiosInstance.post<CreateResponse>('/v1/holdings', data);
    return response.data;
};

export const createHoldingValuation = async (data: RecordValuationRequest & { holdingId: string }): Promise<CreateResponse> => {
    const response = await axiosInstance.post<CreateResponse>(`/v1/holdings/${data.holdingId}/valuations`, data);
    return response.data;
};

export const getHoldings = async (): Promise<HoldingsListResponse> => {
    const response = await axiosInstance.get<HoldingsListResponse>(`/v1/holdings`);
    return response.data;
};

export const getHoldingById = async (data: { holdingId: string }): Promise<Holding> => {
    const response = await axiosInstance.get<Holding>(`/v1/holdings/${data.holdingId}`);
    return response.data;
};

export const getHoldingValuations = async (data: { holdingId: string }): Promise<ValuationsListResponse> => {
    const response = await axiosInstance.get<ValuationsListResponse>(`/v1/holdings/${data.holdingId}/valuations`);
    return response.data;
};

export const updateHoldingById = async (data: UpdateHoldingRequest & { holdingId: string }): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.get<EditOrDeleteResponse>(`/v1/holdings/${data.holdingId}`);
    return response.data;
};

export const deleteHolding = async (data: { holdingId: string }): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.delete<EditOrDeleteResponse>(`/v1/holdings/${data.holdingId}`);
    return response.data;
};

export const deleteValuation = async (data: { valuationId: string }): Promise<EditOrDeleteResponse> => {
    const response = await axiosInstance.delete<EditOrDeleteResponse>(`/v1/holdings/${data.valuationId}`);
    return response.data;
};
