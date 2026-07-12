import React, { useEffect, useState } from 'react';
import { adminService, customerService } from '../services/api';
import { Plus, CreditCard, RefreshCw } from 'lucide-react';
import './Invoices.css';
import './Pages.css';

export const Invoices = () => {
    const [invoices, setInvoices] = useState([]);
    const [loading, setLoading] = useState(true);
    const [generating, setGenerating] = useState(false);
    const [showGenForm, setShowGenForm] = useState(false);
    const [genData, setGenData] = useState({
        customer_id: '',
        month: new Date().getMonth() + 1,
        year: new Date().getFullYear()
    });

    const fetchInvoices = async () => {
        try {
            const response = await adminService.getInvoices();
            if (response.data && response.data.data) {
                setInvoices(response.data.data);
            }
        } catch (error) {
            console.error("Error fetching invoices", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchInvoices();
    }, []);

    // Load Midtrans Snap Script dynamically if paying in frontend
    useEffect(() => {
        const snapSrc = "https://app.sandbox.midtrans.com/snap/snap.js";
        const clientKey = "SB-Mid-client-cGrC44IfUPvwlZ3r"; // Sandbox client key
        
        const script = document.createElement('script');
        script.src = snapSrc;
        script.setAttribute('data-client-key', clientKey);
        document.body.appendChild(script);

        return () => {
            document.body.removeChild(script);
        };
    }, []);

    const handleCreateInvoice = async (e) => {
        e.preventDefault();
        setGenerating(true);
        try {
            const body = {
                customer_id: genData.customer_id,
                period_month: parseInt(genData.month),
                period_year: parseInt(genData.year)
            };
            await adminService.createInvoice(body);
            setShowGenForm(false);
            setGenData({ ...genData, customer_id: '' });
            fetchInvoices();
        } catch (error) {
            alert("Failed to generate: " + (error.response?.data?.error || error.message));
        } finally {
            setGenerating(false);
        }
    };

    const handlePay = async (invoiceId) => {
        try {
            const response = await customerService.getSnapToken(invoiceId);
            const token = response.data.data;
            
            if (window.snap) {
                window.snap.pay(token, {
                    onSuccess: function(result) {
                        alert("Payment successful!");
                        fetchInvoices();
                    },
                    onPending: function(result) {
                        alert("Payment pending. Please complete transaction.");
                    },
                    onError: function(result) {
                        alert("Payment failed!");
                    },
                    onClose: function() {
                        console.log("Payment pop-up closed");
                    }
                });
            } else {
                alert("Midtrans payment library loading... please retry.");
            }
        } catch (error) {
            alert("Failed to fetch payment Snap Token.");
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
                <h2>Billing & Invoices</h2>
                <button onClick={() => setShowGenForm(true)} className="btn btn-primary" disabled={generating}>
                    <Plus size={16} /> Generate Invoice
                </button>
            </div>

            {loading ? (
                <div className="loading">Loading invoices...</div>
            ) : invoices.length === 0 ? (
                <div className="empty-state">No invoices found. Click 'Generate Invoice' to create billing records.</div>
            ) : (
                <div className="table-responsive">
                    <table className="custom-table">
                        <thead>
                            <tr>
                                <th>Invoice ID</th>
                                <th>Subscriber</th>
                                <th>Billing Period</th>
                                <th>Amount</th>
                                <th>Tax (11%)</th>
                                <th>Setup Fee</th>
                                <th>Total Bill</th>
                                <th>Status</th>
                                <th>Payment</th>
                            </tr>
                        </thead>
                        <tbody>
                            {invoices.map((inv) => (
                                <tr key={inv.id}>
                                    <td className="bold">{inv.id}</td>
                                    <td>
                                        <div className="bold">{inv.customer ? inv.customer.user.name : inv.customer_id}</div>
                                        <div className="text-muted small">Due: {new Date(inv.due_date).toLocaleDateString()}</div>
                                    </td>
                                    <td>{inv.period_month} / {inv.period_year}</td>
                                    <td>{formatRupiah(inv.amount)}</td>
                                    <td>{formatRupiah(inv.tax_amount)}</td>
                                    <td>{formatRupiah(inv.installation_fee)}</td>
                                    <td className="bold text-purple">{formatRupiah(inv.total_amount)}</td>
                                    <td>
                                        <span className={`badge-status status-${inv.status}`}>
                                            {inv.status.toUpperCase()}
                                        </span>
                                    </td>
                                    <td>
                                        {inv.status === 'pending' || inv.status === 'owed' ? (
                                            <button 
                                                onClick={() => handlePay(inv.id)} 
                                                className="btn btn-success btn-sm btn-pay"
                                            >
                                                <CreditCard size={14} /> Pay Sandbox
                                            </button>
                                        ) : inv.status === 'paid' ? (
                                            <span className="text-muted text-success small font-medium">Paid ({new Date(inv.paid_at || Date.now()).toLocaleDateString()})</span>
                                        ) : (
                                            <span className="text-muted small">N/A</span>
                                        )}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Gen Invoice Modal */}
            {showGenForm && (
                <div className="modal-backdrop">
                    <div className="custom-modal">
                        <div className="modal-header">
                            <h3>Generate Subscriber Invoice</h3>
                            <button onClick={() => setShowGenForm(false)} className="close-btn">&times;</button>
                        </div>
                        <form onSubmit={handleCreateInvoice} className="modal-form">
                            <div className="form-group">
                                <label>Subscriber ID (e.g. CUST-01234)</label>
                                <input 
                                    type="text" 
                                    value={genData.customer_id} 
                                    onChange={(e) => setGenData({ ...genData, customer_id: e.target.value })} 
                                    required 
                                    placeholder="Enter Customer ID"
                                />
                            </div>
                            <div className="form-group">
                                <label>Period Month</label>
                                <select 
                                    value={genData.month} 
                                    onChange={(e) => setGenData({ ...genData, month: e.target.value })}
                                >
                                    {[...Array(12)].map((_, i) => (
                                        <option key={i+1} value={i+1}>{i+1}</option>
                                    ))}
                                </select>
                            </div>
                            <div className="form-group">
                                <label>Period Year</label>
                                <input 
                                    type="number" 
                                    value={genData.year} 
                                    onChange={(e) => setGenData({ ...genData, year: e.target.value })} 
                                    required
                                />
                            </div>
                            <div className="form-actions">
                                <button type="button" onClick={() => setShowGenForm(false)} className="btn btn-secondary">Cancel</button>
                                <button type="submit" className="btn btn-primary" disabled={generating}>
                                    {generating ? 'Generating...' : 'Create Invoice'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};
