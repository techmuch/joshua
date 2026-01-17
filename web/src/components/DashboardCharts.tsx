import React from 'react';
import { 
    BarChart, 
    Bar, 
    XAxis, 
    YAxis, 
    CartesianGrid, 
    Tooltip, 
    ResponsiveContainer,
    Cell
} from 'recharts';

interface DashboardChartsProps {
    timeData: { name: string; count: number }[];
    agencyData: { name: string; count: number }[];
    dateFilter: string | null;
    agencyFilter: string | null;
    setDateFilter: (val: string | null) => void;
    setAgencyFilter: (val: string | null) => void;
}

const DashboardCharts: React.FC<DashboardChartsProps> = ({
    timeData,
    agencyData,
    dateFilter,
    agencyFilter,
    setDateFilter,
    setAgencyFilter
}) => {
    return (
        <div className="dashboard-grid">
            <div className="chart-card">
                <h3>Timeline</h3>
                <div style={{ width: '100%', height: 250 }}>
                    <ResponsiveContainer>
                        <BarChart data={timeData} onClick={(data) => data && setDateFilter(data.activeLabel ? String(data.activeLabel) : null)}>
                            <CartesianGrid strokeDasharray="3 3" vertical={false} />
                            <XAxis dataKey="name" fontSize={12} tickLine={false} axisLine={false} />
                            <YAxis allowDecimals={false} hide />
                            <Tooltip cursor={{fill: 'transparent'}} />
                            <Bar dataKey="count" radius={[4, 4, 0, 0]} cursor="pointer">
                                {timeData.map((entry, index) => (
                                    <Cell key={`cell-${index}`} fill={entry.name === dateFilter ? "#2980b9" : "#3498db"} />
                                ))}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </div>
            <div className="chart-card">
                <h3>Top Agencies</h3>
                <div style={{ width: '100%', height: 250 }}>
                    <ResponsiveContainer>
                        <BarChart layout="vertical" data={agencyData} onClick={(data) => data && setAgencyFilter(data.activeLabel ? String(data.activeLabel) : null)}>
                            <CartesianGrid strokeDasharray="3 3" horizontal={false} />
                            <XAxis type="number" hide />
                            <YAxis dataKey="name" type="category" width={120} fontSize={11} tickLine={false} axisLine={false} />
                            <Tooltip cursor={{fill: 'transparent'}} />
                            <Bar dataKey="count" radius={[0, 4, 4, 0]} barSize={20} cursor="pointer">
                                {agencyData.map((entry, index) => (
                                    <Cell key={`cell-${index}`} fill={entry.name === agencyFilter ? "#8e44ad" : "#9b59b6"} />
                                ))}
                            </Bar>
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </div>
        </div>
    );
};

export default DashboardCharts;
