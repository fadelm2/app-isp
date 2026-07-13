import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { CreditCard, CheckCircle, AlertTriangle, Search, Loader2, ArrowRight } from 'lucide-react';
import './PublicPayment.css';

// Dynamically fetch configurations
let API_BASE_URL = '/api';
let CLIENT_KEY = 'SB-Mid-client-cGrC44IfUPvwlZ3r'; // fallback default

export function PublicPayment() {
	const { customerId: urlCustomerId } = useParams();
	const navigate = useNavigate();
	const [customerId, setCustomerId] = useState(urlCustomerId || '');
	const [invoices, setInvoices] = useState([]);
	const [loading, setLoading] = useState(false);
	const [payingId, setPayingId] = useState(null);
	const [error, setError] = useState('');
	const [searchVal, setSearchVal] = useState('');

	// Load Midtrans Snap Script
	useEffect(() => {
		const loadSnap = async () => {
			try {
				const response = await axios.get('/config.json');
				const config = response.data;
				const snapUrl = config.payment?.midtrans?.is_production 
					? 'https://app.midtrans.com/snap/snap.js'
					: 'https://app.sandbox.midtrans.com/snap/snap.js';
				const clientKey = config.payment?.midtrans?.client_key || CLIENT_KEY;
				
				// Inject script
				const script = document.createElement('script');
				script.src = snapUrl;
				script.setAttribute('data-client-key', clientKey);
				script.async = true;
				document.body.appendChild(script);
			} catch (e) {
				const script = document.createElement('script');
				script.src = 'https://app.sandbox.midtrans.com/snap/snap.js';
				script.setAttribute('data-client-key', CLIENT_KEY);
				script.async = true;
				document.body.appendChild(script);
			}
		};
		loadSnap();
	}, []);

	useEffect(() => {
		if (urlCustomerId) {
			fetchInvoices(urlCustomerId);
		}
	}, [urlCustomerId]);

	const fetchInvoices = async (id) => {
		setLoading(true);
		setError('');
		try {
			const res = await axios.get(`${API_BASE_URL}/public/customers/${id}/invoices`);
			setInvoices(res.data.data || []);
		} catch (err) {
			const errMsg = err.response?.data?.errors || err.response?.data?.message || 'Failed to fetch bills. Check Customer ID.';
			setError(errMsg);
			setInvoices([]);
		} finally {
			setLoading(false);
		}
	};

	const handleSearchSubmit = (e) => {
		e.preventDefault();
		if (searchVal.trim()) {
			navigate(`/payment/${searchVal.trim()}`);
		}
	};

	const handlePay = async (invoiceId) => {
		setPayingId(invoiceId);
		setError('');
		try {
			const res = await axios.get(`${API_BASE_URL}/public/invoices/${invoiceId}/pay`);
			const snapToken = res.data.data;
			
			if (window.snap) {
				window.snap.pay(snapToken, {
					onSuccess: () => {
						alert('Payment successful!');
						if (urlCustomerId) fetchInvoices(urlCustomerId);
					},
					onPending: () => {
						alert('Payment is pending. Please complete transaction.');
						if (urlCustomerId) fetchInvoices(urlCustomerId);
					},
					onError: () => {
						setError('Payment process failed. Please try again.');
					},
					onClose: () => {
						console.log('Payment modal closed');
					}
				});
			} else {
				setError('Payment gateway is loading, please try again in a few seconds.');
			}
		} catch (err) {
			const errMsg = err.response?.data?.errors || 'Error initiating payment gateway.';
			setError(errMsg);
		} finally {
			setPayingId(null);
		}
	};

	return (
		<div className="pub-payment-container">
			<div className="pub-payment-card">
				<div className="pub-payment-header">
					<div className="brand-logo">GREENET</div>
					<h2>Quick Bill Payment Portal</h2>
					<p>Pay your internet service invoices instantly using safe online payment channels</p>
				</div>

				<form onSubmit={handleSearchSubmit} className="search-form-pub">
					<div className="search-input-wrapper">
						<Search size={18} />
						<input 
							type="text" 
							placeholder="Enter Customer ID (e.g. CUST-12345)" 
							value={searchVal}
							onChange={(e) => setSearchVal(e.target.value)}
							required
						/>
						<button type="submit" disabled={loading}>
							{loading ? <Loader2 className="animate-spin" size={16} /> : <ArrowRight size={16} />}
						</button>
					</div>
				</form>

				{error && (
					<div className="alert-error-pub">
						<AlertTriangle size={18} />
						<span>{error}</span>
					</div>
				)}

				{loading ? (
					<div className="pub-payment-loading">
						<Loader2 className="animate-spin" size={32} />
						<p>Fetching billing records...</p>
					</div>
				) : (
					<div className="pub-invoice-list">
						{urlCustomerId && invoices.length === 0 && !error && (
							<div className="no-bills-card">
								<CheckCircle size={40} className="success-icon" />
								<h3>No Outstanding Bills</h3>
								<p>Customer <strong>{urlCustomerId}</strong> is completely up-to-date. Thank you!</p>
							</div>
						)}

						{invoices.map((inv) => (
							<div key={inv.id} className="pub-invoice-card">
								<div className="pub-invoice-header">
									<div>
										<span className="badge-period">Period: {inv.period_month}/{inv.period_year}</span>
										<h4>Invoice ID: {inv.id}</h4>
									</div>
									<div className="status-badge pending">Owed</div>
								</div>
								
								<div className="pub-invoice-body">
									<div className="info-row">
										<span>Customer Name:</span>
										<strong>{inv.customer?.user?.username || 'Customer'}</strong>
									</div>
									<div className="info-row">
										<span>Internet Plan:</span>
										<strong>{inv.customer?.package?.name || 'Active Package'}</strong>
									</div>
									<div className="info-row border-t pt-2">
										<span>Monthly Charge:</span>
										<span>{inv.amount.toLocaleString()} IDR</span>
									</div>
									<div className="info-row">
										<span>Tax (11%):</span>
										<span>{inv.tax_amount.toLocaleString()} IDR</span>
									</div>
									{inv.installation_fee > 0 && (
										<div className="info-row">
											<span>Installation Fee:</span>
											<span>{inv.installation_fee.toLocaleString()} IDR</span>
										</div>
									)}
									<div className="info-row total-row">
										<span>Total Amount:</span>
										<span className="grand-total">{inv.total_amount.toLocaleString()} IDR</span>
									</div>
								</div>

								<button 
									className="pub-pay-btn" 
									onClick={() => handlePay(inv.id)}
									disabled={payingId === inv.id}
								>
									{payingId === inv.id ? (
										<>
											<Loader2 className="animate-spin" size={16} />
											Processing Gateway...
										</>
									) : (
										<>
											<CreditCard size={16} />
											Pay Bill Now
										</>
									)}
								</button>
							</div>
						))}
					</div>
				)}
			</div>
		</div>
	);
}
