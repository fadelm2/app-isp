import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { adminService } from '../services/api';
import { ClipboardList, Check, X, Eye, MapPin, Loader, Info, FileImage } from 'lucide-react';
import './Pages.css';

export const Registrations = () => {
    const [registrations, setRegistrations] = useState([]);
    const [loading, setLoading] = useState(true);
    const [updatingId, setUpdatingId] = useState(null);
    const [selectedReg, setSelectedReg] = useState(null);
    const [odpNumber, setOdpNumber] = useState('');
    const navigate = useNavigate();

    // Search and Pagination States
    const [search, setSearch] = useState('');
    const [status, setStatus] = useState('');
    const [page, setPage] = useState(1);
    const [size, setSize] = useState(10);
    const [totalItem, setTotalItem] = useState(0);
    const [totalPage, setTotalPage] = useState(0);

    const fetchRegistrations = async (currentPage = page, searchKeyword = search, statusFilter = status, pageSize = size) => {
        setLoading(true);
        try {
            const response = await adminService.getRegistrations({
                search: searchKeyword,
                status: statusFilter,
                page: currentPage,
                size: pageSize
            });
            if (response.data && response.data.data) {
                setRegistrations(response.data.data);
                if (response.data.paging) {
                    setTotalItem(response.data.paging.total_item);
                    setTotalPage(response.data.paging.total_page);
                }
            }
        } catch (error) {
            console.error("Error fetching registrations", error);
        } finally {
            setLoading(false);
        }
    };

    // Debounced search trigger
    useEffect(() => {
        const delayDebounceFn = setTimeout(() => {
            setPage(1);
            fetchRegistrations(1, search, status, size);
        }, 300);

        return () => clearTimeout(delayDebounceFn);
    }, [search]);

    // Direct filter triggers
    const handleStatusChange = (newStatus) => {
        setStatus(newStatus);
        setPage(1);
        fetchRegistrations(1, search, newStatus, size);
    };

    const handleSizeChange = (newSize) => {
        setSize(newSize);
        setPage(1);
        fetchRegistrations(1, search, status, newSize);
    };

    const handlePageChange = (newPage) => {
        setPage(newPage);
        fetchRegistrations(newPage, search, status, size);
    };

    const handleUpdateStatus = async (id, statusVal, odp = '') => {
        setUpdatingId(id);
        try {
            await adminService.updateRegistrationStatus(id, statusVal, odp);
            fetchRegistrations(page, search, status, size);
            if (selectedReg && selectedReg.id === id) {
                setSelectedReg(prev => ({ ...prev, status: statusVal, odp_number: odp || prev.odp_number }));
            }
            if (statusVal === 'approved') {
                setSelectedReg(null);
            }
        } catch (error) {
            console.error("Failed to update status", error);
            alert("Error: " + (error.response?.data?.error || "Failed to update registration"));
        } finally {
            setUpdatingId(null);
        }
    };

    const handleOpenDetail = (reg) => {
        setSelectedReg(reg);
        setOdpNumber(reg.odp_number || '');
    };

    return (
        <div className="page-container">
            <div className="page-header">
                <h2>Customer Registrations</h2>
                <button onClick={() => navigate('/registrations/new')} className="btn btn-primary">
                    + Register Customer
                </button>
            </div>

            {/* Search & Filter Bar */}
            <div className="filter-bar mb-4" style={{ display: 'flex', gap: '12px', flexWrap: 'wrap', alignItems: 'center', marginBottom: '16px' }}>
                <input 
                    type="text" 
                    placeholder="Search by ID, name, email, phone, NIK, ODP..." 
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
                    <option value="pending">Pending</option>
                    <option value="under_review">Under Review</option>
                    <option value="surveying">Surveying</option>
                    <option value="approved">Approved</option>
                    <option value="rejected">Rejected</option>
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
                <div className="loading">Loading registrations...</div>
            ) : registrations.length === 0 ? (
                <div className="empty-state">No customer registrations found.</div>
            ) : (
                <>
                    <div className="table-responsive">
                        <table className="custom-table">
                            <thead>
                                <tr>
                                    <th>Name / Address</th>
                                    <th>NIK</th>
                                    <th>Email / Phone</th>
                                    <th>Package</th>
                                    <th>ODP Number</th>
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
                                            <div className="text-muted small">{reg.phone}</div>
                                        </td>
                                        <td>{reg.package ? reg.package.name : reg.package_id}</td>
                                        <td>
                                            <span className="font-mono text-xs">{reg.odp_number || '-'}</span>
                                        </td>
                                        <td>
                                            <span className={`badge-status status-${reg.status}`}>
                                                {reg.status.toUpperCase()}
                                            </span>
                                        </td>
                                        <td>
                                            <div className="action-buttons">
                                                <button 
                                                    onClick={() => handleOpenDetail(reg)}
                                                    className="btn btn-secondary btn-sm"
                                                    title="View Details"
                                                >
                                                    <Eye size={14} /> Detail
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
                            Showing {registrations.length} of {totalItem} registrations
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

            {/* Registration Details Modal */}
            {selectedReg && (
                <div className="modal-backdrop">
                    <div className="custom-modal wide-modal">
                        <div className="modal-header">
                            <h3>Registration Details: {selectedReg.full_name}</h3>
                            <button onClick={() => setSelectedReg(null)} className="close-btn">&times;</button>
                        </div>
                        <div className="modal-body text-left">
                            <div className="detail-grid">
                                {/* Left Column: Text Information */}
                                <div className="detail-section text-info-section">
                                    <h4>Subscriber Information</h4>
                                    <div className="info-row">
                                        <span className="label">Status:</span>
                                        <span className={`badge-status status-${selectedReg.status}`}>
                                            {selectedReg.status.toUpperCase()}
                                        </span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Full Name:</span>
                                        <span className="val bold">{selectedReg.full_name}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">NIK:</span>
                                        <span className="val">{selectedReg.nik}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Birth Info:</span>
                                        <span className="val">{selectedReg.birth_place}, {selectedReg.birth_date}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Gender:</span>
                                        <span className="val">{selectedReg.gender}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Email:</span>
                                        <span className="val">{selectedReg.email}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Phone:</span>
                                        <span className="val">{selectedReg.phone}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Internet Package:</span>
                                        <span className="val bold text-purple">
                                            {selectedReg.package ? `${selectedReg.package.name} (${selectedReg.package.speed_mbps} Mbps) - Rp ${selectedReg.package.price.toLocaleString()}` : selectedReg.package_id}
                                        </span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Provinsi:</span>
                                        <span className="val">{selectedReg.province || '-'}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Kota / Kabupaten:</span>
                                        <span className="val">{selectedReg.city || '-'}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Kecamatan:</span>
                                        <span className="val">{selectedReg.district || '-'}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Kelurahan / Desa:</span>
                                        <span className="val">{selectedReg.village || '-'}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Alamat Instalasi:</span>
                                        <span className="val text-sm">{selectedReg.installation_address}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Alamat Penagihan:</span>
                                        <span className="val text-sm">{selectedReg.billing_address}</span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">GPS Coordinates:</span>
                                        <span className="val">
                                            {selectedReg.latitude && selectedReg.longitude ? (
                                                <a 
                                                    href={`https://www.google.com/maps?q=${selectedReg.latitude},${selectedReg.longitude}`} 
                                                    target="_blank" 
                                                    rel="noreferrer"
                                                    className="map-link-btn"
                                                >
                                                    <MapPin size={14} /> {selectedReg.latitude.toFixed(6)}, {selectedReg.longitude.toFixed(6)}
                                                </a>
                                            ) : 'N/A'}
                                        </span>
                                    </div>
                                    <div className="info-row">
                                        <span className="label">Notes / Description:</span>
                                        <span className="val italic">{selectedReg.notes || '-'}</span>
                                    </div>

                                    {/* Action Workflows inside Modal */}
                                    <div className="modal-actions-workflow">
                                        <h4>Workflow Actions</h4>
                                        {selectedReg.status === 'pending' && (
                                            <button 
                                                onClick={() => handleUpdateStatus(selectedReg.id, 'under_review')}
                                                className="btn btn-secondary w-full"
                                                disabled={updatingId !== null}
                                            >
                                                Start Review (Under Review)
                                            </button>
                                        )}

                                        {selectedReg.status === 'under_review' && (
                                            <div className="workflow-action-group">
                                                <button 
                                                    onClick={() => handleUpdateStatus(selectedReg.id, 'surveying')}
                                                    className="btn btn-warning w-full"
                                                    disabled={updatingId !== null}
                                                >
                                                    Assign to Survey Phase
                                                </button>
                                            </div>
                                        )}

                                        {selectedReg.status === 'surveying' && (
                                            <div className="workflow-action-group">
                                                <div className="form-group-inline">
                                                    <label>Nomor ODP (Optical Distribution Point)</label>
                                                    <input 
                                                        type="text" 
                                                        placeholder="Contoh: ODP-JKT-01A"
                                                        value={odpNumber} 
                                                        onChange={(e) => setOdpNumber(e.target.value)}
                                                        className="form-input text-white bg-dark"
                                                    />
                                                </div>
                                                <div className="flex-row gap-2 mt-2">
                                                    <button 
                                                        onClick={() => handleUpdateStatus(selectedReg.id, 'approved', odpNumber)}
                                                        className="btn btn-success flex-1"
                                                        disabled={updatingId !== null || !odpNumber}
                                                    >
                                                        <Check size={14} /> Approve & Activate
                                                    </button>
                                                    <button 
                                                        onClick={() => handleUpdateStatus(selectedReg.id, 'rejected')}
                                                        className="btn btn-danger flex-1"
                                                        disabled={updatingId !== null}
                                                    >
                                                        <X size={14} /> Reject Application
                                                    </button>
                                                </div>
                                                {!odpNumber && <p className="text-warning text-xs mt-1">* Mohon masukkan nomor ODP sebelum menyetujui pendaftaran.</p>}
                                            </div>
                                        )}

                                        {selectedReg.status === 'approved' && (
                                            <div className="info-row">
                                                <span className="label">ODP Port assigned:</span>
                                                <span className="val bold text-success">{selectedReg.odp_number || odpNumber || 'Verified'}</span>
                                            </div>
                                        )}

                                        {selectedReg.status === 'rejected' && (
                                            <div className="info-row">
                                                <span className="label">Registration Status:</span>
                                                <span className="val bold text-danger">REJECTED</span>
                                            </div>
                                        )}
                                    </div>
                                </div>

                                {/* Right Column: Document Upload Views */}
                                <div className="detail-section doc-uploads-section">
                                    <h4>Verification Documents</h4>
                                    
                                    <div className="doc-item">
                                        <span className="doc-label"><FileImage size={14} /> KTP Photo</span>
                                        {selectedReg.ktp_path ? (
                                            <div className="image-container">
                                                <a href={selectedReg.ktp_path} target="_blank" rel="noreferrer">
                                                    <img src={selectedReg.ktp_path} alt="KTP Verification" className="doc-preview-img" />
                                                </a>
                                            </div>
                                        ) : <p className="no-doc">No KTP Photo Uploaded</p>}
                                    </div>

                                    {selectedReg.selfie_path && (
                                        <div className="doc-item">
                                            <span className="doc-label"><FileImage size={14} /> Selfie Photo</span>
                                            <div className="image-container">
                                                <a href={selectedReg.selfie_path} target="_blank" rel="noreferrer">
                                                    <img src={selectedReg.selfie_path} alt="Selfie Verification" className="doc-preview-img" />
                                                </a>
                                            </div>
                                        </div>
                                    )}

                                    {selectedReg.house_path && (
                                        <div className="doc-item">
                                            <span className="doc-label"><FileImage size={14} /> House / Site Photo</span>
                                            <div className="image-container">
                                                <a href={selectedReg.house_path} target="_blank" rel="noreferrer">
                                                    <img src={selectedReg.house_path} alt="House Front" className="doc-preview-img" />
                                                </a>
                                            </div>
                                        </div>
                                    )}

                                    {selectedReg.installation_path && (
                                        <div className="doc-item">
                                            <span className="doc-label"><FileImage size={14} /> Installation Site Detail</span>
                                            <div className="image-container">
                                                <a href={selectedReg.installation_path} target="_blank" rel="noreferrer">
                                                    <img src={selectedReg.installation_path} alt="Installation Site Detail" className="doc-preview-img" />
                                                </a>
                                            </div>
                                        </div>
                                    )}

                                    {selectedReg.supporting_doc_path && (
                                        <div className="doc-item">
                                            <span className="doc-label"><FileImage size={14} /> Supporting Document</span>
                                            <div className="doc-file-link">
                                                <a href={selectedReg.supporting_doc_path} target="_blank" rel="noreferrer" className="btn btn-secondary btn-sm">
                                                    Open Document File
                                                </a>
                                            </div>
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};
