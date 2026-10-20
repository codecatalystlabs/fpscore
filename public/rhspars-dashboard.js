sessionStorage.setItem('active_tool', 'rhspars');

const DOMAINS = [
    { key: 'ancPct', label: 'Antenatal care' },
    { key: 'maternityPct', label: 'Maternity' },
    { key: 'fpPct', label: 'Family planning' },
    { key: 'commoditiesPct', label: 'Drugs and supplies' },
    { key: 'healthInfoPct', label: 'Health information' },
    { key: 'generalFacilityPct', label: 'General facility' }
];

const COLORS = {
    proficient: '#006d41',
    competent: '#c9a227',
    poor: '#c0392b',
    track: '#e6eeea',
    blue: '#1a6bb5'
};

let allVisits = [];
const charts = {};

function visitDate(row) {
    const raw = row.supervisionDate || row.createdAt || '';
    const d = new Date(raw);
    return Number.isNaN(d.getTime()) ? null : d;
}

function mean(values) {
    if (!values.length) return 0;
    return values.reduce((sum, n) => sum + n, 0) / values.length;
}

function levelOf(pct) {
    return performanceLevelFromPct(Number(pct) || 0);
}

function barColor(pct) {
    const level = levelOf(pct);
    if (level === 'Proficient') return COLORS.proficient;
    if (level === 'Competent') return COLORS.competent;
    return COLORS.poor;
}

function heatStyle(pct) {
    const level = levelOf(pct);
    if (level === 'Proficient') return 'background:#e5f4ec;color:#006d41';
    if (level === 'Competent') return 'background:#fbf6e6;color:#8a6808';
    return 'background:#fdecea;color:#c0392b';
}

function esc(value) {
    return String(value ?? '').replace(/[&<>"']/g, ch => ({
        '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
    }[ch]));
}

function filteredVisits() {
    const period = document.getElementById('filterPeriod').value;
    const region = document.getElementById('filterRegion').value;
    const district = document.getElementById('filterDistrict').value;
    const cutoff = period === 'all' ? null : Date.now() - Number(period) * 86400000;
    return allVisits.filter(row => {
        if (region && String(row.regionName) !== region) return false;
        if (district && String(row.districtName) !== district) return false;
        if (cutoff) {
            const d = visitDate(row);
            if (!d || d.getTime() < cutoff) return false;
        }
        return true;
    });
}

function fillFilters() {
    const regionSelect = document.getElementById('filterRegion');
    const districtSelect = document.getElementById('filterDistrict');
    const regions = [...new Set(allVisits.map(r => r.regionName).filter(Boolean))].sort();
    const selectedRegion = regionSelect.value;
    regionSelect.innerHTML = '<option value="">All regions</option>' +
        regions.map(name => `<option>${esc(name)}</option>`).join('');
    regionSelect.value = selectedRegion;

    const region = regionSelect.value;
    const districts = [...new Set(allVisits
        .filter(r => !region || r.regionName === region)
        .map(r => r.districtName)
        .filter(Boolean))].sort();
    const selectedDistrict = districtSelect.value;
    districtSelect.innerHTML = '<option value="">All districts</option>' +
        districts.map(name => `<option>${esc(name)}</option>`).join('');
    if (districts.includes(selectedDistrict)) districtSelect.value = selectedDistrict;
}

function latestByFacility(rows) {
    const map = new Map();
    rows.forEach(row => {
        const key = row.facilityId || row.facilityName;
        const current = map.get(key);
        const when = visitDate(row);
        if (!current || (when && visitDate(current) && when > visitDate(current))) map.set(key, row);
    });
    return [...map.values()];
}

function destroyCharts() {
    Object.values(charts).forEach(chart => chart.destroy());
    for (const key of Object.keys(charts)) delete charts[key];
}

function render(rows) {
    destroyCharts();
    const root = document.getElementById('dashboard');
    if (!rows.length) {
        root.innerHTML = '<div class="moh-card"><h5>No visits in this view</h5><p class="mb-0 text-muted">Change the period or area, or record a visit.</p></div>';
        return;
    }

    const grand = rows.map(r => Number(r.grandPercentage) || 0);
    const avg = mean(grand);
    const facilities = latestByFacility(rows);
    const counts = { Proficient: 0, Competent: 0, 'Not Acceptable': 0 };
    rows.forEach(r => {
        const level = r.performanceLevel || levelOf(r.grandPercentage);
        counts[level] = (counts[level] || 0) + 1;
    });
    const domainAvgs = DOMAINS.map(d => ({
        ...d,
        value: mean(rows.map(r => Number(r[d.key]) || 0))
    }));
    const weakest = [...domainAvgs].sort((a, b) => a.value - b.value)[0];
    const strongest = [...domainAvgs].sort((a, b) => b.value - a.value)[0];
    const below = facilities.filter(r => Number(r.grandPercentage) < 70).length;
    const selfCount = rows.filter(r => r.supervisionType === 'self').length;
    const externalCount = rows.filter(r => r.supervisionType === 'external').length;
    const ranked = facilities.slice().sort((a, b) => Number(a.grandPercentage) - Number(b.grandPercentage));

    root.innerHTML = `
        <div class="kpi-grid">
            <div class="kpi"><div class="label">Visits</div><div class="value">${rows.length}</div><div class="hint">${externalCount} external · ${selfCount} self</div></div>
            <div class="kpi"><div class="label">Mean grand score</div><div class="value">${avg.toFixed(1)}%</div><div class="hint">${esc(levelOf(avg))}</div></div>
            <div class="kpi"><div class="label">Proficient</div><div class="value">${counts.Proficient}</div><div class="hint">${rows.length ? Math.round(counts.Proficient / rows.length * 100) : 0}% of visits</div></div>
            <div class="kpi"><div class="label">Competent</div><div class="value">${counts.Competent}</div><div class="hint">70–89%</div></div>
            <div class="kpi ${below ? 'alert' : ''}"><div class="label">Facilities under 70%</div><div class="value">${below}</div><div class="hint">Latest visit in this view</div></div>
        </div>
        <div class="insight">
            <i class="bi bi-lightbulb"></i>
            <div>Strongest block is <strong>${strongest.label}</strong> at ${strongest.value.toFixed(1)}%. The weakest is <strong>${weakest.label}</strong> at ${weakest.value.toFixed(1)}%.</div>
        </div>
        <div class="chart-grid mb-3">
            <div class="moh-card"><h5>Performance mix</h5><canvas id="mixChart"></canvas></div>
            <div class="moh-card"><h5>Average score by domain</h5><canvas id="domainChart"></canvas></div>
        </div>
        <div class="moh-card mb-3"><h5>Grand score over time</h5><canvas id="trendChart"></canvas></div>
        <div class="moh-card mb-3">
            <h5>Latest score sheet</h5>
            <p class="text-muted mb-2">One row per facility, using the most recent visit in this view. Green is proficient, gold is competent, red is not acceptable.</p>
            <div class="sheet table-responsive">
                <table class="table table-sm mb-0">
                    <thead><tr>
                        <th>Facility</th><th>District</th>
                        ${DOMAINS.map(d => `<th>${esc(d.label)}</th>`).join('')}
                        <th>Grand</th>
                    </tr></thead>
                    <tbody>
                        ${ranked.map(r => `<tr>
                            <td><a href="rhspars-view.html?id=${r.id}">${esc(r.facilityName)}</a></td>
                            <td>${esc(r.districtName)}</td>
                            ${DOMAINS.map(d => {
                                const pct = Number(r[d.key]) || 0;
                                return `<td><span class="heat" style="${heatStyle(pct)}">${pct.toFixed(0)}</span></td>`;
                            }).join('')}
                            <td><span class="heat" style="${heatStyle(r.grandPercentage)}">${Number(r.grandPercentage).toFixed(1)}</span></td>
                        </tr>`).join('')}
                    </tbody>
                </table>
            </div>
        </div>
        <div class="split">
            <div class="moh-card">
                <h5>Facilities</h5>
                <div class="table-responsive">
                    <table class="table table-sm mb-0">
                        <thead><tr><th>Facility</th><th>District</th><th>Latest</th><th></th></tr></thead>
                        <tbody>
                            ${ranked.map(r => {
                                const pct = Number(r.grandPercentage);
                                const level = r.performanceLevel || levelOf(pct);
                                return `<tr>
                                    <td>${esc(r.facilityName)}</td>
                                    <td>${esc(r.districtName)}</td>
                                    <td><strong style="color:${barColor(pct)}">${pct.toFixed(1)}%</strong></td>
                                    <td><a href="rhspars-view.html?id=${r.id}">${esc(level)}</a></td>
                                </tr>`;
                            }).join('')}
                        </tbody>
                    </table>
                </div>
            </div>
            <div class="moh-card">
                <h5>Recent visits</h5>
                <div class="table-responsive">
                    <table class="table table-sm mb-0">
                        <thead><tr><th>Date</th><th>Facility</th><th>Score</th><th></th></tr></thead>
                        <tbody>
                            ${rows.slice(0, 12).map(r => {
                                const pct = Number(r.grandPercentage);
                                const when = (r.supervisionDate || r.createdAt || '').toString().slice(0, 10);
                                return `<tr>
                                    <td>${when}</td>
                                    <td>${esc(r.facilityName)}<div class="text-muted">${esc(r.regionName)}</div></td>
                                    <td><strong style="color:${barColor(pct)}">${pct.toFixed(1)}%</strong></td>
                                    <td><a href="rhspars-view.html?id=${r.id}">Open</a></td>
                                </tr>`;
                            }).join('')}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>`;

    charts.mix = new Chart(document.getElementById('mixChart'), {
        type: 'doughnut',
        data: {
            labels: ['Proficient', 'Competent', 'Not Acceptable'],
            datasets: [{
                data: [counts.Proficient, counts.Competent, counts['Not Acceptable']],
                backgroundColor: [COLORS.proficient, COLORS.competent, COLORS.poor],
                borderWidth: 0
            }]
        },
        options: { plugins: { legend: { position: 'bottom' } }, cutout: '62%' }
    });

    charts.domain = new Chart(document.getElementById('domainChart'), {
        type: 'bar',
        data: {
            labels: domainAvgs.map(d => d.label),
            datasets: [{
                data: domainAvgs.map(d => Number(d.value.toFixed(1))),
                backgroundColor: domainAvgs.map(d => barColor(d.value)),
                borderRadius: 4
            }]
        },
        options: {
            indexAxis: 'y',
            plugins: { legend: { display: false } },
            scales: { x: { min: 0, max: 100, ticks: { callback: v => v + '%' } } }
        }
    });

    const months = new Map();
    rows.forEach(r => {
        const d = visitDate(r);
        if (!d) return;
        const key = d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0');
        if (!months.has(key)) months.set(key, []);
        months.get(key).push(Number(r.grandPercentage) || 0);
    });
    const monthKeys = [...months.keys()].sort();
    charts.trend = new Chart(document.getElementById('trendChart'), {
        type: 'line',
        data: {
            labels: monthKeys,
            datasets: [{
                label: 'Mean grand score',
                data: monthKeys.map(key => Number(mean(months.get(key)).toFixed(1))),
                borderColor: COLORS.proficient,
                backgroundColor: 'rgba(0,109,65,0.12)',
                fill: true,
                tension: 0.3,
                pointRadius: 3
            }]
        },
        options: {
            plugins: { legend: { display: false } },
            scales: { y: { min: 0, max: 100, ticks: { callback: v => v + '%' } } }
        }
    });
}

async function load() {
    const res = await fetch('/api/rhspars/assessments', { headers: authHeaders() });
    if (res.status === 401) { logout(); return; }
    if (!res.ok) {
        document.getElementById('dashboard').innerHTML = '<div class="alert alert-danger">Could not load visits.</div>';
        return;
    }
    const data = await res.json();
    allVisits = Array.isArray(data) ? data : [];
    fillFilters();
    render(filteredVisits());
}

['filterPeriod', 'filterRegion', 'filterDistrict'].forEach(id => {
    document.getElementById(id).addEventListener('change', () => {
        if (id === 'filterRegion') fillFilters();
        render(filteredVisits());
    });
});

if (!localStorage.getItem('token')) window.location.href = '/';
load().catch(err => {
    console.error(err);
    document.getElementById('dashboard').textContent = 'Error loading the dashboard.';
});
