import React from 'react';
import { Trophy, Maximize, MoreHorizontal, Menu } from 'lucide-react';
import NotificationDropdown from '../Notifications/NotificationDropdown';
import './TopBar.css';

const TopBar = ({ onToggleSidebar, currentUser }) => {
    return (
        <div className="topbar">
            <div className="greeting-section" style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
                <button className="menu-btn" onClick={onToggleSidebar}>
                    <Menu size={24} />
                </button>
                <div>
                    <h1>Welcome back, {currentUser ? currentUser.username : 'Admin'}!</h1>
                    <p>ISP Management Dashboard</p>
                </div>
            </div>

            <div className="actions-section">
                <NotificationDropdown />

                <div className="action-button">
                    <Trophy size={18} />
                    <div className="notification-dot"></div>
                </div>

                <div className="action-button">
                    <Maximize size={18} />
                </div>

                <div className="action-button">
                    <MoreHorizontal size={18} />
                </div>
            </div>
        </div>
    );
};

export default TopBar;
