import { Route, Routes } from "react-router-dom";
import { ContentLayout } from "../components/core/Layout";
import { Home } from "../components/pages/Home";
import { Login } from "../components/pages/Auth/Login";
import { Register } from "../components/pages/Auth/Register";
import { Accounts } from "../components/pages/Accounts/Accounts";
import { AccountCreation } from "../components/pages/Accounts/AccountCreation";

export const Router = () => {
  return (
    <>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/signup" element={<Register />} />
        <Route element={<ContentLayout />}>
          <Route path="/home" element={<Home />} />
          <Route path="/accounts" element={<Accounts />} />
          <Route path="/accounts/create" element={< AccountCreation />} />
        </Route>
      </Routes>
    </>
  );
};
