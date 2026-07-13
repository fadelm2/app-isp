import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';

export const ProtectedRoute = ({ allowedRoles }) => {
    const token = localStorage.getItem('token');

    if (!token) {
        return <Navigate to="/login" replace />;
    }

    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(
            window.atob(base64)
                .split('')
                .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                .join('')
        );
        const payload = JSON.parse(jsonPayload);
        const userRole = parseInt(payload.role, 10);

        if (allowedRoles && !allowedRoles.includes(userRole)) {
            if (userRole === 2) {
                // Customers redirected to public payment
                return <Navigate to="/payment" replace />;
            }
            // Other unauthorized roles
            localStorage.removeItem('token');
            return <Navigate to="/login" replace />;
        }
    } catch (e) {
        console.error('Invalid token payload', e);
        localStorage.removeItem('token');
        return <Navigate to="/login" replace />;
    }

    return <Outlet />;
};

export const GuestRoute = () => {
    const token = localStorage.getItem('token');

    if (token) {
        try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(
                window.atob(base64)
                    .split('')
                    .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                    .join('')
            );
            const payload = JSON.parse(jsonPayload);
            const userRole = parseInt(payload.role, 10);

            if (userRole === 1 || userRole === 99) {
                return <Navigate to="/" replace />;
            } else if (userRole === 2) {
                return <Navigate to="/payment" replace />;
            }
        } catch (e) {
            localStorage.removeItem('token');
        }
    }

    return <Outlet />;
};
