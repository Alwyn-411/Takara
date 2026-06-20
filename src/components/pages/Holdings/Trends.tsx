import { useQuery } from '@tanstack/react-query';
import { useUserStore } from '../../../store/User';
import { getHoldingById, getHoldingValuations } from '../../../api/holdings';
import { useParams } from 'react-router-dom';

export const HoldingTrends = () => {
    const userId = useUserStore.getState().userId;
    const { holdingId } = useParams<{ holdingId: string }>();

    const HoldingValuationQuery = useQuery({
        queryKey: ['HoldingTrends', userId, holdingId],
        queryFn: () => getHoldingValuations({ holdingId }),
        enabled: !!userId && !!holdingId,
    });

    const HoldingQuery = useQuery({
        queryKey: ['Holding', userId, holdingId],
        queryFn: () => getHoldingById({ holdingId }),
        enabled: !!userId && !!holdingId,
    });

    console.log(HoldingQuery.data);
    console.log(HoldingValuationQuery.data);

    return <></>;
};
