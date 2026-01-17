import { useMemo } from 'react';

// Helper to calc days remaining
export const getDaysRemaining = (dateStr: string) => {
    if (dateStr === "0001-01-01T00:00:00Z") return -999;
    const diff = new Date(dateStr).getTime() - new Date().getTime();
    return Math.ceil(diff / (1000 * 60 * 60 * 24));
};

export const getDateBucket = (days: number) => {
    if (days < 0) return "Expired";
    if (days <= 7) return "0-7 Days";
    if (days <= 14) return "8-14 Days";
    if (days <= 30) return "15-30 Days";
    return "30+ Days";
};

export interface AnalyticsResult {
    timeData: { name: string; count: number }[];
    agencyData: { name: string; count: number }[];
}

export function useAnalytics<T>(
    items: T[],
    getDate: (item: T) => string,
    getAgency: (item: T) => string,
    currentDateFilter: string | null,
    currentAgencyFilter: string | null
): AnalyticsResult {

    // 1. Time Data (Filtered by Agency if active)
    const timeData = useMemo(() => {
        const buckets = {
            "Expired": 0, "0-7 Days": 0, "8-14 Days": 0, "15-30 Days": 0, "30+ Days": 0
        };
        
        const source = currentAgencyFilter 
            ? items.filter(i => getAgency(i) === currentAgencyFilter) 
            : items;

        source.forEach(item => {
            const days = getDaysRemaining(getDate(item));
            if (days === -999) return;
            const bucket = getDateBucket(days);
            if (buckets[bucket as keyof typeof buckets] !== undefined) {
                buckets[bucket as keyof typeof buckets]++;
            }
        });
        return Object.entries(buckets).map(([name, count]) => ({ name, count }));
    }, [items, currentAgencyFilter, getDate, getAgency]);

    // 2. Agency Data (Filtered by Date if active)
    const agencyData = useMemo(() => {
        const counts: Record<string, number> = {};
        
        const source = currentDateFilter
            ? items.filter(i => getDateBucket(getDaysRemaining(getDate(i))) === currentDateFilter)
            : items;

        source.forEach(item => {
            const agency = getAgency(item);
            counts[agency] = (counts[agency] || 0) + 1;
        });

        const sorted = Object.entries(counts).sort((a, b) => b[1] - a[1]);
        // Top 5
        return sorted.slice(0, 5).map(([name, count]) => ({ name, count }));
    }, [items, currentDateFilter, getDate, getAgency]);

    return { timeData, agencyData };
}
