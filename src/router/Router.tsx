import { Route, Routes } from 'react-router-dom';
import { ContentLayout } from '../components/core/Layout';
import { Home } from '../components/pages/Home';
import { Login } from '../components/pages/Auth/Login';
import { Register } from '../components/pages/Auth/Register';
import { AccountCreate } from '../components/pages/Accounts/Create';
import { AccountEdit } from '../components/pages/Accounts/Edit';
import { Accounts as AccountOverview } from '../components/pages/Accounts/Overview';
import { AccountDetails } from '../components/pages/Accounts/Details';

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
                </Route>
            </Routes>
        </>
    );
};
