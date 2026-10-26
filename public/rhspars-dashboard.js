sessionStorage.setItem('active_tool', 'rhspars');

const DOMAINS = [
    { key: 'maternityPct', label: 'Maternity services' },
    { key: 'fpPct', label: 'Family planning services' },
    { key: 'ancPct', label: 'Antenatal care services' },
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
    const raw = String(row.supervisionDate || row.createdAt || '').slice(0, 10);
    if (!/^\d{4}-\d{2}-\d{2}$/.test(raw)) return null;
    const d = new Date(raw + 'T00:00:00');
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

function scored(row, key) {
    const value = row[key];
    if (value === null || value === undefined || value === '') return null;
    const n = Number(value);
    return Number.isFinite(n) ? n : null;
}

function rowGrand(row) {
    const values = DOMAINS.map(d => scored(row, d.key)).filter(v => v !== null);
    if (!values.length) return null;
    return mean(values);
}

function openHref(row) {
    return row.inProgress ? 'rhspars-home.html' : 'rhspars-view.html?id=' + row.id;
}

const CODE_KEY = {
    maternity: 'maternityPct',
    fp: 'fpPct',
    anc: 'ancPct',
    commodities: 'commoditiesPct',
    health_information: 'healthInfoPct',
    general_facility: 'generalFacilityPct'
};

function scoreDraft(structure) {
    let draft = {};
    try {
        draft = JSON.parse(sessionStorage.getItem('rhspars_draft') || '{}');
    } catch (e) {
        return null;
    }
    const responses = draft.responses || {};
    if (!Object.keys(responses).length) return null;

    const domainSections = {};
    (structure || []).forEach(domain => {
        const key = CODE_KEY[domain.code];
        if (!key) return;
        const sections = [];
        (domain.thematicAreas || []).forEach(area => {
            let possible = 0;
            let achieved = 0;
            (area.questions || []).forEach(q => {
                const response = responses[String(q.id)] || responses[q.id];
                const result = scoreResponse(response, Number(q.scoreWeight) || 0);
                if (!result.counts) return;
                possible += Number(q.scoreWeight) || 0;
                achieved += result.points;
            });
            if (possible > 0) sections.push(percentageScore(achieved, possible));
        });
        if (sections.length) domainSections[key] = mean(sections);
    });
    if (!Object.keys(domainSections).length) return null;

    const row = {
        id: 'draft',
        inProgress: true,
        facilityId: draft.facilityId || 'draft',
        facilityName: draft.facilityName || 'Visit in progress',
        districtName: draft.districtName || '',
        regionName: draft.regionName || '',
        supervisionType: draft.supervisionType || '',
        supervisionDate: draft.supervisionDate || ''
    };
    DOMAINS.forEach(d => { row[d.key] = Object.prototype.hasOwnProperty.call(domainSections, d.key) ? domainSections[d.key] : null; });
    row.grandPercentage = rowGrand(row);
    row.performanceLevel = levelOf(row.grandPercentage);
    return row;
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
            if (d && d.getTime() < cutoff) return false;
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

    const grands = rows.map(rowGrand).filter(v => v !== null);
    const avg = mean(grands);
    const facilities = latestByFacility(rows);
    const counts = { Proficient: 0, Competent: 0, 'Not Acceptable': 0 };
    rows.forEach(r => {
        const pct = rowGrand(r);
        if (pct === null) return;
        const level = levelOf(pct);
        counts[level] = (counts[level] || 0) + 1;
    });
    const domainAvgs = DOMAINS.map(d => {
        const values = rows.map(r => scored(r, d.key)).filter(v => v !== null);
        return { ...d, value: values.length ? mean(values) : null };
    }).filter(d => d.value !== null);
    const weakest = [...domainAvgs].sort((a, b) => a.value - b.value)[0];
    const strongest = [...domainAvgs].sort((a, b) => b.value - a.value)[0];
    const below = facilities.filter(r => {
        const pct = rowGrand(r);
        return pct !== null && pct < 70;
    }).length;
    const selfCount = rows.filter(r => r.supervisionType === 'self').length;
    const externalCount = rows.filter(r => r.supervisionType === 'external').length;
    const ranked = facilities.slice().sort((a, b) => Number(a.grandPercentage) - Number(b.grandPercentage));

    const pivotRows = domainAvgs.map((d, i) => `<tr>
        <td>${i + 1}</td>
        <td>${esc(d.label)}</td>
        <td><strong>${d.value.toFixed(1)}%</strong></td>
    </tr>`).join('');

    root.innerHTML = `
        <div class="moh-card mb-3">
            <h5>Average percentage score</h5>
            <p class="text-muted mb-2">Average of the domains that have been scored. Domains not yet entered are left out.</p>
            <table class="table table-sm mb-0">
                <thead><tr><th>SNo</th><th>Indicator / building block</th><th>Average percentage score</th></tr></thead>
                <tbody>
                    ${pivotRows}
                    <tr><td></td><td><strong>Grand facility percentage score</strong></td><td><strong>${avg.toFixed(1)}%</strong></td></tr>
                </tbody>
            </table>
        </div>
        <div class="kpi-grid">
            <div class="kpi"><div class="label">Visits</div><div class="value">${rows.length}</div><div class="hint">${externalCount} external · ${selfCount} self</div></div>
            <div class="kpi"><div class="label">Mean grand score</div><div class="value">${avg.toFixed(1)}%</div><div class="hint">${esc(levelOf(avg))}</div></div>
            <div class="kpi"><div class="label">Proficient</div><div class="value">${counts.Proficient}</div><div class="hint">${rows.length ? Math.round(counts.Proficient / rows.length * 100) : 0}% of visits</div></div>
            <div class="kpi"><div class="label">Competent</div><div class="value">${counts.Competent}</div><div class="hint">70–89%</div></div>
            <div class="kpi ${below ? 'alert' : ''}"><div class="label">Facilities under 70%</div><div class="value">${below}</div><div class="hint">Latest visit in this view</div></div>
        </div>
        <div class="insight">
            <i class="bi bi-lightbulb"></i>
            <div>${domainAvgs.length ? `Showing ${domainAvgs.length} scored domain${domainAvgs.length === 1 ? '' : 's'}. Strongest is <strong>${esc(strongest.label)}</strong> at ${strongest.value.toFixed(1)}%. Weakest is <strong>${esc(weakest.label)}</strong> at ${weakest.value.toFixed(1)}%.` : 'No domain scores in this view yet.'}</div>
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
                        ${domainAvgs.map(d => `<th>${esc(d.label)}</th>`).join('')}
                        <th>Grand</th>
                    </tr></thead>
                    <tbody>
                        ${ranked.map(r => {
                            const grandPct = rowGrand(r);
                            return `<tr>
                            <td><a href="${openHref(r)}">${esc(r.facilityName)}</a>${r.inProgress ? ' <span class="text-muted">in progress</span>' : ''}</td>
                            <td>${esc(r.districtName)}</td>
                            ${domainAvgs.map(d => {
                                const pct = scored(r, d.key);
                                if (pct === null) return '<td class="text-muted">—</td>';
                                return `<td><span class="heat" style="${heatStyle(pct)}">${pct.toFixed(0)}</span></td>`;
                            }).join('')}
                            <td>${grandPct === null ? '<span class="text-muted">—</span>' : `<span class="heat" style="${heatStyle(grandPct)}">${grandPct.toFixed(1)}</span>`}</td>
                        </tr>`;
                        }).join('')}
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
                                const pct = rowGrand(r);
                                const level = pct === null ? '' : levelOf(pct);
                                return `<tr>
                                    <td>${esc(r.facilityName)}</td>
                                    <td>${esc(r.districtName)}</td>
                                    <td>${pct === null ? '—' : `<strong style="color:${barColor(pct)}">${pct.toFixed(1)}%</strong>`}</td>
                                    <td><a href="${openHref(r)}">${r.inProgress ? 'Continue' : esc(level)}</a></td>
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
                                const pct = rowGrand(r);
                                const when = (r.supervisionDate || r.createdAt || '').toString().slice(0, 10);
                                return `<tr>
                                    <td>${when}</td>
                                    <td>${esc(r.facilityName)}${r.inProgress ? ' <span class="text-muted">in progress</span>' : ''}<div class="text-muted">${esc(r.regionName)}</div></td>
                                    <td>${pct === null ? '—' : `<strong style="color:${barColor(pct)}">${pct.toFixed(1)}%</strong>`}</td>
                                    <td><a href="${openHref(r)}">${r.inProgress ? 'Continue' : 'Open'}</a></td>
                                </tr>`;
                            }).join('')}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>`;

    try {
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
        const pct = rowGrand(r);
        if (pct === null) return;
        months.get(key).push(pct);
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
    } catch (err) {
        console.error(err);
    }
}

async function load() {
    const [res, structureRes] = await Promise.all([
        fetch('/api/rhspars/assessments', { headers: authHeaders() }),
        fetch('/api/rhspars/structure', { headers: authHeaders() })
    ]);
    if (res.status === 401 || structureRes.status === 401) { logout(); return; }
    if (!res.ok) {
        document.getElementById('dashboard').innerHTML = '<div class="alert alert-danger">Could not load visits.</div>';
        return;
    }
    const data = await res.json();
    const structure = structureRes.ok ? await structureRes.json() : [];
    allVisits = Array.isArray(data) ? data : [];
    const draftRow = scoreDraft(Array.isArray(structure) ? structure : []);
    if (draftRow) allVisits = [draftRow, ...allVisits];
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
