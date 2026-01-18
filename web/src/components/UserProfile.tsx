import React, { useState, useRef } from 'react';
import { useAuth } from '../context/AuthContext';
import { User, Lock, Save, Camera } from 'lucide-react';

const UserProfile: React.FC = () => {
    const { user, refreshUser } = useAuth();
    
    // Profile State
    const [fullName, setFullName] = useState(user?.full_name || "");
    const [email, setEmail] = useState(user?.email || "");
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const [avatarPreview, setAvatarPreview] = useState<string | null>(user?.avatar_url || null);
    const [profileStatus, setProfileStatus] = useState<{msg: string, type: 'success' | 'error'} | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    // Password State
    const [oldPassword, setOldPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [passwordStatus, setPasswordStatus] = useState<{msg: string, type: 'success' | 'error'} | null>(null);
    
    const [isSubmitting, setIsSubmitting] = useState(false);

    if (!user) return <div>Please login to view profile.</div>;

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            const file = e.target.files[0];
            setAvatarFile(file);
            setAvatarPreview(URL.createObjectURL(file));
        }
    };

    const handleProfileUpdate = async (e: React.FormEvent) => {
        e.preventDefault();
        setProfileStatus(null);
        setIsSubmitting(true);

        try {
            let avatarUrl = user.avatar_url || "";

            // 1. Upload Avatar if changed
            if (avatarFile) {
                const formData = new FormData();
                formData.append('avatar', avatarFile);
                const res = await fetch('/api/user/avatar', { method: 'POST', body: formData });
                if (!res.ok) throw new Error("Failed to upload avatar");
                const data = await res.json();
                avatarUrl = data.url;
            }

            // 2. Update Profile
            const res = await fetch('/api/user/profile', {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ 
                    email: email, 
                    full_name: fullName, 
                    avatar_url: avatarUrl 
                }),
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || "Failed to update profile");
            }

            setProfileStatus({ msg: "Profile updated successfully", type: 'success' });
            await refreshUser();
        } catch (err: any) {
            setProfileStatus({ msg: err.message, type: 'error' });
        } finally {
            setIsSubmitting(false);
        }
    };

    const handlePasswordChange = async (e: React.FormEvent) => {
        e.preventDefault();
        setPasswordStatus(null);

        if (newPassword !== confirmPassword) {
            setPasswordStatus({ msg: "New passwords do not match", type: 'error' });
            return;
        }

        setIsSubmitting(true);
        try {
            const res = await fetch('/api/auth/password', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || "Failed to update password");
            }

            setPasswordStatus({ msg: "Password updated successfully", type: 'success' });
            setOldPassword("");
            setNewPassword("");
            setConfirmPassword("");
        } catch (err: any) {
            setPasswordStatus({ msg: err.message, type: 'error' });
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="solicitation-list" style={{maxWidth: '900px', margin: '0 auto'}}>
            <div style={{borderBottom: '1px solid #eee', paddingBottom: '1rem', marginBottom: '2rem'}}>
                <h2>User Profile</h2>
            </div>

            <div style={{display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2rem'}}>
                
                {/* --- Profile Info & Avatar --- */}
                <div className="chart-card">
                    <div style={{display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '1.5rem'}}>
                        <User size={20} color="#34495e" />
                        <h3 style={{margin: 0, color: '#34495e'}}>Public Profile</h3>
                    </div>

                    <form onSubmit={handleProfileUpdate}>
                        {/* Avatar Upload */}
                        <div style={{textAlign: 'center', marginBottom: '1.5rem', position: 'relative'}}>
                            <div 
                                style={{
                                    width: '100px', height: '100px', borderRadius: '50%', 
                                    background: '#e1f5fe', color: '#3498db', 
                                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                                    margin: '0 auto', overflow: 'hidden', border: '3px solid #fff',
                                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)', cursor: 'pointer',
                                    position: 'relative'
                                }}
                                onClick={() => fileInputRef.current?.click()}
                            >
                                {avatarPreview ? (
                                    <img src={avatarPreview} alt="Avatar" style={{width: '100%', height: '100%', objectFit: 'cover'}} />
                                ) : (
                                    <User size={50} />
                                )}
                                
                                {/* Overlay */}
                                <div style={{
                                    position: 'absolute', bottom: 0, left: 0, right: 0,
                                    background: 'rgba(0,0,0,0.5)', padding: '4px',
                                    display: 'flex', justifyContent: 'center'
                                }}>
                                    <Camera size={16} color="white" />
                                </div>
                            </div>
                            <input 
                                type="file" 
                                ref={fileInputRef} 
                                onChange={handleFileChange} 
                                style={{display: 'none'}} 
                                accept="image/*"
                            />
                            <p style={{fontSize: '0.8rem', color: '#7f8c8d', marginTop: '0.5rem'}}>Click to upload</p>
                        </div>

                        {/* Fields */}
                        <div style={{marginBottom: '1rem'}}>
                            <label style={{display: 'block', marginBottom: '0.5rem'}}>Full Name</label>
                            <input 
                                type="text" 
                                className="search-input"
                                style={{border: '1px solid #ddd', padding: '0.75rem', borderRadius: '4px', width: '100%', boxSizing: 'border-box'}}
                                value={fullName}
                                onChange={(e) => setFullName(e.target.value)}
                                required
                            />
                        </div>
                        <div style={{marginBottom: '1.5rem'}}>
                            <label style={{display: 'block', marginBottom: '0.5rem'}}>Email Address</label>
                            <input 
                                type="email" 
                                className="search-input"
                                style={{border: '1px solid #ddd', padding: '0.75rem', borderRadius: '4px', width: '100%', boxSizing: 'border-box'}}
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                            />
                        </div>

                        {profileStatus && (
                            <div style={{
                                padding: '0.75rem', borderRadius: '4px', marginBottom: '1rem',
                                backgroundColor: profileStatus.type === 'success' ? '#d4edda' : '#f8d7da',
                                color: profileStatus.type === 'success' ? '#155724' : '#721c24'
                            }}>
                                {profileStatus.msg}
                            </div>
                        )}

                        <button 
                            type="submit" 
                            className="btn-primary" 
                            disabled={isSubmitting}
                            style={{width: '100%', justifyContent: 'center'}}
                        >
                            <Save size={18} /> Save Profile
                        </button>
                    </form>
                </div>

                {/* --- Change Password --- */}
                <div className="chart-card">
                    <div style={{display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '1.5rem'}}>
                        <Lock size={20} color="#34495e" />
                        <h3 style={{margin: 0, color: '#34495e'}}>Security</h3>
                    </div>

                    <form onSubmit={handlePasswordChange}>
                        <div style={{marginBottom: '1rem'}}>
                            <label style={{display: 'block', marginBottom: '0.5rem'}}>Current Password</label>
                            <input 
                                type="password" 
                                className="search-input"
                                style={{border: '1px solid #ddd', padding: '0.75rem', borderRadius: '4px', width: '100%', boxSizing: 'border-box'}}
                                value={oldPassword}
                                onChange={(e) => setOldPassword(e.target.value)}
                                required
                            />
                        </div>
                        <div style={{marginBottom: '1rem'}}>
                            <label style={{display: 'block', marginBottom: '0.5rem'}}>New Password</label>
                            <input 
                                type="password" 
                                className="search-input"
                                style={{border: '1px solid #ddd', padding: '0.75rem', borderRadius: '4px', width: '100%', boxSizing: 'border-box'}}
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                                required
                            />
                        </div>
                        <div style={{marginBottom: '1.5rem'}}>
                            <label style={{display: 'block', marginBottom: '0.5rem'}}>Confirm New Password</label>
                            <input 
                                type="password" 
                                className="search-input"
                                style={{border: '1px solid #ddd', padding: '0.75rem', borderRadius: '4px', width: '100%', boxSizing: 'border-box'}}
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                required
                            />
                        </div>

                        {passwordStatus && (
                            <div style={{
                                padding: '0.75rem', borderRadius: '4px', marginBottom: '1rem',
                                backgroundColor: passwordStatus.type === 'success' ? '#d4edda' : '#f8d7da',
                                color: passwordStatus.type === 'success' ? '#155724' : '#721c24'
                            }}>
                                {passwordStatus.msg}
                            </div>
                        )}

                        <button 
                            type="submit" 
                            className="btn-primary" 
                            disabled={isSubmitting}
                            style={{width: '100%', justifyContent: 'center'}}
                        >
                            <Lock size={18} /> Update Password
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default UserProfile;