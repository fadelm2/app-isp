import React from 'react';
import { Bus, Car, User, Users } from 'lucide-react';
import './StatCard.css';

const StatCard = () => {
    return (
        <div className="stat-card">
            <div className="stat-main">
                <Bus size={24} className="stat-icon" />
                <div className="stat-info">
                    <span className="label">Total Vehicles</span>
                    <span className="value">154</span>
                    <span className="dots">•••••••••</span>
                </div>
            </div>

            <div className="stat-metrics">
                <div className="metric-item">
                    <Car size={16} />
                    <div className="metric-text">
                        <span className="on-duty">Online</span>
                    </div>
                </div>

                <div className="metric-item">
                    <User size={16} />
                    <div className="metric-text">
                        <span className="label">On Duty</span>
                        <span className="value">On Duty</span>
                    </div>
                </div>

                <div className="metric-item">
                    <Users size={16} />
                    <div className="metric-text">
                        <span className="label">Active</span>
                        <span className="value">Active</span>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default StatCard;
