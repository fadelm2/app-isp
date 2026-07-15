import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { authService, publicService } from '../services/api';
import { Shield, Lock, User, AlertCircle } from 'lucide-react';
import './Login.css';

export const Login = () => {
    const [id, setId] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const [ispName, setIspName] = useState(localStorage.getItem('isp_name') || 'GREENET');

    useEffect(() => {
        publicService.getIspInfo()
            .then(res => {
                if (res.data && res.data.isp_name) {
                    setIspName(res.data.isp_name);
                    localStorage.setItem('isp_name', res.data.isp_name);
                }
            })
            .catch(err => console.error("Error fetching ISP info:", err));
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            const response = await authService.login({ id, password });
            if (response.data && response.data.data) {
                const token = response.data.data.token;
                localStorage.setItem('token', token);
                navigate('/');
            } else {
                setError('Invalid login response. Please contact support.');
            }
        } catch (err) {
            console.error("Login error", err);
            setError(err.response?.data?.errors || err.response?.data?.error || 'Invalid credentials or connection error.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="login-wrapper">
            <div className="login-gradient-bg"></div>
            <div className="login-card">
                <div className="login-header">
                    <div className="login-logo">
                        <Shield size={32} className="logo-icon-purple" />
                    </div>
                    <h2>{ispName.toUpperCase()}</h2>
                    <p className="subtitle">ISP Management Portal</p>
                </div>

                {error && (
                    <div className="login-error">
                        <AlertCircle size={16} />
                        <span>{error}</span>
                    </div>
                )}

                <form onSubmit={handleSubmit} className="login-form">
                    <div className="input-group">
                        <label>Username / ID</label>
                        <div className="input-field">
                            <User size={18} className="input-icon" />
                            <input 
                                type="text" 
                                value={id} 
                                onChange={(e) => setId(e.target.value)} 
                                placeholder="Enter your ID (e.g. admin1)"
                                required 
                            />
                        </div>
                    </div>

                    <div className="input-group">
                        <label>Password</label>
                        <div className="input-field">
                            <Lock size={18} className="input-icon" />
                            <input 
                                type="password" 
                                value={password} 
                                onChange={(e) => setPassword(e.target.value)} 
                                placeholder="Enter password"
                                required 
                            />
                        </div>
                    </div>

                    <button type="submit" className="login-btn" disabled={loading}>
                        {loading ? 'Authenticating...' : 'Sign In'}
                    </button>
                </form>

                <div className="login-footer">
                    <p className="text-muted small">{ispName.toUpperCase()} OSS/BSS System v1.0.0</p>
                </div>
            </div>
        </div>
    );
};
