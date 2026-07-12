import React from 'react';
import { MoreHorizontal, Monitor, Cloud, Database, AlertCircle } from 'lucide-react';
import './DeploymentCard.css';

const DeploymentCard = ({ title, icon: Icon, status = 'success' }) => {
    return (
        <div className="deployment-card">
            <div className="card-header">
                <div className={`status-dot ${status}`}></div>
                <button className="more-btn"><MoreHorizontal size={16} /></button>
            </div>

            <div className="card-body">
                <div className="icon-wrapper">
                    <Icon size={32} />
                </div>
                <span className="card-title">{title}</span>
            </div>
        </div>
    );
};

export default DeploymentCard;
