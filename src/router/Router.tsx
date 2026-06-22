import { Route, Routes } from 'react-router-dom';
import { ContentLayout } from '../components/core/Layout';
import { AccountCreate } from '../components/pages/Accounts/Create';
import { AccountDetails } from '../components/pages/Accounts/Details';
import { AccountEdit } from '../components/pages/Accounts/Edit';
import { Accounts as AccountOverview } from '../components/pages/Accounts/Overview';
import { Login } from '../components/pages/Auth/Login';
import { Register } from '../components/pages/Auth/Register';
import { HoldingsCreate } from '../components/pages/Holdings/Create';
import { ValuationCreate } from '../components/pages/Holdings/CreateValuations';
import { Holdings as HoldingsOverview } from '../components/pages/Holdings/Overview';
import { HoldingTrends } from '../components/pages/Holdings/Trends';
import { Home } from '../components/pages/Home';
import { TransactionsCreate } from '../components/pages/Transactions/Create';
import { Profile } from '../components/pages/User/Profile';

export const Router = () => {
    return (
        <>
            <Routes>
                <Route path="/" element={<Login />} />
                <Route path="/signup" element={<Register />} />
                <Route element={<ContentLayout />}>
                    <Route path="/home" element={<Home />} />
                    <Route path="/accounts" element={<AccountOverview />} />
                    <Route path="/accounts/create" element={<AccountCreate />} />
                    <Route path="/accounts/:accountId/edit" element={<AccountEdit />} />
                    <Route path="/accounts/:accountId/details" element={<AccountDetails />} />
                    <Route path="/accounts/:accountId/transactions/create" element={<TransactionsCreate />} />
                    <Route path="/holdings" element={<HoldingsOverview />} />
                    <Route path="/holdings/create" element={<HoldingsCreate />} />
                    <Route path="/holdings/:holdingId/trends" element={<HoldingTrends />} />
                    <Route path="/holdings/:holdingId/trends/add" element={<ValuationCreate />} />
                    <Route path="/profile" element={<Profile />} />
                </Route>
            </Routes>
        </>
    );
};
