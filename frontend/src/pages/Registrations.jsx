import React, { useEffect, useState } from 'react';
import { adminService } from '../services/api';
import { ClipboardList, Check, X, FileText, MapPin, Loader } from 'lucide-react';
import './Pages.css';

export const Registrations = () => {
    const [registrations, setRegistrations] = useState([]);
    const [loading, setLoading] = useState(true);
    const [updatingId, setUpdatingId] = useState(null);

    const fetchRegistrations = async () => {
        try {
            const response = await adminService.getRegistrations();
            if (response.data && response.data.data) {
                setRegistrations(response.data.data);
            }
        } catch (error) {
            console.error("Error fetching registrations", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchRegistrations();
    }, []);

    const handleUpdateStatus = async (id, status) => {
        setUpdatingId(id);
        try {
            await adminService.updateRegistrationStatus(id, status);
            fetchRegistrations();
        } catch (error) {
            console.error("Failed to update status", error);
            alert("Error: " + (error.response?.data?.error || "Failed to update registration"));
        } finally {
            setUpdatingId(null);
        }
    };

    return (
        <div className="page-container">
            <div className="page-header">
                <h2>Customer Registrations</h2>
            </div>

            {loading ? (
                <div className="loading">Loading registrations...</div>
            ) : registrations.length === 0 ? (
                <div className="empty-state">No customer registrations found.</div>
            ) : (
                <div className="table-responsive">
                    <table className="custom-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>NIK</th>
                                <th>Email / Phone</th>
                                <th>Package</th>
                                <th>Location</th>
                                <th>Status</th>
                                <th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            {registrations.map((reg) => (
                                <tr key={reg.id}>
                                    <td>
                                        <div className="bold">{reg.full_name}</div>
                                        <div className="text-muted small">{reg.installation_address}</div>
                                    </td>
                                    <td>{reg.nik}</td>
                                    <td>
                                        <div>{reg.email}</div>
                                        <div className="text-muted">{reg.phone}</div>
                                    </td>
                                    <td>{reg.package ? reg.package.name : reg.package_id}</td>
                                    <td>
                                        {reg.latitude && reg.longitude ? (
                                            <a 
                                                href={`https://www.google.com/maps?q=${reg.latitude},${reg.longitude}`} 
                                                target="_blank" 
                                                rel="noreferrer"
                                                className="map-link"
                                            >
                                                <MapPin size={14} /> Open Maps
                                            </a>
                                        ) : 'N/A'}
                                    </td>
                                    <td>
                                        <span className={`badge-status status-${reg.status}`}>
                                            {reg.status.toUpperCase()}
                                        </span>
                                    </td>
                                    <td>
                                        {reg.status === 'pending' && (
                                            <div className="action-buttons">
                                                <button 
                                                    onClick={() => handleUpdateStatus(reg.id, 'under_review')}
                                                    className="btn btn-secondary btn-sm"
                                                    disabled={updatingId !== null}
                                                >
                                                    Under Review
                                                </button>
                                            </div>
                                        )}
                                        {reg.status === 'under_review' && (
                                            <div className="action-buttons">
                                                <button 
                                                    onClick={() => handleUpdateStatus(reg.id, 'surveying')}
                                                    className="btn btn-warning btn-sm"
                                                    disabled={updatingId !== null}
                                                >
                                                    Survey
                                                </button>
                                            </div>
                                        )}
                                        {reg.status === 'surveying' && (
                                            <div className="action-buttons">
                                                <button 
                                                    onClick={() => handleUpdateStatus(reg.id, 'approved')}
                                                    className="btn btn-success btn-sm btn-icon"
                                                    disabled={updatingId !== null}
                                                    title="Approve"
                                                >
                                                    <Check size={14} /> Approve
                                                </button>
                                                <button 
                                                    onClick={() => handleUpdateStatus(reg.id, 'rejected')}
                                                    className="btn btn-danger btn-sm btn-icon"
                                                    disabled={updatingId !== null}
                                                    title="Reject"
                                                >
                                                    <X size={14} /> Reject
                                                </button>
                                            </div>
                                        )}
                                        {reg.status === 'approved' && (
                                            <span className="text-success small font-medium">Activated</span>
                                        )}
                                        {reg.status === 'rejected' && (
                                            <span className="text-danger small font-medium">Rejected</span>
                                        )}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
};
