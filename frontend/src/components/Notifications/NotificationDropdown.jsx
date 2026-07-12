import React, { useState } from 'react';
import { Bell, X } from 'lucide-react';
import './NotificationDropdown.css';

const NotificationDropdown = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [activeTab, setActiveTab] = useState('all');

    const notifications = [
        {
            id: 1,
            avatar: 'https://ui-avatars.com/api/?name=Ram&background=random',
            message: 'Penjual menanggapi komentar Anda di Ram.',
            time: '16 jam',
            unread: false,
            icon: '💬'
        },
        {
            id: 2,
            avatar: 'https://ui-avatars.com/api/?name=Marketplace&background=FF9900',
            message: 'Kami menolak tawaran Marketplace Anda. Lihat alasannya.',
            time: '15 jam',
            unread: false,
            icon: '⚠️'
        },
        {
            id: 3,
            avatar: 'https://ui-avatars.com/api/?name=Ram&background=random',
            message: 'Ram tersedia. Kirim pesan kepada penjual jika Anda masih tertarik.',
            time: '16 jam',
            unread: true,
            icon: '💬'
        },
        {
            id: 4,
            avatar: 'https://ui-avatars.com/api/?name=Geri+Ferial&background=random',
            message: 'Geri Ferial menyebut Anda di sebuah komentar di Jual Beli Pc Gaming Dan Perlengkapan Pc.',
            time: '16 jam',
            unread: false,
            icon: '💬'
        },
    ];

    const toggleDropdown = () => {
        setIsOpen(!isOpen);
    };

    return (
        <div className="notification-wrapper">
            <div className="action-button notification-bell" onClick={toggleDropdown}>
                <Bell size={18} />
                <span className="notification-count">3</span>
            </div>

            {isOpen && (
                <>
                    <div className="notification-overlay" onClick={toggleDropdown}></div>
                    <div className="notification-dropdown">
                        <div className="notification-header">
                            <h3>Notifikasi</h3>
                            <button className="close-btn" onClick={toggleDropdown}>
                                <X size={20} />
                            </button>
                        </div>

                        <div className="notification-tabs">
                            <button
                                className={`tab ${activeTab === 'all' ? 'active' : ''}`}
                                onClick={() => setActiveTab('all')}
                            >
                                Semua
                            </button>
                            <button
                                className={`tab ${activeTab === 'unread' ? 'active' : ''}`}
                                onClick={() => setActiveTab('unread')}
                            >
                                Belum Dibaca
                            </button>
                        </div>

                        <div className="notification-section">
                            <div className="section-header">
                                <span>Hari ini</span>
                            </div>

                            <div className="notification-list">
                                {notifications
                                    .filter(n => activeTab === 'all' || (activeTab === 'unread' && n.unread))
                                    .map(notification => (
                                        <div
                                            key={notification.id}
                                            className={`notification-item ${notification.unread ? 'unread' : ''}`}
                                        >
                                            <div className="notification-avatar">
                                                <img src={notification.avatar} alt="Avatar" />
                                                <span className="notification-icon">{notification.icon}</span>
                                            </div>
                                            <div className="notification-content">
                                                <p>{notification.message}</p>
                                                <span className="notification-time">{notification.time}</span>
                                            </div>
                                            {notification.unread && <div className="unread-dot"></div>}
                                        </div>
                                    ))}
                            </div>
                        </div>

                        <div className="notification-footer">
                            <button className="see-all-btn">Lihat semua</button>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
};

export default NotificationDropdown;
