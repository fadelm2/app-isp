import React from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { authService } from '../../services/api';
import {
    LayoutDashboard,
    ClipboardList,
    Users,
    Wifi,
    Network,
    Receipt,
    Settings,
    LogOut
} from 'lucide-react';
import './Sidebar.css';

const Sidebar = ({ isOpen, currentUser }) => {
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            await authService.logout();
        } catch (e) {
            console.error("Logout API failed", e);
        }
        localStorage.removeItem('token');
        navigate('/login');
    };

    const ispName = localStorage.getItem('isp_name') || 'GREENET';

    return (
        <div className={`sidebar ${isOpen ? '' : 'collapsed'}`}>
            <div className="brand">
                <div className="logo-icon">
                    <div className="grid-icon"></div>
                </div>
                <span className="brand-name">{ispName.toUpperCase()}</span>
            </div>

            <nav className="nav-menu">
                <NavLink to="/" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
                    <LayoutDashboard size={20} />
                    <span>Dashboard</span>
                </NavLink>
                <NavLink to="/registrations" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
                    <ClipboardList size={20} />
                    <span>Registrations</span>
                </NavLink>
                <NavLink to="/customers" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
                    <Users size={20} />
                    <span>Customers</span>
                </NavLink>
                <NavLink to="/packages" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
                    <Wifi size={20} />
                    <span>Internet Packages</span>
                </NavLink>
                <NavLink to="/routers" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
                    <Network size={20} />
                    <span>Router Management</span>
                </NavLink>
                <NavLink to="/invoices" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
                    <Receipt size={20} />
                    <span>Invoices & Billing</span>
                </NavLink>
                <NavLink to="/settings" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''} settings-item`}>
                    <Settings size={20} />
                    <span>Settings</span>
                </NavLink>
                <button onClick={handleLogout} className="nav-item logout-item" style={{ background: 'none', border: 'none', width: '100%', textAlign: 'left', cursor: 'pointer' }}>
                    <LogOut size={20} />
                    <span>Logout</span>
                </button>
            </nav>

            <div className="user-profile-bottom">
                <img src={`https://ui-avatars.com/api/?name=${currentUser ? encodeURIComponent(currentUser.username) : `Admin+${encodeURIComponent(ispName)}`}&background=0D8ABC&color=fff`} alt="User" />
                <div className="user-info">
                    <div className="name">{currentUser ? currentUser.username : `Admin ${ispName}`}</div>
                    <div className="role">{currentUser ? currentUser.id : 'Administrator'}</div>
                </div>
            </div>
        </div>
    );
};

export default Sidebar;
