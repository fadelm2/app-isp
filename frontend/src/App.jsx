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
import { AdminRegister } from './pages/AdminRegister';
import { PublicPayment } from './pages/PublicPayment';
import { ProtectedRoute, GuestRoute } from './components/ProtectedRoute';
import { useEffect } from 'react';
import { publicService } from './services/api';

function App() {
    useEffect(() => {
        publicService.getIspInfo()
            .then(res => {
                if (res.data && res.data.isp_name) {
                    localStorage.setItem('isp_name', res.data.isp_name);
                    document.title = `${res.data.isp_name} - ISP Management System`;
                }
            })
            .catch(err => console.error("Error fetching ISP info:", err));
    }, []);

    return (
        <Router>
            <Routes>
                {/* Guest-only Routes (Redirects to dashboard if already logged in) */}
                <Route element={<GuestRoute />}>
                    <Route path="/login" element={<Login />} />
                </Route>

                {/* Public Guest Routes */}
                <Route path="/register" element={<PublicRegister />} />
                <Route path="/payment" element={<PublicPayment />} />
                <Route path="/payment/:customerId" element={<PublicPayment />} />

                {/* Admin/NOC Protected Routes */}
                <Route element={<ProtectedRoute allowedRoles={[1, 99]} />}>
                    <Route path="/" element={<Layout />}>
                        <Route index element={<Dashboard />} />
                        <Route path="registrations" element={<Registrations />} />
                        <Route path="registrations/new" element={<AdminRegister />} />
                        <Route path="customers" element={<Customers />} />
                        <Route path="packages" element={<Packages />} />
                        <Route path="routers" element={<Routers />} />
                        <Route path="invoices" element={<Invoices />} />
                    </Route>
                </Route>
            </Routes>
        </Router>
    );
}

export default App;
