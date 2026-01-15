import React, { useEffect, useState } from 'react';
import type { Match } from '../types';
import { useAuth } from '../context/AuthContext';
import { ChevronDown, ChevronRight, ExternalLink, FileText, Sparkles } from 'lucide-react';

const PersonalInbox: React.FC = () => {
    const { user } = useAuth();
    const [matches, setMatches] = useState<Match[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [expandedRow, setExpandedRow] = useState<string | null>(null);

    useEffect(() => {
        if (!user) return;

        const fetchMatches = async () => {
            try {
                const response = await fetch('/api/matches');
                if (response.ok) {
                    const data = await response.json();
                    setMatches(data || []);
                }
            } finally {
                setLoading(false);
            }
        };

        fetchMatches();
    }, [user]);

    if (!user) return <div>Please login to view your inbox.</div>;
    if (loading) return <div className="loading">Loading AI matches...</div>;

    if (matches.length === 0) {
        return (
            <div style={{textAlign: 'center', padding: '3rem', color: '#7f8c8d'}}>
                <Sparkles size={48} style={{marginBottom: '1rem', opacity: 0.5}} />
                <h3>No matches found yet.</h3>
                <p>Ensure you have set your narrative and run the matching engine.</p>
            </div>
        );
    }

    return (
        <div className="solicitation-list">
            <h2>Personalized Opportunities</h2>
            <div className="table-responsive">
                <table>
                    <thead>
                        <tr>
                            <th style={{width: 40}}></th>
                            <th style={{width: 80}}>Score</th>
                            <th>Agency</th>
                            <th>Title</th>
                            <th>Reasoning</th>
                            <th>Action</th>
                        </tr>
                    </thead>
                    <tbody>
                        {matches.map((match) => (
                            <React.Fragment key={match.match_id}>
                                <tr onClick={() => setExpandedRow(expandedRow === match.solicitation.source_id ? null : match.solicitation.source_id)}>
                                    <td className="chevron-cell">
                                        {expandedRow === match.solicitation.source_id ? <ChevronDown size={20} /> : <ChevronRight size={20} />}
                                    </td>
                                    <td>
                                        <span className="badge" style={{
                                            backgroundColor: match.score >= 80 ? '#27ae60' : match.score >= 50 ? '#f39c12' : '#e74c3c'
                                        }}>
                                            {match.score}
                                        </span>
                                    </td>
                                    <td>{match.solicitation.agency}</td>
                                    <td>{match.solicitation.title}</td>
                                    <td>
                                        <div style={{
                                            maxWidth: '400px', 
                                            whiteSpace: 'nowrap', 
                                            overflow: 'hidden', 
                                            textOverflow: 'ellipsis',
                                            color: '#555'
                                        }}>
                                            {match.explanation}
                                        </div>
                                    </td>
                                    <td>
                                        <a 
                                            href={match.solicitation.url} 
                                            target="_blank" 
                                            rel="noreferrer" 
                                            className="btn-link"
                                            onClick={(e) => e.stopPropagation()}
                                        >
                                            <ExternalLink size={16} />
                                        </a>
                                    </td>
                                </tr>
                                {expandedRow === match.solicitation.source_id && (
                                    <tr className="expanded-row">
                                        <td colSpan={6}>
                                            <div className="details-panel">
                                                <div style={{marginBottom: '1.5rem', padding: '1rem', backgroundColor: '#fdfdfd', borderLeft: '4px solid #3498db'}}>
                                                    <strong>AI Analysis:</strong>
                                                    <p style={{marginTop: '0.5rem'}}>{match.explanation}</p>
                                                </div>
                                                
                                                {match.solicitation.description && (
                                                    <div style={{marginBottom: '1rem'}}>
                                                        <strong>Description:</strong>
                                                        <p>{match.solicitation.description}</p>
                                                    </div>
                                                )}

                                                {match.solicitation.documents?.length > 0 && (
                                                    <>
                                                        <strong>Documents:</strong>
                                                        <ul className="doc-grid">
                                                            {match.solicitation.documents.map((doc, idx) => (
                                                                <li key={idx}>
                                                                    <a href={doc.url} target="_blank" rel="noreferrer" className="doc-link">
                                                                        <FileText size={16} /> {doc.title || "File"}
                                                                    </a>
                                                                </li>
                                                            ))}
                                                        </ul>
                                                    </>
                                                )}
                                            </div>
                                        </td>
                                    </tr>
                                )}
                            </React.Fragment>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default PersonalInbox;
