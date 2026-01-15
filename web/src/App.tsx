import './App.css'
import SolicitationList from './components/SolicitationList'
import NarrativeEditor from './components/NarrativeEditor'
import PersonalInbox from './components/PersonalInbox'
import { AuthProvider, useAuth } from './context/AuthContext'
import { LoginButton } from './components/LoginButton'
import { useState } from 'react'
import { LayoutGrid, UserCircle, Inbox } from 'lucide-react'

function AppContent() {
  const [activeTab, setActiveTab] = useState<'library' | 'narrative' | 'inbox'>('library');
  const { user } = useAuth();

  return (
    <div className="app-container">
      <header className="app-header">
        <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem'}}>
          <h1>BD_Bot Intelligence Portal</h1>
          <LoginButton />
        </div>
        
        <nav className="nav-tabs" style={{display: 'flex', gap: '1rem', borderBottom: '1px solid #eee'}}>
          <button 
            className={`nav-tab ${activeTab === 'library' ? 'active' : ''}`}
            onClick={() => setActiveTab('library')}
            style={{
              padding: '0.8rem 1.2rem',
              background: 'transparent',
              border: 'none',
              borderBottom: activeTab === 'library' ? '2px solid #3498db' : '2px solid transparent',
              color: activeTab === 'library' ? '#3498db' : '#7f8c8d',
              cursor: 'pointer',
              fontWeight: 500,
              display: 'flex',
              alignItems: 'center',
              gap: '0.5rem'
            }}
          >
            <LayoutGrid size={18} /> Opportunity Library
          </button>
          
          {user && (
            <>
              <button 
                className={`nav-tab ${activeTab === 'inbox' ? 'active' : ''}`}
                onClick={() => setActiveTab('inbox')}
                style={{
                  padding: '0.8rem 1.2rem',
                  background: 'transparent',
                  border: 'none',
                  borderBottom: activeTab === 'inbox' ? '2px solid #3498db' : '2px solid transparent',
                  color: activeTab === 'inbox' ? '#3498db' : '#7f8c8d',
                  cursor: 'pointer',
                  fontWeight: 500,
                  display: 'flex',
                  alignItems: 'center',
                  gap: '0.5rem'
                }}
              >
                <Inbox size={18} /> Personal Inbox
              </button>
              <button 
                className={`nav-tab ${activeTab === 'narrative' ? 'active' : ''}`}
                onClick={() => setActiveTab('narrative')}
                style={{
                  padding: '0.8rem 1.2rem',
                  background: 'transparent',
                  border: 'none',
                  borderBottom: activeTab === 'narrative' ? '2px solid #3498db' : '2px solid transparent',
                  color: activeTab === 'narrative' ? '#3498db' : '#7f8c8d',
                  cursor: 'pointer',
                  fontWeight: 500,
                  display: 'flex',
                  alignItems: 'center',
                  gap: '0.5rem'
                }}
              >
                <UserCircle size={18} /> My Narrative
              </button>
            </>
          )}
        </nav>
      </header>
      
      <main>
        {activeTab === 'library' && <SolicitationList />}
        {activeTab === 'inbox' && <PersonalInbox />}
        {activeTab === 'narrative' && <NarrativeEditor />}
      </main>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  )
}

export default App
