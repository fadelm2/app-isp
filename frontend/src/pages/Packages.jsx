import React, { useEffect, useState } from 'react';
import { adminService } from '../services/api';
import { Plus, Trash, Check, X } from 'lucide-react';
import './Pages.css';

export const Packages = () => {
    const [packages, setPackages] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [newPkg, setNewPkg] = useState({
        name: '',
        speed_mbps: 10,
        price: 150000,
        installation_fee: 500000,
        tax_rate: 0.11
    });

    const fetchPackages = async () => {
        try {
            const response = await adminService.getPackages();
            if (response.data && response.data.data) {
                setPackages(response.data.data);
            }
        } catch (error) {
            console.error("Error fetching packages", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPackages();
    }, []);

    const handleCreate = async (e) => {
        e.preventDefault();
        try {
            const body = {
                name: newPkg.name,
                speed_mbps: parseInt(newPkg.speed_mbps),
                price: parseFloat(newPkg.price),
                installation_fee: parseFloat(newPkg.installation_fee),
                tax_rate: parseFloat(newPkg.tax_rate)
            };
            await adminService.createPackage(body);
            setShowCreateForm(false);
            setNewPkg({ name: '', speed_mbps: 10, price: 150000, installation_fee: 500000, tax_rate: 0.11 });
            fetchPackages();
        } catch (error) {
            alert("Failed to create package: " + (error.response?.data?.error || error.message));
        }
    };

    const handleToggleActive = async (id, currentVal) => {
        try {
            await adminService.updatePackage(id, { is_active: !currentVal });
            fetchPackages();
        } catch (error) {
            alert("Failed to update status");
        }
    };

    const handleDelete = async (id) => {
        if (!confirm("Are you sure you want to delete this package?")) return;
        try {
            await adminService.deletePackage(id);
            fetchPackages();
        } catch (error) {
            alert("Failed to delete package. It might be assigned to a subscriber.");
        }
    };

    const formatRupiah = (val) => {
        return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0
        }).format(val);
    };

    return (
        <div className="page-container">
            <div className="page-header">
                <h2>Internet Packages</h2>
                <button onClick={() => setShowCreateForm(true)} className="btn btn-primary">
                    <Plus size={16} /> New Package
                </button>
            </div>

            {loading ? (
                <div className="loading">Loading packages...</div>
            ) : packages.length === 0 ? (
                <div className="empty-state">No packages defined. Click 'New Package' to add one.</div>
            ) : (
                <div className="table-responsive">
                    <table className="custom-table">
                        <thead>
                            <tr>
                                <th>Package Name</th>
                                <th>Speed Limit</th>
                                <th>Monthly Price</th>
                                <th>Installation Fee</th>
                                <th>Tax Rate</th>
                                <th>Status</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {packages.map((pkg) => (
                                <tr key={pkg.id}>
                                    <td className="bold">{pkg.name}</td>
                                    <td>{pkg.speed_mbps} Mbps</td>
                                    <td>{formatRupiah(pkg.price)}</td>
                                    <td>{formatRupiah(pkg.installation_fee)}</td>
                                    <td>{pkg.tax_rate * 100}%</td>
                                    <td>
                                        <button 
                                            onClick={() => handleToggleActive(pkg.id, pkg.is_active)}
                                            className={`badge-status ${pkg.is_active ? 'status-active' : 'status-suspended'}`}
                                            title="Click to toggle status"
                                        >
                                            {pkg.is_active ? 'ACTIVE' : 'INACTIVE'}
                                        </button>
                                    </td>
                                    <td>
                                        <button onClick={() => handleDelete(pkg.id)} className="btn btn-danger btn-sm">
                                            <Trash size={14} /> Delete
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Create Modal */}
            {showCreateForm && (
                <div className="modal-backdrop">
                    <div className="custom-modal">
                        <div className="modal-header">
                            <h3>Create Internet Package</h3>
                            <button onClick={() => setShowCreateForm(false)} className="close-btn">&times;</button>
                        </div>
                        <form onSubmit={handleCreate} className="modal-form">
                            <div className="form-group">
                                <label>Package Name</label>
                                <input 
                                    type="text" 
                                    value={newPkg.name} 
                                    onChange={(e) => setNewPkg({ ...newPkg, name: e.target.value })} 
                                    required 
                                    placeholder="e.g. Greenet Home 20M"
                                />
                            </div>
                            <div className="form-group">
                                <label>Speed (Mbps)</label>
                                <input 
                                    type="number" 
                                    value={newPkg.speed_mbps} 
                                    onChange={(e) => setNewPkg({ ...newPkg, speed_mbps: e.target.value })} 
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>Price (IDR / Month)</label>
                                <input 
                                    type="number" 
                                    value={newPkg.price} 
                                    onChange={(e) => setNewPkg({ ...newPkg, price: e.target.value })} 
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>Installation Fee (IDR)</label>
                                <input 
                                    type="number" 
                                    value={newPkg.installation_fee} 
                                    onChange={(e) => setNewPkg({ ...newPkg, installation_fee: e.target.value })} 
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>VAT / Tax Rate (Decimal, e.g. 0.11 for 11%)</label>
                                <input 
                                    type="number" 
                                    step="0.01" 
                                    value={newPkg.tax_rate} 
                                    onChange={(e) => setNewPkg({ ...newPkg, tax_rate: e.target.value })} 
                                    required
                                />
                            </div>
                            <div className="form-actions">
                                <button type="button" onClick={() => setShowCreateForm(false)} className="btn btn-secondary">Cancel</button>
                                <button type="submit" className="btn btn-primary">Create Package</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};
