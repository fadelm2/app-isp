import React, { useEffect, useState } from 'react';
import { adminService } from '../services/api';
import { ShieldAlert, ShieldCheck, PowerOff, History, Eye, Edit, Save, X, FileImage, MapPin } from 'lucide-react';
import './Pages.css';

export const Customers = () => {
    const [customers, setCustomers] = useState([]);
    const [packages, setPackages] = useState([]);
    const [routers, setRouters] = useState([]);
    const [loading, setLoading] = useState(true);

    // Search and Pagination States
    const [search, setSearch] = useState('');
    const [status, setStatus] = useState('');
    const [page, setPage] = useState(1);
    const [size, setSize] = useState(10);
    const [totalItem, setTotalItem] = useState(0);
    const [totalPage, setTotalPage] = useState(0);

    // History Modal
    const [selectedHistory, setSelectedHistory] = useState(null);
    const [historyList, setHistoryList] = useState([]);
    const [loadingHistory, setLoadingHistory] = useState(false);

    // Detail & Edit Modal
    const [selectedCust, setSelectedCust] = useState(null);
    const [isEditing, setIsEditing] = useState(false);
    const [saving, setSaving] = useState(false);

    // Edit Form State
    const [formPackageId, setFormPackageId] = useState('');
    const [formRouterId, setFormRouterId] = useState('');
    const [formPppUsername, setFormPppUsername] = useState('');
    const [formPppPassword, setFormPppPassword] = useState('');
    const [formDueDateDay, setFormDueDateDay] = useState(10);
    const [formOdpNumber, setFormOdpNumber] = useState('');

    const fetchCustomers = async (currentPage = page, searchKeyword = search, statusFilter = status, pageSize = size) => {
        setLoading(true);
        try {
            const response = await adminService.getCustomers({
                search: searchKeyword,
                status: statusFilter,
                page: currentPage,
                size: pageSize
            });
            if (response.data && response.data.data) {
                setCustomers(response.data.data);
                if (response.data.paging) {
                    setTotalItem(response.data.paging.total_item);
                    setTotalPage(response.data.paging.total_page);
                }
            }
        } catch (error) {
            console.error("Error fetching customers", error);
        } finally {
            setLoading(false);
        }
    };

    const fetchDropdowns = async () => {
        try {
            const [pkgResp, rtrResp] = await Promise.all([
                adminService.getPackages(),
                adminService.getRouters()
            ]);
            if (pkgResp.data && pkgResp.data.data) setPackages(pkgResp.data.data);
            if (rtrResp.data && rtrResp.data.data) setRouters(rtrResp.data.data);
        } catch (error) {
            console.error("Error fetching dropdown options", error);
        }
    };

    // Initial load
    useEffect(() => {
        fetchDropdowns();
    }, []);

    // Debounced search trigger
    useEffect(() => {
        const delayDebounceFn = setTimeout(() => {
            setPage(1);
            fetchCustomers(1, search, status, size);
        }, 300);

        return () => clearTimeout(delayDebounceFn);
    }, [search]);

    // Direct filters triggers
    const handleStatusChange = (newStatus) => {
        setStatus(newStatus);
        setPage(1);
        fetchCustomers(1, search, newStatus, size);
    };

    const handleSizeChange = (newSize) => {
        setSize(newSize);
        setPage(1);
        fetchCustomers(1, search, status, newSize);
    };

    const handlePageChange = (newPage) => {
        setPage(newPage);
        fetchCustomers(newPage, search, status, size);
    };

    const handleSuspend = async (id) => {
        const notes = prompt("Enter suspension reason:");
        if (notes === null) return;
        try {
            await adminService.suspendCustomer(id, notes);
            fetchCustomers(page, search, status, size);
            if (selectedCust && selectedCust.id === id) {
                setSelectedCust(prev => ({ ...prev, status: 'suspended' }));
            }
        } catch (error) {
            alert("Failed to suspend customer: " + (error.response?.data?.error || error.message));
        }
    };

    const handleUnsuspend = async (id) => {
        const notes = prompt("Enter reactivation notes (optional):") || "Reactivated by admin";
        try {
            await adminService.unsuspendCustomer(id, notes);
            fetchCustomers(page, search, status, size);
            if (selectedCust && selectedCust.id === id) {
                setSelectedCust(prev => ({ ...prev, status: 'active' }));
            }
        } catch (error) {
            alert("Failed to unsuspend customer: " + (error.response?.data?.error || error.message));
        }
    };

    const handleTerminate = async (id) => {
        if (!confirm("Are you sure you want to TERMINATE this subscriber? This will delete their RADIUS access.")) return;
        const notes = prompt("Enter termination reason:") || "Contract terminated";
        try {
            await adminService.terminateCustomer(id, notes);
            fetchCustomers(page, search, status, size);
            if (selectedCust && selectedCust.id === id) {
                setSelectedCust(prev => ({ ...prev, status: 'terminated' }));
            }
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

    const handleOpenDetail = (cust) => {
        setSelectedCust(cust);
        setIsEditing(false);
        setFormPackageId(cust.package_id || '');
        setFormRouterId(cust.router_id || '');
        setFormPppUsername(cust.ppp_username || '');
        setFormPppPassword(cust.ppp_password || '');
        setFormDueDateDay(cust.due_date_day || 10);
        setFormOdpNumber(cust.odp_number || '');
    };

    const handleSaveChanges = async (e) => {
        e.preventDefault();
        setSaving(true);
        try {
            const payload = {
                package_id: formPackageId,
                router_id: formRouterId,
                ppp_username: formPppUsername,
                ppp_password: formPppPassword,
                due_date_day: parseInt(formDueDateDay, 10),
                odp_number: formOdpNumber
            };
            const response = await adminService.updateCustomer(selectedCust.id, payload);
            if (response.data && response.data.data) {
                alert("Customer profile updated successfully!");
                // Refresh list
                fetchCustomers(page, search, status, size);
                // Update selected details modal
                setSelectedCust(response.data.data);
                setIsEditing(false);
            }
        } catch (error) {
            console.error("Failed to save changes", error);
            alert("Error: " + (error.response?.data?.error || "Failed to save customer changes"));
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="page-container">
            <div className="page-header">
                <h2>Active Subscribers</h2>
            </div>

            {/* Search & Filter Bar */}
            <div className="filter-bar mb-4" style={{ display: 'flex', gap: '12px', flexWrap: 'wrap', alignItems: 'center', marginBottom: '16px' }}>
                <input 
                    type="text" 
                    placeholder="Search by ID, name, email, ODP, PPP..." 
                    value={search} 
                    onChange={(e) => setSearch(e.target.value)} 
                    className="form-input"
                    style={{ flex: 1, minWidth: '250px' }}
                />
                
                <select 
                    value={status} 
                    onChange={(e) => handleStatusChange(e.target.value)}
                    className="form-input"
                    style={{ minWidth: '150px', padding: '8px 12px' }}
                >
                    <option value="">All Statuses</option>
                    <option value="active">Active</option>
                    <option value="suspended">Suspended</option>
                    <option value="terminated">Terminated</option>
                </select>

                <select 
                    value={size} 
                    onChange={(e) => handleSizeChange(parseInt(e.target.value, 10))}
                    className="form-input"
                    style={{ minWidth: '130px', padding: '8px 12px' }}
                >
                    <option value={10}>10 per page</option>
                    <option value={20}>20 per page</option>
                    <option value={50}>50 per page</option>
                </select>
            </div>

            {loading ? (
                <div className="loading">Loading subscribers list...</div>
            ) : customers.length === 0 ? (
                <div className="empty-state">No subscribers found.</div>
            ) : (
                <>
                    <div className="table-responsive">
                        <table className="custom-table">
                            <thead>
                                <tr>
                                    <th>Cust ID</th>
                                    <th>Subscriber</th>
                                    <th>PPP / RADIUS Username</th>
                                    <th>Internet Package</th>
                                    <th>ODP Number</th>
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
                                            <div className="text-muted small italic mt-1">{cust.registration ? cust.registration.installation_address : ''}</div>
                                        </td>
                                        <td>
                                            <code>{cust.ppp_username}</code>
                                        </td>
                                        <td>{cust.package ? `${cust.package.name} (${cust.package.speed_mbps} Mbps)` : 'N/A'}</td>
                                        <td>
                                            <span className="font-mono text-xs">{cust.odp_number || '-'}</span>
                                        </td>
                                        <td>Day {cust.due_date_day}</td>
                                        <td>
                                            <span className={`badge-status status-${cust.status}`}>
                                                {cust.status.toUpperCase()}
                                            </span>
                                        </td>
                                        <td>
                                            <div className="action-buttons">
                                                <button 
                                                    onClick={() => handleOpenDetail(cust)}
                                                    className="btn btn-secondary btn-sm"
                                                    title="View/Edit Profile"
                                                >
                                                    <Eye size={14} /> Detail / Edit
                                                </button>
                                                <button 
                                                    onClick={() => handleViewHistory(cust.id, cust.user?.name)} 
                                                    className="btn btn-secondary btn-sm"
                                                    title="View History Logs"
                                                >
                                                    <History size={14} /> Log
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>

                    {/* Pagination Controls */}
                    <div className="pagination-bar" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '16px', padding: '0 8px' }}>
                        <span className="text-muted small">
                            Showing {customers.length} of {totalItem} subscribers
                        </span>
                        <div className="pagination-buttons" style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                            <button 
                                className="btn btn-secondary btn-sm" 
                                onClick={() => handlePageChange(page - 1)} 
                                disabled={page <= 1}
                            >
                                Previous
                            </button>
                            <span className="bold small" style={{ color: 'var(--text-primary)' }}>Page {page} of {totalPage || 1}</span>
                            <button 
                                className="btn btn-secondary btn-sm" 
                                onClick={() => handlePageChange(page + 1)} 
                                disabled={page >= totalPage}
                            >
                                Next
                            </button>
                        </div>
                    </div>
                </>
            )}

            {/* Customer Details & Edit Modal */}
            {selectedCust && (
                <div className="modal-backdrop">
                    <div className="custom-modal wide-modal text-left">
                        <div className="modal-header">
                            <h3>Subscriber Profile: {selectedCust.user?.name} ({selectedCust.id})</h3>
                            <button onClick={() => setSelectedCust(null)} className="close-btn">&times;</button>
                        </div>
                        <div className="modal-body">
                            <form onSubmit={handleSaveChanges} className="detail-grid">
                                
                                {/* Left Column: Info & Edit Form */}
                                <div className="detail-section">
                                    <div className="flex-row justify-between items-center mb-4 border-b pb-2">
                                        <h4 className="m-0 border-none pb-0">Account Details</h4>
                                        {!isEditing ? (
                                            <button 
                                                type="button"
                                                onClick={() => setIsEditing(true)}
                                                className="btn btn-primary btn-sm"
                                            >
                                                <Edit size={14} /> Edit details
                                            </button>
                                        ) : (
                                            <div className="flex-row gap-2">
                                                <button 
                                                    type="submit"
                                                    className="btn btn-success btn-sm"
                                                    disabled={saving}
                                                >
                                                    <Save size={14} /> Save
                                                </button>
                                                <button 
                                                    type="button"
                                                    onClick={() => setIsEditing(false)}
                                                    className="btn btn-secondary btn-sm"
                                                >
                                                    <X size={14} /> Cancel
                                                </button>
                                            </div>
                                        )}
                                    </div>

                                    <div className="info-row">
                                        <span className="label">Status:</span>
                                        <span className="val flex-row gap-2 items-center">
                                            <span className={`badge-status status-${selectedCust.status}`}>
                                                {selectedCust.status.toUpperCase()}
                                            </span>
                                            {selectedCust.status === 'active' && (
                                                <button type="button" onClick={() => handleSuspend(selectedCust.id)} className="btn btn-warning btn-sm py-1">
                                                    Suspend
                                                </button>
                                            )}
                                            {selectedCust.status === 'suspended' && (
                                                <button type="button" onClick={() => handleUnsuspend(selectedCust.id)} className="btn btn-success btn-sm py-1">
                                                    Unsuspend
                                                </button>
                                            )}
                                            {selectedCust.status !== 'terminated' && (
                                                <button type="button" onClick={() => handleTerminate(selectedCust.id)} className="btn btn-danger btn-sm py-1">
                                                    Terminate
                                                </button>
                                            )}
                                        </span>
                                    </div>

                                    <div className="info-row">
                                        <span className="label">Internet Package:</span>
                                        <span className="val">
                                            {isEditing ? (
                                                <select 
                                                    value={formPackageId} 
                                                    onChange={(e) => setFormPackageId(e.target.value)}
                                                    className="form-input w-full text-white bg-dark"
                                                    required
                                                >
                                                    <option value="">Select Package</option>
                                                    {packages.map(p => (
                                                        <option key={p.id} value={p.id}>{p.name} ({p.speed_mbps} Mbps) - Rp {p.price.toLocaleString()}</option>
                                                    ))}
                                                </select>
                                            ) : (
                                                <span className="bold text-purple">
                                                    {selectedCust.package ? `${selectedCust.package.name} (${selectedCust.package.speed_mbps} Mbps) - Rp ${selectedCust.package.price.toLocaleString()}` : selectedCust.package_id}
                                                </span>
                                            )}
                                        </span>
                                    </div>

                                    <div className="info-row">
                                        <span className="label">Nomor ODP (Port):</span>
                                        <span className="val">
                                            {isEditing ? (
                                                <input 
                                                    type="text" 
                                                    value={formOdpNumber} 
                                                    onChange={(e) => setFormOdpNumber(e.target.value)}
                                                    className="form-input w-full text-white bg-dark"
                                                    placeholder="ODP Number, e.g. ODP-MDN-03"
                                                />
                                            ) : (
                                                <span className="bold text-success">{selectedCust.odp_number || '-'}</span>
                                            )}
                                        </span>
                                    </div>

                                    <div className="info-row">
                                        <span className="label">Billing Day:</span>
                                        <span className="val">
                                            {isEditing ? (
                                                <input 
                                                    type="number" 
                                                    min="1" 
                                                    max="28"
                                                    value={formDueDateDay} 
                                                    onChange={(e) => setFormDueDateDay(e.target.value)}
                                                    className="form-input w-full text-white bg-dark"
                                                    required
                                                />
                                            ) : (
                                                <span>Day {selectedCust.due_date_day} (Setiap bulan)</span>
                                            )}
                                        </span>
                                    </div>

                                    <div className="info-row">
                                        <span className="label">Router Assignment:</span>
                                        <span className="val">
                                            {isEditing ? (
                                                <select 
                                                    value={formRouterId} 
                                                    onChange={(e) => setFormRouterId(e.target.value)}
                                                    className="form-input w-full text-white bg-dark"
                                                >
                                                    <option value="">None / Auto Assign</option>
                                                    {routers.map(r => (
                                                        <option key={r.id} value={r.id}>{r.name} ({r.host})</option>
                                                    ))}
                                                </select>
                                            ) : (
                                                <span>{selectedCust.router ? `${selectedCust.router.name} (${selectedCust.router.host})` : 'Unassigned'}</span>
                                            )}
                                        </span>
                                    </div>

                                    <div className="info-row">
                                        <span className="label">PPP / Radius Username:</span>
                                        <span className="val">
                                            {isEditing ? (
                                                <input 
                                                    type="text" 
                                                    value={formPppUsername} 
                                                    onChange={(e) => setFormPppUsername(e.target.value)}
                                                    className="form-input w-full text-white bg-dark"
                                                    required
                                                />
                                            ) : (
                                                <code>{selectedCust.ppp_username}</code>
                                            )}
                                        </span>
                                    </div>

                                    <div className="info-row">
                                        <span className="label">PPP / Radius Password:</span>
                                        <span className="val">
                                            {isEditing ? (
                                                <input 
                                                    type="text" 
                                                    value={formPppPassword} 
                                                    onChange={(e) => setFormPppPassword(e.target.value)}
                                                    className="form-input w-full text-white bg-dark"
                                                    required
                                                />
                                            ) : (
                                                <span className="font-mono">•••••••• (Secured)</span>
                                            )}
                                        </span>
                                    </div>

                                    <div className="info-row">
                                        <span className="label">Email:</span>
                                        <span className="val">{selectedCust.user?.email || '-'}</span>
                                    </div>
                                    
                                    <div className="info-row">
                                        <span className="label">Company Name:</span>
                                        <span className="val">{selectedCust.user?.company_name || 'GREENET'}</span>
                                    </div>
                                </div>

                                {/* Right Column: Registration & Documents */}
                                <div className="detail-section doc-uploads-section">
                                    <h4>Installation & KTP Verification</h4>
                                    
                                    {selectedCust.registration ? (
                                        <>
                                            <div className="info-row">
                                                <span className="label">Provinsi:</span>
                                                <span className="val">{selectedCust.registration.province || '-'}</span>
                                            </div>
                                            <div className="info-row">
                                                <span className="label">Kota / Kabupaten:</span>
                                                <span className="val">{selectedCust.registration.city || '-'}</span>
                                            </div>
                                            <div className="info-row">
                                                <span className="label">Kecamatan:</span>
                                                <span className="val">{selectedCust.registration.district || '-'}</span>
                                            </div>
                                            <div className="info-row">
                                                <span className="label">Kelurahan / Desa:</span>
                                                <span className="val">{selectedCust.registration.village || '-'}</span>
                                            </div>
                                            <div className="info-row">
                                                <span className="label">Alamat Instalasi:</span>
                                                <span className="val text-sm">{selectedCust.registration.installation_address}</span>
                                            </div>
                                            <div className="info-row">
                                                <span className="label">No HP / WhatsApp:</span>
                                                <span className="val">{selectedCust.registration.phone}</span>
                                            </div>
                                            <div className="info-row">
                                                <span className="label">NIK / KTP Number:</span>
                                                <span className="val">{selectedCust.registration.nik}</span>
                                            </div>
                                            <div className="info-row">
                                                <span className="label">Map Pin Location:</span>
                                                <span className="val">
                                                    {selectedCust.registration.latitude && selectedCust.registration.longitude ? (
                                                        <a 
                                                            href={`https://www.google.com/maps?q=${selectedCust.registration.latitude},${selectedCust.registration.longitude}`} 
                                                            target="_blank" 
                                                            rel="noreferrer"
                                                            className="map-link-btn"
                                                        >
                                                            <MapPin size={12} /> View Map Location
                                                        </a>
                                                    ) : 'N/A'}
                                                </span>
                                            </div>
                                            <div className="doc-item mt-4">
                                                <span className="doc-label"><FileImage size={14} /> Registered KTP Image</span>
                                                {selectedCust.registration.ktp_path ? (
                                                    <div className="image-container">
                                                        <a href={selectedCust.registration.ktp_path} target="_blank" rel="noreferrer">
                                                            <img src={selectedCust.registration.ktp_path} alt="KTP Photo" className="doc-preview-img" />
                                                        </a>
                                                    </div>
                                                ) : <p className="no-doc">No KTP Photo available</p>}
                                            </div>
                                        </>
                                    ) : (
                                        <div className="empty-state text-sm p-4">
                                            No original registration record linked to this customer account.
                                        </div>
                                    )}
                                </div>

                            </form>
                        </div>
                    </div>
                </div>
            )}

            {/* History Modal */}
            {selectedHistory && (
                <div className="modal-backdrop">
                    <div className="custom-modal text-left">
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
