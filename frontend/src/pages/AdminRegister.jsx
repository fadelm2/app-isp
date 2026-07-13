import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { publicService } from '../services/api';
import { Upload, FileCheck, CheckCircle2, ShieldAlert, ArrowLeft } from 'lucide-react';
import './Pages.css';

export const AdminRegister = () => {
    const [packages, setPackages] = useState([]);
    const [form, setForm] = useState({
        full_name: '',
        nik: '',
        birth_place: '',
        birth_date: '',
        gender: 'Laki-laki',
        email: '',
        phone: '',
        installation_address: '',
        billing_address: '',
        package_id: '',
        latitude: '0.0',
        longitude: '0.0',
        notes: ''
    });

    const [files, setFiles] = useState({
        ktp: null,
        house: null,
        installation: null,
        supporting_doc: null
    });

    const [loadingPackages, setLoadingPackages] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [submitted, setSubmitted] = useState(false);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const fetchPackages = async () => {
            try {
                const response = await publicService.getPackages();
                if (response.data && response.data.data) {
                    const activePkgs = response.data.data.filter(p => p.is_active);
                    setPackages(activePkgs);
                    if (activePkgs.length > 0) {
                        setForm(f => ({ ...f, package_id: activePkgs[0].id }));
                    }
                }
            } catch (err) {
                console.error("Failed to load packages", err);
            } finally {
                setLoadingPackages(false);
            }
        };
        fetchPackages();
    }, []);

    const handleFileChange = (e, key) => {
        if (e.target.files && e.target.files[0]) {
            setFiles(f => ({ ...f, [key]: e.target.files[0] }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSubmitting(true);

        // --- Data Validations ---

        // 1. NIK format check (exactly 16 digits)
        if (!/^\d{16}$/.test(form.nik)) {
            setError('NIK must be exactly 16 digits.');
            setSubmitting(false);
            return;
        }

        // 2. Birth Date check (cannot be in the future, must be min 17 years old)
        if (!form.birth_date) {
            setError('Birth Date is required.');
            setSubmitting(false);
            return;
        }
        const birthDateObj = new Date(form.birth_date);
        const today = new Date();
        if (birthDateObj > today) {
            setError('Birth Date cannot be in the future.');
            setSubmitting(false);
            return;
        }

        let age = today.getFullYear() - birthDateObj.getFullYear();
        const m = today.getMonth() - birthDateObj.getMonth();
        if (m < 0 || (m === 0 && today.getDate() < birthDateObj.getDate())) {
            age--;
        }
        if (age < 17) {
            setError('Customer must be at least 17 years old to register.');
            setSubmitting(false);
            return;
        }

        // 3. Email validation
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
            setError('Invalid email address format.');
            setSubmitting(false);
            return;
        }

        // 4. Phone validation (10 to 15 digits, allows optional + prefix)
        if (!/^\+?[0-9]{10,15}$/.test(form.phone)) {
            setError('Invalid Phone number (must be 10-15 digits).');
            setSubmitting(false);
            return;
        }

        // 5. Files validation
        if (!files.ktp) {
            setError('KTP Photo is required.');
            setSubmitting(false);
            return;
        }
        if (!files.installation) {
            setError('Installation site photo is required.');
            setSubmitting(false);
            return;
        }

        const formData = new FormData();
        Object.keys(form).forEach(key => {
            formData.append(key, form[key]);
        });
        Object.keys(files).forEach(key => {
            if (files[key]) {
                formData.append(key, files[key]);
            }
        });

        try {
            await publicService.register(formData);
            setSubmitted(true);
        } catch (err) {
            setError(err.response?.data?.errors || err.response?.data?.error || 'Failed to submit registration form.');
        } finally {
            setSubmitting(false);
        }
    };

    if (submitted) {
        return (
            <div className="page-container">
                <div className="empty-state">
                    <CheckCircle2 size={48} className="text-success mb-4" style={{ display: 'inline-block' }} />
                    <h3>Registration Submitted Successfully!</h3>
                    <p className="text-muted" style={{ margin: '12px 0 24px 0' }}>
                        The new customer registration has been recorded and is ready for verifications in the Registrations list.
                    </p>
                    <button onClick={() => navigate('/registrations')} className="btn btn-primary">
                        Back to Registrations
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="page-container">
            <div className="page-header">
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                    <button onClick={() => navigate('/registrations')} className="btn btn-secondary btn-icon" style={{ padding: '8px' }}>
                        <ArrowLeft size={16} />
                    </button>
                    <h2>Register New Subscriber</h2>
                </div>
            </div>

            {error && (
                <div className="pub-reg-error" style={{ marginBottom: '24px' }}>
                    <ShieldAlert size={20} />
                    <span>{error}</span>
                </div>
            )}

            <div className="table-responsive" style={{ padding: '32px' }}>
                <form onSubmit={handleSubmit} className="pub-reg-form">
                    <div className="form-section">
                        <h3>1. Subscriber Personal Details</h3>
                        <div className="form-row">
                            <div className="form-group-pub">
                                <label>Full Name (KTP Name)</label>
                                <input 
                                    type="text" 
                                    value={form.full_name} 
                                    onChange={(e) => setForm({ ...form, full_name: e.target.value })} 
                                    required 
                                />
                            </div>
                            <div className="form-group-pub">
                                <label>National Identity Number (NIK)</label>
                                <input 
                                    type="text" 
                                    placeholder="Must be 16 digits"
                                    value={form.nik} 
                                    onChange={(e) => setForm({ ...form, nik: e.target.value })} 
                                    required 
                                />
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group-pub">
                                <label>Birth Place</label>
                                <input 
                                    type="text" 
                                    value={form.birth_place} 
                                    onChange={(e) => setForm({ ...form, birth_place: e.target.value })} 
                                    required 
                                />
                            </div>
                            <div className="form-group-pub">
                                <label>Birth Date</label>
                                <input 
                                    type="date" 
                                    value={form.birth_date} 
                                    onChange={(e) => setForm({ ...form, birth_date: e.target.value })} 
                                    required 
                                />
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group-pub">
                                <label>Gender</label>
                                <select 
                                    value={form.gender} 
                                    onChange={(e) => setForm({ ...form, gender: e.target.value })}
                                >
                                    <option value="Laki-laki">Laki-laki</option>
                                    <option value="Perempuan">Perempuan</option>
                                </select>
                            </div>
                            <div className="form-group-pub">
                                <label>Phone / WhatsApp Number</label>
                                <input 
                                    type="tel" 
                                    placeholder="e.g. 08123456789"
                                    value={form.phone} 
                                    onChange={(e) => setForm({ ...form, phone: e.target.value })} 
                                    required 
                                />
                            </div>
                        </div>

                        <div className="form-group-pub">
                            <label>Email Address</label>
                            <input 
                                type="email" 
                                placeholder="name@example.com"
                                value={form.email} 
                                onChange={(e) => setForm({ ...form, email: e.target.value })} 
                                required 
                            />
                        </div>
                    </div>

                    <div className="form-section">
                        <h3>2. Installation & Service Plan</h3>
                        <div className="form-group-pub">
                            <label>Installation Address</label>
                            <textarea 
                                rows="3" 
                                value={form.installation_address} 
                                onChange={(e) => setForm({ ...form, installation_address: e.target.value })} 
                                required
                            />
                        </div>

                        <div className="form-group-pub">
                            <div className="flex-between">
                                <label>Billing Address</label>
                                <button 
                                    type="button" 
                                    className="link-btn-pub"
                                    onClick={() => setForm(f => ({ ...f, billing_address: f.installation_address }))}
                                >
                                    Same as Installation Address
                                </button>
                            </div>
                            <textarea 
                                rows="3" 
                                value={form.billing_address} 
                                onChange={(e) => setForm({ ...form, billing_address: e.target.value })} 
                                required
                            />
                        </div>

                        <div className="form-group-pub">
                            <label>Internet Service Package</label>
                            {loadingPackages ? (
                                <p className="small text-muted">Loading packages...</p>
                            ) : (
                                <select 
                                    value={form.package_id} 
                                    onChange={(e) => setForm({ ...form, package_id: e.target.value })}
                                >
                                    {packages.map(p => (
                                        <option key={p.id} value={p.id}>
                                            {p.name} - {p.speed_mbps} Mbps ({p.price.toLocaleString()} IDR/month)
                                        </option>
                                    ))}
                                </select>
                            )}
                        </div>

                        <div className="form-group-pub">
                            <label>Additional Notes (Optional)</label>
                            <textarea 
                                rows="2" 
                                value={form.notes} 
                                onChange={(e) => setForm({ ...form, notes: e.target.value })} 
                            />
                        </div>
                    </div>

                    <div className="form-section">
                        <h3>3. Verification Documents</h3>
                        
                        <div className="upload-grid-pub">
                            <div className="upload-item">
                                <label className="upload-label">
                                    <Upload size={20} />
                                    <span>KTP Photo *</span>
                                    <input 
                                        type="file" 
                                        accept="image/*,application/pdf"
                                        onChange={(e) => handleFileChange(e, 'ktp')} 
                                    />
                                </label>
                                {files.ktp && <div className="file-name"><FileCheck size={14} /> {files.ktp.name}</div>}
                            </div>

                            <div className="upload-item">
                                <label className="upload-label">
                                    <Upload size={20} />
                                    <span>House Photo (Optional)</span>
                                    <input 
                                        type="file" 
                                        accept="image/*"
                                        onChange={(e) => handleFileChange(e, 'house')} 
                                    />
                                </label>
                                {files.house && <div className="file-name"><FileCheck size={14} /> {files.house.name}</div>}
                            </div>

                            <div className="upload-item">
                                <label className="upload-label">
                                    <Upload size={20} />
                                    <span>Installation Site Photo *</span>
                                    <input 
                                        type="file" 
                                        accept="image/*"
                                        onChange={(e) => handleFileChange(e, 'installation')} 
                                    />
                                </label>
                                {files.installation && <div className="file-name"><FileCheck size={14} /> {files.installation.name}</div>}
                            </div>

                            <div className="upload-item">
                                <label className="upload-label">
                                    <Upload size={20} />
                                    <span>Other Supporting Doc (Optional)</span>
                                    <input 
                                        type="file" 
                                        accept="image/*,application/pdf"
                                        onChange={(e) => handleFileChange(e, 'supporting_doc')} 
                                    />
                                </label>
                                {files.supporting_doc && <div className="file-name"><FileCheck size={14} /> {files.supporting_doc.name}</div>}
                            </div>
                        </div>
                    </div>

                    <button type="submit" className="submit-btn-pub" disabled={submitting}>
                        {submitting ? 'Registering Customer...' : 'Submit New Registration'}
                    </button>
                </form>
            </div>
        </div>
    );
};
