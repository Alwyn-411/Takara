import { Route, Routes } from 'react-router-dom';
import { ContentLayout } from '../components/core/Layout';
import { Home } from '../components/pages/Home';
import { Login } from '../components/pages/Auth/Login';
import { Register } from '../components/pages/Auth/Register';
import { Accounts } from '../components/pages/Accounts/Accounts';
import { AccountCreate } from '../components/pages/Accounts/AccountCreate';
import { AccountEdit } from '../components/pages/Accounts/AccountEdit';

export const Router = () => {
    return (
        <>
            <Routes>
                <Route path="/" element={<Login />} />
                <Route path="/signup" element={<Register />} />
                <Route element={<ContentLayout />}>
                    <Route path="/home" element={<Home />} />
                    <Route path="/accounts" element={<Accounts />} />
                    <Route path="/accounts/create" element={<AccountCreate />} />
                    <Route path="/accounts/:accountId/edit" element={<AccountEdit />} />
                </Route>
            </Routes>
        </>
    );
};
