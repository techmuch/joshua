import React from 'react';
import { useAuth } from '../context/AuthContext';
import { LogIn, LogOut, User } from 'lucide-react';

export const LoginButton: React.FC = () => {
    const { user, login, logout, isLoading } = useAuth();

    const handleLogin = () => {
        // Mock Login
        login("demo@example.com", "Demo User");
    };

    if (isLoading) return <div className="text-muted">...</div>;

    if (user) {
        return (
            <div className="user-menu">
                <span className="user-name">
                    <User size={18} />
                    {user.full_name}
                </span>
                <button 
                    onClick={logout} 
                    className="btn-outline"
                >
                    <LogOut size={16} /> Logout
                </button>
            </div>
        );
    }

    return (
        <button 
            onClick={handleLogin} 
            className="btn-primary"
        >
            <LogIn size={16} /> Login (Mock)
        </button>
    );
};