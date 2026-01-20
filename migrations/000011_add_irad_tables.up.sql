-- Strategic Capability Objectives (Leadership Defined)
CREATE TABLE irad_scos (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    target_spend_percent DECIMAL(5,2), -- e.g. 40.00
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- IRAD Projects (Researchers/PIs)
CREATE TABLE irad_projects (
    id SERIAL PRIMARY KEY,
    sco_id INT REFERENCES irad_scos(id),
    title TEXT NOT NULL,
    description TEXT,
    pi_id INT REFERENCES users(id), -- Principal Investigator
    status TEXT DEFAULT 'concept', -- concept, proposal, active, completed, transitioned
    total_budget DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Roadmap / Phased Budget (Multi-Year)
CREATE TABLE irad_roadmaps (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES irad_projects(id) ON DELETE CASCADE,
    fiscal_year INT NOT NULL,
    labor_cost DECIMAL(15,2) DEFAULT 0,
    odc_cost DECIMAL(15,2) DEFAULT 0,
    sub_cost DECIMAL(15,2) DEFAULT 0,
    milestones JSONB DEFAULT '[]', -- list of {title, date, status}
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Transition Logs (ROI tracking)
CREATE TABLE irad_transitions (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES irad_projects(id) ON DELETE CASCADE,
    outcome_type TEXT, -- DARPA BAA, Internal Pivot, IP Filing, Contract
    description TEXT,
    captured_funding DECIMAL(15,2) DEFAULT 0, -- ROI mapping
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
