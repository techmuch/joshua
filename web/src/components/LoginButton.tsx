import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { LogIn, LogOut, User, X } from 'lucide-react';

export const LoginButton: React.FC = () => {
    const { user, login, logout, isLoading } = useAuth();
    const [isOpen, setIsOpen] = useState(false);
    const [email, setEmail] = useState("admin@example.com"); // Default for dev
    const [password, setPassword] = useState("secret123"); // Default for dev
    const [error, setError] = useState("");

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        try {
            await login(email, password);
            setIsOpen(false);
        } catch (err: any) {
            setError(err.message || "Invalid email or password");
        }
    };

    if (isLoading) return <div className="text-muted">...</div>;

    if (user) {
        return (
            <div className="user-menu">
                <span className="user-name">
                    <User size={18} />
                    {user.full_name || user.email}
                </span>
                <button onClick={logout} className="btn-outline">
                    <LogOut size={16} /> Logout
                </button>
            </div>
        );
    }

    return (
        <>
            <button onClick={() => setIsOpen(true)} className="btn-primary">
                <LogIn size={16} /> Login
            </button>

            {isOpen && (
                <div style={modalOverlayStyle}>
                    <div style={modalStyle}>
                        <div style={{display: 'flex', justifyContent: 'space-between', marginBottom: '1rem'}}>
                            <h3>Login</h3>
                            <button onClick={() => setIsOpen(false)} style={{background: 'none', border: 'none', cursor: 'pointer'}}>
                                <X size={20} />
                            </button>
                        </div>
                        
                        <form onSubmit={handleLogin} style={{display: 'flex', flexDirection: 'column', gap: '1rem'}}>
                            <div>
                                <label style={{display: 'block', marginBottom: '0.5rem', fontSize: '0.9rem', color: '#2c3e50', fontWeight: '600'}}>Email Address</label>
                                <input 
                                    type="email" 
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    style={{
                                        width: '100%',
                                        padding: '0.75rem', 
                                        borderRadius: '6px', 
                                        border: '2px solid #ddd',
                                        background: '#f9f9f9', 
                                        color: '#000',
                                        fontSize: '1rem',
                                        boxSizing: 'border-box'
                                    }}
                                    required
                                />
                            </div>
                            <div>
                                <label style={{display: 'block', marginBottom: '0.5rem', fontSize: '0.9rem', color: '#2c3e50', fontWeight: '600'}}>Password</label>
                                <input 
                                    type="password" 
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    style={{
                                        width: '100%',
                                        padding: '0.75rem', 
                                        borderRadius: '6px', 
                                        border: '2px solid #ddd',
                                        background: '#f9f9f9', 
                                        color: '#000',
                                        fontSize: '1rem',
                                        boxSizing: 'border-box'
                                    }}
                                    required
                                />
                            </div>
                            
                            {error && <div style={{color: '#e74c3c', fontSize: '0.9rem', fontWeight: '500', padding: '0.5rem', background: '#fdeded', borderRadius: '4px'}}>{error}</div>}
                            
                            <div style={{display: 'flex', gap: '1rem', justifyContent: 'flex-end', marginTop: '1rem'}}>
                                <button 
                                    type="button" 
                                    onClick={() => setIsOpen(false)}
                                    className="btn-outline"
                                    style={{justifyContent: 'center'}}
                                >
                                    Cancel
                                </button>
                                <button 
                                    type="submit" 
                                    className="btn-primary" 
                                    style={{justifyContent: 'center'}}
                                >
                                    Sign In
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </>
    );
};

const modalOverlayStyle: React.CSSProperties = {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 2000,
};

const modalStyle: React.CSSProperties = {
    backgroundColor: 'white',
    padding: '2rem',
    borderRadius: '8px',
    width: '100%',
    maxWidth: '400px',
    boxShadow: '0 10px 25px rgba(0,0,0,0.1)',
    color: '#333', // Ensure modal text is dark
};
