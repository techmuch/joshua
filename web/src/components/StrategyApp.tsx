import React, { useState } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { Network, Plus, ChevronRight, ChevronDown, Trash2 } from 'lucide-react';

// --- Types ---
interface Goal {
    id: string;
    title: string;
    level: 'Strategic' | 'Operational' | 'Tactical';
    children: Goal[];
    expanded?: boolean;
}

// --- Components ---

const GoalItem: React.FC<{ 
    goal: Goal; 
    onToggle: (id: string) => void; 
    onDelete: (id: string) => void; 
    level: number 
}> = ({ goal, onToggle, onDelete, level }) => {
    return (
        <div style={{ marginLeft: `${level * 20}px`, marginBottom: '0.5rem' }}>
            <div style={{ 
                display: 'flex', alignItems: 'center', gap: '0.5rem', 
                padding: '0.5rem', background: 'var(--bg-card)', 
                border: '1px solid var(--border-color)', borderRadius: '4px' 
            }}>
                <button 
                    onClick={() => onToggle(goal.id)}
                    style={{ background: 'none', border: 'none', cursor: 'pointer', padding: 0, display: 'flex' }}
                >
                    {goal.children.length > 0 ? (
                        goal.expanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />
                    ) : <span style={{ width: 16 }} />}
                </button>
                
                <span className={`badge`} style={{ 
                    background: goal.level === 'Strategic' ? 'var(--primary-color)' : goal.level === 'Operational' ? '#f1c40f' : '#2ecc71',
                    fontSize: '0.7rem', color: 'var(--bg-card)'
                }}>
                    {goal.level[0]}
                </span>
                
                <span style={{ flex: 1, fontWeight: 500 }}>{goal.title}</span>
                
                <button onClick={() => onDelete(goal.id)} style={{ background: 'none', border: 'none', color: 'var(--error-color)', cursor: 'pointer' }}>
                    <Trash2 size={14} />
                </button>
            </div>
            
            {goal.expanded && goal.children.map(child => (
                <GoalItem key={child.id} goal={child} onToggle={onToggle} onDelete={onDelete} level={level + 1} />
            ))}
        </div>
    );
};

const GoalDefinition: React.FC = () => {
    const [goals, setGoals] = useState<Goal[]>([
        {
            id: '1', title: 'Dominate Government AI Market', level: 'Strategic', expanded: true, children: [
                { id: '2', title: 'Win 3 Major DoD Contracts', level: 'Operational', children: [], expanded: false },
                { id: '3', title: 'Launch Autonomous Drone Platform', level: 'Operational', children: [], expanded: false }
            ]
        }
    ]);

    const toggleGoal = (id: string) => {
        // Recursive toggle (simplified for top-level demo)
        const newGoals = [...goals];
        const goal = newGoals.find(g => g.id === id);
        if (goal) goal.expanded = !goal.expanded;
        setGoals(newGoals);
    };

    return (
        <div className="chart-card">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
                <h3 style={{ margin: 0 }}>Organizational Goals</h3>
                <button className="btn-primary"><Plus size={16} /> Add Goal</button>
            </div>
            <div>
                {goals.map(g => (
                    <GoalItem key={g.id} goal={g} onToggle={toggleGoal} onDelete={() => {}} level={0} />
                ))}
            </div>
        </div>
    );
};

const BayesianNetwork: React.FC = () => {
    return (
        <div className="chart-card" style={{ height: '60vh', display: 'flex', alignItems: 'center', justifyContent: 'center', border: '2px dashed var(--border-color)' }}>
            <div style={{ textAlign: 'center', color: 'var(--text-secondary)' }}>
                <Network size={48} style={{ marginBottom: '1rem', opacity: 0.5 }} />
                <h3>Bayesian Network Designer</h3>
                <p>Define causal relationships and probabilities here.</p>
                <button className="btn-primary" style={{ marginTop: '1rem' }}>Open Editor</button>
            </div>
        </div>
    );
};

const StrategyApp: React.FC = () => {
    return (
        <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
            <Routes>
                <Route path="goals" element={<GoalDefinition />} />
                <Route path="network" element={<BayesianNetwork />} />
                <Route path="/" element={<Navigate to="goals" replace />} />
            </Routes>
        </div>
    );
};

export default StrategyApp;