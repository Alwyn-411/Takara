import { Route, Routes } from "react-router-dom";
import { ContentLayout } from "../components/core/Layout";
import { Home } from "../components/pages/Home";
import { Login } from "../components/pages/Login";
import { Register } from "../components/pages/Register";

export const Router = () => {
  return (
    <>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/signup" element={<Register />} />
        <Route element={<ContentLayout />}>
          <Route path="/home" element={<Home />} />
        </Route>
      </Routes>
    </>
  );
};
