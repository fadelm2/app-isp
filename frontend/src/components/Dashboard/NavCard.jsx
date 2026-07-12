import React from 'react';
import { useNavigate } from 'react-router-dom';
import './NavCard.css';

const NavCard = ({ icon: Icon, title, subtitle, variant, path = '#' }) => {
    const navigate = useNavigate();
    return (
        <div className={`nav-card ${variant}`} onClick={() => navigate(path)}>
            <div className="nav-card-icon">
                <Icon size={32} />
            </div>
            <div className="nav-card-content">
                <h3>{title}</h3>
                {subtitle && <p className="subtitle">{subtitle}</p>}
            </div>
        </div>
    );
};

export default NavCard;
