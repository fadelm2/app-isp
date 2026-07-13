import React, { useState, useEffect } from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import TopBar from './TopBar';
import { authService } from '../../services/api';
import './Layout.css';

export const Layout = () => {
    const [isSidebarOpen, setIsSidebarOpen] = useState(true);
    const [currentUser, setCurrentUser] = useState(null);

    const toggleSidebar = () => {
        setIsSidebarOpen(!isSidebarOpen);
    };

    const token = localStorage.getItem('token');

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await authService.getCurrentUser();
                if (response.data && response.data.data) {
                    setCurrentUser(response.data.data);
                }
            } catch (err) {
                console.error("Failed to load user profile in Layout", err);
            }
        };
        if (token) {
            fetchUser();
        }
    }, [token]);

    if (!token) {
        return <Navigate to="/login" replace />;
    }

    return (
        <div className={`app-container ${isSidebarOpen ? '' : 'sidebar-collapsed'}`}>
            <Sidebar isOpen={isSidebarOpen} currentUser={currentUser} />
            <main className="main-content">
                <TopBar onToggleSidebar={toggleSidebar} currentUser={currentUser} />
                <div className="content-area">
                    <Outlet />
                </div>
            </main>
        </div>
    );
};
