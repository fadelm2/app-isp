import React, { useEffect, useState } from 'react';
import { adminService } from '../services/api';
import { ShieldAlert, ShieldCheck, PowerOff, History } from 'lucide-react';
import './Pages.css';

export const Customers = () => {
    const [customers, setCustomers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [selectedHistory, setSelectedHistory] = useState(null);
    const [historyList, setHistoryList] = useState([]);
    const [loadingHistory, setLoadingHistory] = useState(false);

    const fetchCustomers = async () => {
        try {
            const response = await adminService.getCustomers();
            if (response.data && response.data.data) {
                setCustomers(response.data.data);
            }
        } catch (error) {
            console.error("Error fetching customers", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCustomers();
    }, []);

    const handleSuspend = async (id) => {
        const notes = prompt("Enter suspension reason:");
        if (notes === null) return;
        try {
            await adminService.suspendCustomer(id, notes);
            fetchCustomers();
        } catch (error) {
            alert("Failed to suspend customer: " + (error.response?.data?.error || error.message));
        }
    };

    const handleUnsuspend = async (id) => {
        const notes = prompt("Enter reactivation notes (optional):") || "Reactivated by admin";
        try {
            await adminService.unsuspendCustomer(id, notes);
            fetchCustomers();
        } catch (error) {
            alert("Failed to unsuspend customer: " + (error.response?.data?.error || error.message));
        }
    };

    const handleTerminate = async (id) => {
        if (!confirm("Are you sure you want to TERMINATE this subscriber? This will delete their RADIUS access.")) return;
        const notes = prompt("Enter termination reason:") || "Contract terminated";
        try {
            await adminService.terminateCustomer(id, notes);
            fetchCustomers();
        } catch (error) {
            alert("Failed to terminate customer: " + (error.response?.data?.error || error.message));
        }
    };

    const handleViewHistory = async (id, name) => {
        setSelectedHistory({ id, name });
        setLoadingHistory(true);
        try {
            const response = await adminService.getCustomerHistory(id);
            if (response.data && response.data.data) {
                setHistoryList(response.data.data);
            }
        } catch (error) {
            console.error("Failed to load customer history", error);
        } finally {
            setLoadingHistory(false);
        }
    };

    return (
        <div className="page-container">
            <div className="page-header">
                <h2>Active Subscribers</h2>
            </div>

            {loading ? (
                <div className="loading">Loading subscribers list...</div>
            ) : customers.length === 0 ? (
                <div className="empty-state">No subscribers found.</div>
            ) : (
                <div className="table-responsive">
                    <table className="custom-table">
                        <thead>
                            <tr>
                                <th>Cust ID</th>
                                <th>Subscriber</th>
                                <th>PPP / RADIUS Username</th>
                                <th>Internet Package</th>
                                <th>Billing Day</th>
                                <th>Status</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {customers.map((cust) => (
                                <tr key={cust.id}>
                                    <td className="bold">{cust.id}</td>
                                    <td>
                                        <div className="bold">{cust.user ? cust.user.name : 'N/A'}</div>
                                        <div className="text-muted small">{cust.user ? cust.user.email : ''}</div>
                                    </td>
                                    <td>
                                        <code>{cust.ppp_username}</code>
                                    </td>
                                    <td>{cust.package ? `${cust.package.name} (${cust.package.speed_mbps} Mbps)` : 'N/A'}</td>
                                    <td>Day {cust.due_date_day}</td>
                                    <td>
                                        <span className={`badge-status status-${cust.status}`}>
                                            {cust.status.toUpperCase()}
                                        </span>
                                    </td>
                                    <td>
                                        <div className="action-buttons">
                                            {cust.status === 'active' ? (
                                                <button 
                                                    onClick={() => handleSuspend(cust.id)} 
                                                    className="btn btn-warning btn-sm"
                                                    title="Suspend Service"
                                                >
                                                    <ShieldAlert size={14} /> Suspend
                                                </button>
                                            ) : cust.status === 'suspended' ? (
                                                <button 
                                                    onClick={() => handleUnsuspend(cust.id)} 
                                                    className="btn btn-success btn-sm"
                                                    title="Unsuspend Service"
                                                >
                                                    <ShieldCheck size={14} /> Unsuspend
                                                </button>
                                            ) : null}
                                            {cust.status !== 'terminated' && (
                                                <button 
                                                    onClick={() => handleTerminate(cust.id)} 
                                                    className="btn btn-danger btn-sm"
                                                    title="Terminate Account"
                                                >
                                                    <PowerOff size={14} /> Terminate
                                                </button>
                                            )}
                                            <button 
                                                onClick={() => handleViewHistory(cust.id, cust.user?.name)} 
                                                className="btn btn-secondary btn-sm"
                                                title="View History Logs"
                                            >
                                                <History size={14} /> History
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {/* History Modal */}
            {selectedHistory && (
                <div className="modal-backdrop">
                    <div className="custom-modal">
                        <div className="modal-header">
                            <h3>Audit Log: {selectedHistory.name}</h3>
                            <button onClick={() => setSelectedHistory(null)} className="close-btn">&times;</button>
                        </div>
                        <div className="modal-body">
                            {loadingHistory ? (
                                <p>Loading log history...</p>
                            ) : historyList.length === 0 ? (
                                <p className="text-muted">No state history recorded.</p>
                            ) : (
                                <div className="history-timeline">
                                    {historyList.map((log) => (
                                        <div className="timeline-item" key={log.id}>
                                            <div className="timeline-meta">
                                                <span className={`badge-status status-${log.action}`}>{log.action.toUpperCase()}</span>
                                                <span className="time">{new Date(log.created_at).toLocaleString()}</span>
                                            </div>
                                            <div className="timeline-content">
                                                <p>{log.notes}</p>
                                                <small className="text-muted">Executed by: {log.user?.name || log.created_by}</small>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};
