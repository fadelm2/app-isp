import React, { useEffect, useState } from 'react';
import { adminService } from '../services/api';
import { Users, Wifi, ShieldAlert, DollarSign, Activity, RefreshCw } from 'lucide-react';
import './Dashboard.css';

export const Dashboard = () => {
    const [stats, setStats] = useState({
        total_customers: 0,
        active_customers: 0,
        suspended_customers: 0,
        owed_customers: 0,
        today_payments: 0,
        monthly_revenue: 0,
        router_status: 'offline',
        online_users: 0,
        offline_users: 0
    });
    const [loading, setLoading] = useState(true);

    const fetchStats = async () => {
        setLoading(true);
        try {
            const response = await adminService.getDashboardStats();
            if (response.data && response.data.data) {
                setStats(response.data.data);
            }
        } catch (error) {
            console.error("Error fetching dashboard statistics", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchStats();
    }, []);

    const formatRupiah = (val) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0
        }).format(val);
    };

    return (
        <div className="dashboard-container">
            <div className="dashboard-header">
                <h2>ISP Operations Overview</h2>
                <button onClick={fetchStats} className="refresh-btn">
                    <RefreshCw size={16} className={loading ? 'spin' : ''} />
                    Refresh
                </button>
            </div>

            {loading ? (
                <div className="loading-container">Loading system dashboard metrics...</div>
            ) : (
                <>
                    <div className="metrics-grid">
                        <div className="metric-card">
                            <div className="card-header">
                                <Users size={20} className="icon-purple" />
                                <span>Total Subscribers</span>
                            </div>
                            <div className="card-value">{stats.total_customers}</div>
                            <div className="card-subtitle">Active: {stats.active_customers} • Suspended: {stats.suspended_customers}</div>
                        </div>

                        <div className="metric-card">
                            <div className="card-header">
                                <Activity size={20} className="icon-green" />
                                <span>Live PPPoE Sessions</span>
                            </div>
                            <div className="card-value">{stats.online_users}</div>
                            <div className="card-subtitle">Offline RADIUS Users: {stats.offline_users}</div>
                        </div>

                        <div className="metric-card">
                            <div className="card-header">
                                <DollarSign size={20} className="icon-blue" />
                                <span>Monthly Revenue</span>
                            </div>
                            <div className="card-value">{formatRupiah(stats.monthly_revenue)}</div>
                            <div className="card-subtitle">Today's Collections: {formatRupiah(stats.today_payments)}</div>
                        </div>

                        <div className="metric-card">
                            <div className="card-header">
                                <ShieldAlert size={20} className="icon-red" />
                                <span>Overdue Accounts</span>
                            </div>
                            <div className="card-value">{stats.owed_customers}</div>
                            <div className="card-subtitle">Router API connection: <strong className={stats.router_status === 'online' ? 'text-green' : 'text-red'}>{stats.router_status.toUpperCase()}</strong></div>
                        </div>
                    </div>

                    <div className="system-health">
                        <h3>System Status Summary</h3>
                        <div className="status-list">
                            <div className="status-item">
                                <div className="status-info">
                                    <span className="dot dot-green"></span>
                                    <span>FreeRADIUS Authentication Database</span>
                                </div>
                                <span className="badge badge-success">Online & Synced</span>
                            </div>
                            <div className="status-item">
                                <div className="status-info">
                                    <span className={`dot ${stats.router_status === 'online' ? 'dot-green' : 'dot-red'}`}></span>
                                    <span>MikroTik PPP Active Monitor Gateway</span>
                                </div>
                                <span className={`badge ${stats.router_status === 'online' ? 'badge-success' : 'badge-danger'}`}>
                                    {stats.router_status === 'online' ? 'Active' : 'Offline / Mock Mode'}
                                </span>
                            </div>
                            <div className="status-item">
                                <div className="status-info">
                                    <span className="dot dot-green"></span>
                                    <span>Midtrans Snap Webhook Endpoint</span>
                                </div>
                                <span className="badge badge-success">Listening</span>
                            </div>
                            <div className="status-item">
                                <div className="status-info">
                                    <span className="dot dot-green"></span>
                                    <span>Automated Suspension Scheduler</span>
                                </div>
                                <span className="badge badge-success">Active (Hourly)</span>
                            </div>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
};
