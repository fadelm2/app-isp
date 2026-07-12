import React, { useState, useEffect } from 'react';
import { publicService } from '../services/api';
import { MapPin, Upload, FileCheck, CheckCircle2, ShieldAlert } from 'lucide-react';
import './PublicRegister.css';

export const PublicRegister = () => {
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
        latitude: '',
        longitude: '',
        notes: ''
    });

    const [files, setFiles] = useState({
        ktp: null,
        selfie: null,
        house: null,
        installation: null,
        supporting_doc: null
    });

    const [loadingPackages, setLoadingPackages] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [submitted, setSubmitted] = useState(false);
    const [error, setError] = useState('');

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
                console.error("Failed to load internet packages", err);
            } finally {
                setLoadingPackages(false);
            }
        };
        fetchPackages();
    }, []);

    const handleGetLocation = () => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                (position) => {
                    setForm(f => ({
                        ...f,
                        latitude: position.coords.latitude,
                        longitude: position.coords.longitude
                    }));
                },
                (err) => {
                    alert("Gagal mendapatkan lokasi GPS: " + err.message);
                }
            );
        } else {
            alert("Geolocation tidak didukung oleh browser Anda.");
        }
    };

    const handleFileChange = (e, key) => {
        if (e.target.files && e.target.files[0]) {
            setFiles(f => ({ ...f, [key]: e.target.files[0] }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSubmitting(true);

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
            console.error("Registration error", err);
            setError(err.response?.data?.errors || err.response?.data?.error || 'Gagal mengirim pendaftaran. Silakan periksa kembali berkas Anda.');
        } finally {
            setSubmitting(false);
        }
    };

    if (submitted) {
        return (
            <div className="pub-reg-wrapper">
                <div className="pub-reg-success-card">
                    <CheckCircle2 size={64} className="success-icon" />
                    <h2>Pendaftaran Berhasil!</h2>
                    <p className="success-message">
                        Terima kasih telah memilih <strong>GREENET</strong>. Data dan dokumen pendaftaran Anda telah berhasil dikirim dan sedang dalam proses review oleh tim verifikator kami.
                    </p>
                    <div className="next-steps">
                        <h4>Langkah Selanjutnya:</h4>
                        <ol>
                            <li>Tim kami akan menghubungi Anda melalui WhatsApp / Telepon untuk konfirmasi survei lokasi.</li>
                            <li>Teknisi melakukan instalasi kabel & pemasangan perangkat ONU.</li>
                            <li>Lakukan pembayaran tagihan pertama untuk mengaktifkan internet Anda.</li>
                        </ol>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="pub-reg-wrapper">
            <div className="pub-reg-card">
                <div className="pub-reg-header">
                    <h2>GREENET</h2>
                    <p className="subtitle">Formulir Pendaftaran Pelanggan Baru</p>
                </div>

                {error && (
                    <div className="pub-reg-error">
                        <ShieldAlert size={20} />
                        <span>{error}</span>
                    </div>
                )}

                <form onSubmit={handleSubmit} className="pub-reg-form">
                    {/* Section 1: Personal Info */}
                    <div className="form-section">
                        <h3>1. Data Diri Pelanggan</h3>
                        <div className="form-row">
                            <div className="form-group-pub">
                                <label>Nama Lengkap (Sesuai KTP)</label>
                                <input 
                                    type="text" 
                                    value={form.full_name} 
                                    onChange={(e) => setForm({ ...form, full_name: e.target.value })} 
                                    required 
                                />
                            </div>
                            <div className="form-group-pub">
                                <label>Nomor Induk Kependudukan (NIK)</label>
                                <input 
                                    type="text" 
                                    value={form.nik} 
                                    onChange={(e) => setForm({ ...form, nik: e.target.value })} 
                                    required 
                                />
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group-pub">
                                <label>Tempat Lahir</label>
                                <input 
                                    type="text" 
                                    value={form.birth_place} 
                                    onChange={(e) => setForm({ ...form, birth_place: e.target.value })} 
                                    required 
                                />
                            </div>
                            <div className="form-group-pub">
                                <label>Tanggal Lahir</label>
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
                                <label>Jenis Kelamin</label>
                                <select 
                                    value={form.gender} 
                                    onChange={(e) => setForm({ ...form, gender: e.target.value })}
                                >
                                    <option value="Laki-laki">Laki-laki</option>
                                    <option value="Perempuan">Perempuan</option>
                                </select>
                            </div>
                            <div className="form-group-pub">
                                <label>Nomor HP / WhatsApp</label>
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
                            <label>Alamat Email</label>
                            <input 
                                type="email" 
                                value={form.email} 
                                onChange={(e) => setForm({ ...form, email: e.target.value })} 
                                required 
                            />
                        </div>
                    </div>

                    {/* Section 2: Installation */}
                    <div className="form-section">
                        <h3>2. Alamat & Pilihan Layanan</h3>
                        <div className="form-group-pub">
                            <label>Alamat Instalasi (Tempat Internet Dipasang)</label>
                            <textarea 
                                rows="3" 
                                value={form.installation_address} 
                                onChange={(e) => setForm({ ...form, installation_address: e.target.value })} 
                                required
                            />
                        </div>

                        <div className="form-group-pub">
                            <div className="flex-between">
                                <label>Alamat Penagihan (Billing Address)</label>
                                <button 
                                    type="button" 
                                    className="link-btn-pub"
                                    onClick={() => setForm(f => ({ ...f, billing_address: f.installation_address }))}
                                >
                                    Sama dengan Alamat Instalasi
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
                            <label>Paket Internet yang Dipilih</label>
                            {loadingPackages ? (
                                <p className="small text-muted">Memuat paket internet...</p>
                            ) : (
                                <select 
                                    value={form.package_id} 
                                    onChange={(e) => setForm({ ...form, package_id: e.target.value })}
                                >
                                    {packages.map(p => (
                                        <option key={p.id} value={p.id}>
                                            {p.name} - {p.speed_mbps} Mbps (Rp {p.price.toLocaleString()}/bln)
                                        </option>
                                    ))}
                                </select>
                            )}
                        </div>

                        <div className="form-group-pub">
                            <label>Lokasi Pemasangan (GPS Coordinates)</label>
                            <div className="location-picker">
                                <input 
                                    type="number" 
                                    step="any" 
                                    placeholder="Latitude" 
                                    value={form.latitude} 
                                    onChange={(e) => setForm({ ...form, latitude: e.target.value })}
                                    required
                                />
                                <input 
                                    type="number" 
                                    step="any" 
                                    placeholder="Longitude" 
                                    value={form.longitude} 
                                    onChange={(e) => setForm({ ...form, longitude: e.target.value })}
                                    required
                                />
                                <button type="button" onClick={handleGetLocation} className="btn-pub btn-secondary-pub">
                                    <MapPin size={16} /> Deteksi GPS
                                </button>
                            </div>
                        </div>

                        <div className="form-group-pub">
                            <label>Catatan Tambahan (Opsional)</label>
                            <textarea 
                                rows="2" 
                                value={form.notes} 
                                onChange={(e) => setForm({ ...form, notes: e.target.value })} 
                            />
                        </div>
                    </div>

                    {/* Section 3: Document Upload */}
                    <div className="form-section">
                        <h3>3. Unggah Dokumen Verifikasi</h3>
                        <p className="section-desc">Format dokumen yang didukung: JPG, JPEG, PNG, PDF (maks 5MB per berkas)</p>
                        
                        <div className="upload-grid-pub">
                            <div className="upload-item">
                                <label className="upload-label">
                                    <Upload size={20} />
                                    <span>Foto KTP *</span>
                                    <input 
                                        type="file" 
                                        accept="image/*,application/pdf"
                                        onChange={(e) => handleFileChange(e, 'ktp')} 
                                        required 
                                    />
                                </label>
                                {files.ktp && <div className="file-name"><FileCheck size={14} /> {files.ktp.name}</div>}
                            </div>

                            <div className="upload-item">
                                <label className="upload-label">
                                    <Upload size={20} />
                                    <span>Selfie dengan KTP (Opsional)</span>
                                    <input 
                                        type="file" 
                                        accept="image/*"
                                        onChange={(e) => handleFileChange(e, 'selfie')} 
                                    />
                                </label>
                                {files.selfie && <div className="file-name"><FileCheck size={14} /> {files.selfie.name}</div>}
                            </div>

                            <div className="upload-item">
                                <label className="upload-label">
                                    <Upload size={20} />
                                    <span>Foto Depan Rumah (Opsional)</span>
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
                                    <span>Foto Titik Instalasi *</span>
                                    <input 
                                        type="file" 
                                        accept="image/*"
                                        onChange={(e) => handleFileChange(e, 'installation')} 
                                        required 
                                    />
                                </label>
                                {files.installation && <div className="file-name"><FileCheck size={14} /> {files.installation.name}</div>}
                            </div>

                            <div className="upload-item">
                                <label className="upload-label">
                                    <Upload size={20} />
                                    <span>Dokumen Pendukung Lain (Opsional)</span>
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
                        {submitting ? 'Mengirim Data Pendaftaran...' : 'Kirim Pendaftaran Pelanggan'}
                    </button>
                </form>
            </div>
        </div>
    );
};
