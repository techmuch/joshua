import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { Send, ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';

const FeedbackApp: React.FC = () => {
    const { user } = useAuth();
    const [appName, setAppName] = useState("BD_Bot");
    const [viewName, setViewName] = useState("Landing Page");
    const [content, setContent] = useState("");
    const [status, setStatus] = useState<{msg: string, type: 'success' | 'error'} | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const apps = ["BD_Bot", "Feedback App", "Developer Tools"];
    const views = ["Landing Page", "Library", "Inbox", "Detail View", "Profile", "Login"];

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);
        setStatus(null);

        try {
            const res = await fetch('/api/feedback', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ app_name: appName, view_name: viewName, content }),
            });

            if (!res.ok) throw new Error("Failed to submit feedback");

            setStatus({ msg: "Feedback submitted successfully!", type: 'success' });
            setContent("");
        } catch (err: any) {
            setStatus({ msg: "Error submitting feedback.", type: 'error' });
        } finally {
            setIsSubmitting(false);
        }
    };

    if (!user) return <div>Please login to submit feedback.</div>;

    return (
        <div style={{maxWidth: '600px', margin: '0 auto', padding: '2rem'}}>
            <Link to="/" style={{display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--text-secondary)', textDecoration: 'none', marginBottom: '2rem'}}>
                <ArrowLeft size={16} /> Back to Hub
            </Link>

            <div className="chart-card">
                <h2 style={{marginTop: 0, borderBottom: '1px solid var(--border-color)', paddingBottom: '1rem', marginBottom: '1.5rem', color: 'var(--text-primary)'}}>Provide Feedback</h2>
                
                <form onSubmit={handleSubmit}>
                    <div style={{display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginBottom: '1rem'}}>
                        <div>
                            <label style={{display: 'block', marginBottom: '0.5rem', fontWeight: 500, color: 'var(--text-body)'}}>Application</label>
                            <select 
                                className="search-input"
                                value={appName}
                                onChange={(e) => setAppName(e.target.value)}
                                style={{border: '1px solid var(--border-input)', padding: '0.75rem', borderRadius: '4px', width: '100%', background: 'var(--bg-input)', color: 'var(--text-body)'}}
                            >
                                {apps.map(a => <option key={a} value={a}>{a}</option>)}
                            </select>
                        </div>
                        <div>
                            <label style={{display: 'block', marginBottom: '0.5rem', fontWeight: 500, color: 'var(--text-body)'}}>View / Context</label>
                            <select 
                                className="search-input"
                                value={viewName}
                                onChange={(e) => setViewName(e.target.value)}
                                style={{border: '1px solid var(--border-input)', padding: '0.75rem', borderRadius: '4px', width: '100%', background: 'var(--bg-input)', color: 'var(--text-body)'}}
                            >
                                {views.map(v => <option key={v} value={v}>{v}</option>)}
                            </select>
                        </div>
                    </div>

                    <div style={{marginBottom: '1.5rem'}}>
                        <label style={{display: 'block', marginBottom: '0.5rem', fontWeight: 500, color: 'var(--text-body)'}}>Feedback & Suggestions</label>
                        <textarea
                            value={content}
                            onChange={(e) => setContent(e.target.value)}
                            placeholder="Describe your issue or idea..."
                            style={{
                                width: '100%', height: '150px', padding: '1rem', 
                                border: '1px solid var(--border-input)', borderRadius: '4px', 
                                fontFamily: 'inherit', fontSize: '1rem', boxSizing: 'border-box',
                                background: 'var(--bg-input)', color: 'var(--text-body)'
                            }}
                            required
                        />
                    </div>

                    {status && (
                        <div style={{
                            padding: '0.75rem', borderRadius: '4px', marginBottom: '1rem',
                            backgroundColor: status.type === 'success' ? 'var(--success-color)' : 'var(--error-color)',
                            color: 'white'
                        }}>
                            {status.msg}
                        </div>
                    )}

                    <button 
                        type="submit" 
                        className="btn-primary" 
                        disabled={isSubmitting}
                        style={{width: '100%', justifyContent: 'center'}}
                    >
                        <Send size={18} /> Submit Feedback
                    </button>
                </form>
            </div>
        </div>
    );
};

export default FeedbackApp;