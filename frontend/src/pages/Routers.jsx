import React, { useEffect, useState } from 'react';
import { adminService } from '../services/api';
import { Plus, Trash, Server, Activity } from 'lucide-react';
import './Pages.css';

export const Routers = () => {
    const [routers, setRouters] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [newRouter, setNewRouter] = useState({
        name: '',
        host: '',
        port: 8728,
        username: '',
        password: ''
    });

    const fetchRouters = async () => {
        try {
            const response = await adminService.getRouters();
            if (response.data && response.data.data) {
                setRouters(response.data.data);
            }
        } catch (error) {
            console.error("Error fetching routers", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchRouters();
    }, []);

    const handleCreate = async (e) => {
        e.preventDefault();
        try {
            const body = {
                name: newRouter.name,
                host: newRouter.host,
                port: parseInt(newRouter.port),
                username: newRouter.username,
                password: newRouter.password
            };
            await adminService.createRouter(body);
            setShowCreateForm(false);
            setNewRouter({ name: '', host: '', port: 8728, username: '', password: '' });
            fetchRouters();
        } catch (error) {
            alert("Failed to save router config: " + (error.response?.data?.error || error.message));
        }
    };

    const handleDelete = async (id) => {
        if (!confirm("Are you sure you want to delete this router configuration?")) return;
        try {
            await adminService.deleteRouter(id);
            fetchRouters();
        } catch (error) {
            alert("Failed to delete router config.");
        }
    };

    return (
        <div className="page-container">
            <div className="page-header">
                <h2>MikroTik NAS Routers</h2>
                <button onClick={() => setShowCreateForm(true)} className="btn btn-primary">
                    <Plus size={16} /> Add Router
                </button>
            </div>

            {loading ? (
                <div className="loading">Loading router configs...</div>
            ) : routers.length === 0 ? (
                <div className="empty-state">No MikroTik NAS routers defined. Click 'Add Router' to integrate.</div>
            ) : (
                <div className="table-responsive">
                    <table className="custom-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Host / IP Address</th>
                                <th>API Port</th>
                                <th>API Username</th>
                                <th>Status</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {routers.map((router) => (
                                <tr key={router.id}>
                                    <td className="bold">
                                        <div className="bold-flex">
                                            <Server size={16} className="text-secondary mr-2" />
                                            {router.name}
                                        </div>
                                    </td>
                                    <td><code>{router.host}</code></td>
                                    <td>{router.port}</td>
                                    <td><code>{router.username}</code></td>
                                    <td>
                                        <span className={`badge-status status-${router.status}`}>
                                            {router.status.toUpperCase()}
                                        </span>
                                    </td>
                                    <td>
                                        <button onClick={() => handleDelete(router.id)} className="btn btn-danger btn-sm">
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
                            <h3>Add MikroTik NAS Router</h3>
                            <button onClick={() => setShowCreateForm(false)} className="close-btn">&times;</button>
                        </div>
                        <form onSubmit={handleCreate} className="modal-form">
                            <div className="form-group">
                                <label>Router Name</label>
                                <input 
                                    type="text" 
                                    value={newRouter.name} 
                                    onChange={(e) => setNewRouter({ ...newRouter, name: e.target.value })} 
                                    required 
                                    placeholder="e.g. Core Router Mikrotik"
                                />
                            </div>
                            <div className="form-group">
                                <label>Host / IP Address</label>
                                <input 
                                    type="text" 
                                    value={newRouter.host} 
                                    onChange={(e) => setNewRouter({ ...newRouter, host: e.target.value })} 
                                    required
                                    placeholder="e.g. 192.168.88.1"
                                />
                            </div>
                            <div className="form-group">
                                <label>API Port (REST Default is 443 / API is 8728)</label>
                                <input 
                                    type="number" 
                                    value={newRouter.port} 
                                    onChange={(e) => setNewRouter({ ...newRouter, port: e.target.value })} 
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>API Username</label>
                                <input 
                                    type="text" 
                                    value={newRouter.username} 
                                    onChange={(e) => setNewRouter({ ...newRouter, username: e.target.value })} 
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>API Password</label>
                                <input 
                                    type="password" 
                                    value={newRouter.password} 
                                    onChange={(e) => setNewRouter({ ...newRouter, password: e.target.value })} 
                                    required
                                />
                            </div>
                            <div className="form-actions">
                                <button type="button" onClick={() => setShowCreateForm(false)} className="btn btn-secondary">Cancel</button>
                                <button type="submit" className="btn btn-primary">Connect Router</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};
