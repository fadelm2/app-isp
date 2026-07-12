import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout/Layout';
import { Dashboard } from './pages/Dashboard';
import { Registrations } from './pages/Registrations';
import { Customers } from './pages/Customers';
import { Packages } from './pages/Packages';
import { Routers } from './pages/Routers';
import { Invoices } from './pages/Invoices';
import { Login } from './pages/Login';
import { PublicRegister } from './pages/PublicRegister';

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<PublicRegister />} />
                <Route path="/" element={<Layout />}>
                    <Route index element={<Dashboard />} />
                    <Route path="registrations" element={<Registrations />} />
                    <Route path="customers" element={<Customers />} />
                    <Route path="packages" element={<Packages />} />
                    <Route path="routers" element={<Routers />} />
                    <Route path="invoices" element={<Invoices />} />
                </Route>
            </Routes>
        </Router>
    );
}

export default App;
