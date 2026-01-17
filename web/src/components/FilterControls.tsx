import React from 'react';
import { Search, X } from 'lucide-react';

interface FilterControlsProps {
    filterText: string;
    setFilterText: (text: string) => void;
    dateFilter: string | null;
    setDateFilter: (val: string | null) => void;
    agencyFilter: string | null;
    setAgencyFilter: (val: string | null) => void;
    count: number;
    total: number;
}

const FilterControls: React.FC<FilterControlsProps> = ({
    filterText,
    setFilterText,
    dateFilter,
    setDateFilter,
    agencyFilter,
    setAgencyFilter,
    count,
    total
}) => {
    return (
        <div className="controls-section">
            <div className="search-group">
                <div className="search-box">
                    <Search size={18} className="search-icon" />
                    <input 
                        type="text" 
                        placeholder="Search opportunities..." 
                        value={filterText}
                        onChange={(e) => setFilterText(e.target.value)}
                        className="search-input"
                    />
                </div>
                {/* Active Filters */}
                <div className="active-filters">
                    {dateFilter && (
                        <span className="filter-chip" onClick={() => setDateFilter(null)}>
                            {dateFilter} <X size={14} />
                        </span>
                    )}
                    {agencyFilter && (
                        <span className="filter-chip" onClick={() => setAgencyFilter(null)}>
                            {agencyFilter} <X size={14} />
                        </span>
                    )}
                    {(dateFilter || agencyFilter) && (
                        <span className="clear-all" onClick={() => {setDateFilter(null); setAgencyFilter(null)}}>
                            Clear all
                        </span>
                    )}
                </div>
            </div>
            <div className="stats-box">
                Showing {count} of {total} results
            </div>
        </div>
    );
};

export default FilterControls;
