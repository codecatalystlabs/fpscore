sessionStorage.setItem('active_tool', 'proficiency');

const API_BASE = '/api';
let previewRows = [];
let previewPage = 1;
const PREVIEW_PAGE_SIZE = 25;

if (!localStorage.getItem('token')) window.location.href = '/';

function getAuthHeaders() {
    return { Authorization: 'Bearer ' + localStorage.getItem('token') };
}

function asList(data) {
    return Array.isArray(data) ? data : [];
}

function fillSelect(el, items, placeholder, keepDisabled) {
    if (!el) return;
    const list = asList(items);
    el.innerHTML = `<option value="">${placeholder}</option>`;
    list.forEach(item => {
        const opt = document.createElement('option');
        opt.value = item.id;
        opt.textContent = item.name || item.fullName || item.full_name || ('#' + item.id);
        el.appendChild(opt);
    });
    el.disabled = !!keepDisabled;
}

function setDateRangeFromPeriod(period) {
    const start = document.getElementById('filterStartDate');
    const end = document.getElementById('filterEndDate');
    const today = new Date();
    const iso = d => d.toISOString().slice(0, 10);
    end.value = iso(today);
    if (period === '30d') {
        const s = new Date(today); s.setDate(s.getDate() - 30); start.value = iso(s);
    } else if (period === '90d') {
        const s = new Date(today); s.setDate(s.getDate() - 90); start.value = iso(s);
    } else if (period === 'ytd') {
        start.value = today.getFullYear() + '-01-01';
    } else if (period === 'all') {
        start.value = '';
        end.value = '';
    }
}

function buildFilterParams() {
    const params = new URLSearchParams();
    const map = {
        regionId: 'filterRegion',
        districtId: 'filterDistrict',
        subcountyId: 'filterSubcounty',
        facilityId: 'filterFacility',
        healthWorkerId: 'filterHealthWorker',
        assessmentTypeId: 'filterAssessmentType',
        thematicAreaId: 'filterThematicArea',
        performanceLevel: 'filterPerformance',
        startDate: 'filterStartDate',
        endDate: 'filterEndDate'
    };
    Object.keys(map).forEach(key => {
        const el = document.getElementById(map[key]);
        if (el && el.value) params.append(key, el.value);
    });
    return params;
}

async function loadBaseFilters() {
    const [regions, types] = await Promise.all([
        fetch(`${API_BASE}/regions`, { headers: getAuthHeaders() }).then(r => r.json()),
        fetch(`${API_BASE}/assessment-types`, { headers: getAuthHeaders() }).then(r => r.json())
    ]);
    const regionSelect = document.getElementById('filterRegion');
    const isRestricted = regionSelect.disabled && regionSelect.classList.contains('bg-light');
    if (!isRestricted) {
        fillSelect(regionSelect, regions, 'All regions');
    } else {
        asList(regions).forEach(r => {
            if (![...regionSelect.options].some(o => o.value === String(r.id))) {
                const opt = document.createElement('option');
                opt.value = r.id;
                opt.textContent = r.name;
                regionSelect.appendChild(opt);
            }
        });
    }
    fillSelect(document.getElementById('filterAssessmentType'), types, 'All types');
}

document.getElementById('filterPeriod').addEventListener('change', function() {
    if (this.value !== 'custom') setDateRangeFromPeriod(this.value);
});

document.getElementById('filterRegion').addEventListener('change', async function() {
    const district = document.getElementById('filterDistrict');
    const subcounty = document.getElementById('filterSubcounty');
    const facility = document.getElementById('filterFacility');
    const hw = document.getElementById('filterHealthWorker');
    fillSelect(district, [], 'All districts', true);
    fillSelect(subcounty, [], 'All subcounties', true);
    fillSelect(facility, [], 'All facilities', true);
    fillSelect(hw, [], 'All health workers', true);
    if (!this.value) return;
    const districts = await fetch(`${API_BASE}/districts/${this.value}`, { headers: getAuthHeaders() }).then(r => r.json());
    fillSelect(district, districts, 'All districts');
});

document.getElementById('filterDistrict').addEventListener('change', async function() {
    const subcounty = document.getElementById('filterSubcounty');
    const facility = document.getElementById('filterFacility');
    const hw = document.getElementById('filterHealthWorker');
    fillSelect(subcounty, [], 'All subcounties', true);
    fillSelect(facility, [], 'All facilities', true);
    fillSelect(hw, [], 'All health workers', true);
    if (!this.value) return;
    const subs = await fetch(`${API_BASE}/subcounties/${this.value}`, { headers: getAuthHeaders() }).then(r => r.json());
    fillSelect(subcounty, subs, 'All subcounties');
});

document.getElementById('filterSubcounty').addEventListener('change', async function() {
    const facility = document.getElementById('filterFacility');
    const hw = document.getElementById('filterHealthWorker');
    fillSelect(facility, [], 'All facilities', true);
    fillSelect(hw, [], 'All health workers', true);
    if (!this.value) return;
    const facilities = await fetch(`${API_BASE}/facilities/${this.value}`, { headers: getAuthHeaders() }).then(r => r.json());
    fillSelect(facility, facilities, 'All facilities');
});

document.getElementById('filterFacility').addEventListener('change', async function() {
    const hw = document.getElementById('filterHealthWorker');
    fillSelect(hw, [], 'All health workers', true);
    if (!this.value) return;
    const workers = await fetch(`${API_BASE}/health-workers?facilityId=${this.value}`, { headers: getAuthHeaders() }).then(r => r.json());
    fillSelect(hw, workers, 'All health workers');
});

document.getElementById('filterAssessmentType').addEventListener('change', async function() {
    const thematic = document.getElementById('filterThematicArea');
    fillSelect(thematic, [], 'All areas', true);
    if (!this.value) return;
    const areas = await fetch(`${API_BASE}/assessment-types/${this.value}/thematic-areas`, { headers: getAuthHeaders() }).then(r => r.json());
    fillSelect(thematic, areas, 'All areas');
});

function updateSummary(rows) {
    const n = rows.length;
    document.getElementById('statCount').textContent = n;
    document.getElementById('matchBadge').textContent = n + ' assessment' + (n === 1 ? '' : 's');
    if (!n) {
        document.getElementById('statAvg').textContent = '—';
        document.getElementById('statProf').textContent = '0';
        document.getElementById('statComp').textContent = '0';
        document.getElementById('statNA').textContent = '0';
        return;
    }
    const avg = rows.reduce((s, r) => s + Number(r.percentage || r.percentageScore || 0), 0) / n;
    document.getElementById('statAvg').textContent = avg.toFixed(1) + '%';
    document.getElementById('statProf').textContent = rows.filter(r => r.performanceLevel === 'Proficient').length;
    document.getElementById('statComp').textContent = rows.filter(r => r.performanceLevel === 'Competent').length;
    document.getElementById('statNA').textContent = rows.filter(r => r.performanceLevel === 'Not Acceptable').length;
}

function perfClass(level) {
    if (level === 'Proficient') return 'perf-proficient';
    if (level === 'Competent') return 'perf-competent';
    return 'perf-na';
}

function renderPreviewPage() {
    const wrap = document.getElementById('previewTableWrap');
    const pager = document.getElementById('previewPager');
    if (!previewRows.length) {
        wrap.innerHTML = '<p class="text-muted small mb-0">No assessments match these filters.</p>';
        pager.style.cssText = 'display:none!important';
        document.getElementById('previewMeta').textContent = '0 rows';
        return;
    }
    const totalPages = Math.max(1, Math.ceil(previewRows.length / PREVIEW_PAGE_SIZE));
    if (previewPage > totalPages) previewPage = totalPages;
    const start = (previewPage - 1) * PREVIEW_PAGE_SIZE;
    const slice = previewRows.slice(start, start + PREVIEW_PAGE_SIZE);
    wrap.innerHTML = `
        <table class="table table-sm table-hover mb-0">
            <thead><tr>
                <th>Date</th><th>Health worker</th><th>Facility</th><th>Type</th>
                <th>Score</th><th>Level</th><th></th>
            </tr></thead>
            <tbody>
                ${slice.map(r => `
                    <tr>
                        <td>${(r.createdAt || '').toString().slice(0, 10)}</td>
                        <td>${r.healthWorkerName || ''}</td>
                        <td>${r.facilityName || ''}</td>
                        <td>${r.assessmentType || r.assessmentTypeName || ''}</td>
                        <td><strong>${Number(r.percentage || r.percentageScore || 0).toFixed(1)}%</strong></td>
                        <td class="${perfClass(r.performanceLevel)}">${r.performanceLevel || ''}</td>
                        <td><a href="view-assessment.html?id=${r.id}">View</a></td>
                    </tr>
                `).join('')}
            </tbody>
        </table>`;
    pager.style.cssText = '';
    document.getElementById('pageLabel').textContent = `Page ${previewPage} of ${totalPages}`;
    document.getElementById('prevPage').disabled = previewPage <= 1;
    document.getElementById('nextPage').disabled = previewPage >= totalPages;
    document.getElementById('previewMeta').textContent =
        `Showing ${start + 1}–${Math.min(start + PREVIEW_PAGE_SIZE, previewRows.length)} of ${previewRows.length}`;
}

async function runPreview() {
    const params = buildFilterParams();
    document.getElementById('previewTableWrap').innerHTML =
        '<div class="text-center py-3"><div class="spinner-border spinner-border-sm text-success"></div></div>';
    try {
        const res = await fetch(`${API_BASE}/assessments?${params.toString()}`, { headers: getAuthHeaders() });
        if (!res.ok) throw new Error('Failed to load assessments');
        let rows = asList(await res.json());
        const perf = document.getElementById('filterPerformance').value;
        if (perf) rows = rows.filter(r => r.performanceLevel === perf);
        previewRows = rows;
        previewPage = 1;
        updateSummary(rows);
        renderPreviewPage();
    } catch (e) {
        document.getElementById('previewTableWrap').innerHTML =
            `<div class="alert alert-danger py-2 small mb-0">${e.message}</div>`;
    }
}

async function exportFile(kind) {
    const params = buildFilterParams();
    const url = `${API_BASE}/reports/assessments/${kind}?` + params.toString();
    const resp = await fetch(url, { headers: getAuthHeaders() });
    if (!resp.ok) {
        alert('Export failed: ' + (await resp.text() || resp.statusText));
        return;
    }
    const blob = await resp.blob();
    if (!blob.size) {
        alert('Export is empty — no assessments match your filters.');
        return;
    }
    const ext = kind === 'pdf' ? 'pdf' : 'xlsx';
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `assessments-report-${new Date().toISOString().slice(0, 10)}.${ext}`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    setTimeout(() => URL.revokeObjectURL(link.href), 100);
}

document.getElementById('btnPreview').addEventListener('click', runPreview);
document.getElementById('btnPdf').addEventListener('click', () => exportFile('pdf'));
document.getElementById('btnXls').addEventListener('click', () => exportFile('xls'));
document.getElementById('prevPage').addEventListener('click', () => { previewPage--; renderPreviewPage(); });
document.getElementById('nextPage').addEventListener('click', () => { previewPage++; renderPreviewPage(); });

document.getElementById('btnClear').addEventListener('click', () => {
    document.getElementById('filterPeriod').value = 'all';
    setDateRangeFromPeriod('all');
    document.getElementById('filterPerformance').value = '';
    document.getElementById('filterAssessmentType').value = '';
    fillSelect(document.getElementById('filterThematicArea'), [], 'All areas', true);
    const region = document.getElementById('filterRegion');
    if (!(region.disabled && region.classList.contains('bg-light'))) region.value = '';
    fillSelect(document.getElementById('filterDistrict'), [], 'All districts', true);
    fillSelect(document.getElementById('filterSubcounty'), [], 'All subcounties', true);
    fillSelect(document.getElementById('filterFacility'), [], 'All facilities', true);
    fillSelect(document.getElementById('filterHealthWorker'), [], 'All health workers', true);
    previewRows = [];
    updateSummary([]);
    document.getElementById('previewTableWrap').innerHTML =
        '<p class="text-muted small mb-0">No preview loaded yet.</p>';
    document.getElementById('previewPager').style.cssText = 'display:none!important';
    document.getElementById('previewMeta').textContent = 'Click Preview to load matching assessments';
});

loadBaseFilters().then(async () => {
    setTimeout(async () => {
        if (typeof window.applyAdminRestrictions === 'function') {
            try {
                const response = await fetch('/api/user/admin-areas', { headers: getAuthHeaders() });
                if (response.ok) {
                    const data = await response.json();
                    const isAdmin = data.isAdmin === true || data.isAdmin === 1 || data.isAdmin === 'true';
                    if (!isAdmin) {
                        window.restrictionsApplied = false;
                        window.applyAdminRestrictions();
                    }
                }
            } catch (e) {}
        }
    }, 250);
});
