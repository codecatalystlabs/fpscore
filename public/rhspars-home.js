sessionStorage.setItem('active_tool', 'rhspars');

const DRAFT_KEY = 'rhspars_draft';

function getDraft() {
    try {
        return JSON.parse(sessionStorage.getItem(DRAFT_KEY) || '{}');
    } catch (e) {
        return {};
    }
}

function saveDraft(patch) {
    const draft = { ...getDraft(), ...patch };
    sessionStorage.setItem(DRAFT_KEY, JSON.stringify(draft));
    return draft;
}

function answeredCountForDomain(draft, domain) {
    const responses = draft.responses || {};
    let total = 0;
    let answered = 0;
    (domain.thematicAreas || []).forEach(ta => {
        (ta.questions || []).forEach(q => {
            total++;
            if (responses[String(q.id)] || responses[q.id]) answered++;
        });
    });
    return { total, answered };
}

async function api(path) {
    const res = await fetch(path, { headers: authHeaders() });
    if (res.status === 401) { logout(); throw new Error('Unauthorized'); }
    if (!res.ok) throw new Error(await res.text());
    const data = await res.json();
    return Array.isArray(data) ? data : [];
}

function fillSelect(el, items, placeholder, selectedValue) {
    const list = Array.isArray(items) ? items : [];
    el.innerHTML = `<option value="">${placeholder}</option>` +
        list.map(i => `<option value="${i.id}">${i.name}</option>`).join('');
    el.disabled = false;
    if (selectedValue) el.value = String(selectedValue);
}

function syncVisitFieldsToDraft() {
    const facilityEl = document.getElementById('facilitySelect');
    saveDraft({
        regionId: document.getElementById('regionSelect').value || '',
        districtId: document.getElementById('districtSelect').value || '',
        subcountyId: document.getElementById('subcountySelect').value || '',
        facilityId: facilityEl.value ? parseInt(facilityEl.value, 10) : null,
        facilityName: facilityEl.selectedOptions[0]?.text || '',
        supervisionType: document.getElementById('supervisionType').value,
        supervisionDate: document.getElementById('supervisionDate').value,
        nextSupervisionDate: document.getElementById('nextSupervisionDate').value,
        rhSupplier: document.getElementById('rhSupplier').value,
        facilityTypeDetail: document.getElementById('facilityTypeDetail').value,
        facilityOwnership: document.getElementById('facilityOwnership').value,
        primaryContactName: document.getElementById('primaryContactName').value,
        primaryContactPhone: document.getElementById('primaryContactPhone').value
    });
    renderDomainCards();
}

let domainsCache = [];

async function loadDomains() {
    domainsCache = await api('/api/rhspars/structure');
    // structure returns array of domains (not empty-array helper needed if object)
    if (!Array.isArray(domainsCache)) domainsCache = [];
    renderDomainCards();
}

function renderDomainCards() {
    const draft = getDraft();
    const hasFacility = !!draft.facilityId;
    const box = document.getElementById('domainCards');
    if (!domainsCache.length) {
        box.innerHTML = '<div class="col-12 text-muted small">No domains loaded.</div>';
        return;
    }

    let domainsDone = 0;
    box.innerHTML = domainsCache.map(d => {
        const { total, answered } = answeredCountForDomain(draft, d);
        let status = 'Not started';
        let cls = '';
        if (!hasFacility) cls = 'disabled-card';
        else if (answered === 0) status = `${total} questions`;
        else if (answered >= total) {
            status = `Complete · ${answered}/${total}`;
            cls = 'done';
            domainsDone++;
        } else {
            status = `In progress · ${answered}/${total}`;
            cls = 'partial';
        }
        return `
            <div class="col-md-4 col-lg-4">
                <div class="moh-card domain-card ${cls}" data-domain-id="${d.id}" data-domain-code="${d.code}">
                    <div class="d-flex justify-content-between align-items-start">
                        <div class="moh-tile-icon" style="width:28px;height:28px;font-size:0.85rem;background:#e8f1fa;color:#1a6bb5">
                            <i class="bi bi-folder2-open"></i>
                        </div>
                        <span class="small text-muted">${status}</span>
                    </div>
                    <h3 class="mt-2 mb-1" style="font-size:0.95rem">${d.name}</h3>
                    <p class="small text-muted mb-2">${(d.thematicAreas || []).length} thematic areas</p>
                    <span class="moh-tile-cta" style="color:var(--moh-accent-blue)">Open thematic areas →</span>
                </div>
            </div>`;
    }).join('');

    document.getElementById('domainProgress').textContent = hasFacility
        ? `${domainsDone}/${domainsCache.length} domains complete`
        : 'Select a facility first';

    const responses = draft.responses || {};
    const anyAnswers = Object.keys(responses).length > 0;
    document.getElementById('submitVisitBtn').disabled = !(hasFacility && anyAnswers);

    box.querySelectorAll('.domain-card:not(.disabled-card)').forEach(card => {
        card.addEventListener('click', () => {
            syncVisitFieldsToDraft();
            if (!getDraft().facilityId) {
                alert('Select a facility first.');
                return;
            }
            const domainId = card.getAttribute('data-domain-id');
            window.location.href = 'rhspars-assessment.html?domainId=' + domainId;
        });
    });
}

async function initGeo() {
    const draft = getDraft();
    const regionSelect = document.getElementById('regionSelect');
    const districtSelect = document.getElementById('districtSelect');
    const subcountySelect = document.getElementById('subcountySelect');
    const facilitySelect = document.getElementById('facilitySelect');

    const regions = await api('/api/regions');
    fillSelect(regionSelect, regions, 'Select region', draft.regionId);

    ['supervisionType','supervisionDate','nextSupervisionDate','rhSupplier','facilityTypeDetail','facilityOwnership','primaryContactName','primaryContactPhone']
        .forEach(id => {
            const el = document.getElementById(id);
            if (el && draft[id]) el.value = draft[id];
        });

    regionSelect.addEventListener('change', async function() {
        districtSelect.innerHTML = '<option value="">Select district</option>';
        districtSelect.disabled = true;
        subcountySelect.innerHTML = '<option value="">Select subcounty</option>';
        subcountySelect.disabled = true;
        facilitySelect.innerHTML = '<option value="">Select facility</option>';
        facilitySelect.disabled = true;
        saveDraft({
            regionId: this.value || '',
            districtId: '',
            subcountyId: '',
            facilityId: null,
            facilityName: '',
            responses: getDraft().responses || {}
        });
        renderDomainCards();
        if (!this.value) return;
        const districts = await api('/api/districts/' + this.value);
        fillSelect(districtSelect, districts, 'Select district');
    });

    districtSelect.addEventListener('change', async function() {
        subcountySelect.innerHTML = '<option value="">Select subcounty</option>';
        subcountySelect.disabled = true;
        facilitySelect.innerHTML = '<option value="">Select facility</option>';
        facilitySelect.disabled = true;
        saveDraft({ districtId: this.value || '', subcountyId: '', facilityId: null, facilityName: '' });
        renderDomainCards();
        if (!this.value) return;
        const subs = await api('/api/subcounties/' + this.value);
        fillSelect(subcountySelect, subs, 'Select subcounty');
    });

    subcountySelect.addEventListener('change', async function() {
        facilitySelect.innerHTML = '<option value="">Select facility</option>';
        facilitySelect.disabled = true;
        saveDraft({ subcountyId: this.value || '', facilityId: null, facilityName: '' });
        renderDomainCards();
        if (!this.value) return;
        const facilities = await api('/api/facilities/' + this.value);
        fillSelect(facilitySelect, facilities, 'Select facility');
    });

    facilitySelect.addEventListener('change', syncVisitFieldsToDraft);
    ['supervisionType','supervisionDate','nextSupervisionDate','rhSupplier','facilityTypeDetail','facilityOwnership','primaryContactName','primaryContactPhone']
        .forEach(id => document.getElementById(id).addEventListener('change', syncVisitFieldsToDraft));

    // Restore cascade from draft
    if (draft.regionId) {
        const districts = await api('/api/districts/' + draft.regionId);
        fillSelect(districtSelect, districts, 'Select district', draft.districtId);
        if (draft.districtId) {
            const subs = await api('/api/subcounties/' + draft.districtId);
            fillSelect(subcountySelect, subs, 'Select subcounty', draft.subcountyId);
            if (draft.subcountyId) {
                const facilities = await api('/api/facilities/' + draft.subcountyId);
                fillSelect(facilitySelect, facilities, 'Select facility', draft.facilityId);
            }
        }
    }

    renderDomainCards();

    // Apply admin-area UI lock after options exist
    if (typeof applyAdminRestrictions === 'function') {
        if (typeof window !== 'undefined') window.restrictionsApplied = false;
        await applyAdminRestrictions();
    }
}

document.getElementById('submitVisitBtn').addEventListener('click', async function() {
    syncVisitFieldsToDraft();
    const draft = getDraft();
    if (!draft.facilityId) {
        alert('Select a facility.');
        return;
    }
    const responses = draft.responses || {};
    if (!Object.keys(responses).length) {
        alert('Complete at least one thematic area before submitting.');
        return;
    }

    let incomplete = 0;
    domainsCache.forEach(d => {
        const { total, answered } = answeredCountForDomain(draft, d);
        if (answered < total) incomplete++;
    });
    if (incomplete > 0) {
        if (!confirm(`${incomplete} domain(s) still incomplete. Unanswered items are omitted from scoring. Submit anyway?`)) {
            return;
        }
    }

    this.disabled = true;
    this.textContent = 'Submitting…';

    try {
        const body = {
            facilityId: draft.facilityId,
            supervisionType: draft.supervisionType,
            supervisionDate: draft.supervisionDate,
            nextSupervisionDate: draft.nextSupervisionDate,
            rhSupplier: draft.rhSupplier,
            facilityTypeDetail: draft.facilityTypeDetail,
            facilityOwnership: draft.facilityOwnership,
            primaryContactName: draft.primaryContactName,
            primaryContactPhone: draft.primaryContactPhone,
            completedByName: draft.completedByName || '',
            completedByDesignation: draft.completedByDesignation || '',
            completedByInstitution: draft.completedByInstitution || '',
            assessmentCompletionDate: draft.assessmentCompletionDate || '',
            notes: draft.notes || '',
            responses,
            actionPlan: draft.actionPlan || [],
            attendees: draft.attendees || []
        };
        const res = await fetch('/api/rhspars/assessments', {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify(body)
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err.error || 'Submit failed');
        }
        const data = await res.json();
        sessionStorage.removeItem(DRAFT_KEY);
        window.location.href = 'rhspars-view.html?id=' + data.id;
    } catch (e) {
        alert(e.message);
        this.disabled = false;
        this.innerHTML = '<i class="bi bi-check2"></i> Submit visit';
    }
});

async function loadRecent() {
    try {
        const rows = await api('/api/rhspars/assessments');
        const box = document.getElementById('recentList');
        if (!rows.length) {
            box.textContent = 'No RH SPARS visits yet.';
            return;
        }
        box.innerHTML = `<table class="table table-sm mb-0"><thead><tr>
            <th>Facility</th><th>Date</th><th>Grand %</th><th>Level</th><th></th></tr></thead><tbody>
            ${rows.slice(0, 10).map(r => {
                const level = r.performanceLevel || performanceLevelFromPct(Number(r.grandPercentage));
                return `<tr>
                <td>${r.facilityName}</td>
                <td>${(r.supervisionDate || r.createdAt || '').toString().slice(0,10)}</td>
                <td><strong>${Number(r.grandPercentage).toFixed(1)}%</strong></td>
                <td><span class="badge ${performanceBadgeClass(level)}">${level}</span></td>
                <td><a href="rhspars-view.html?id=${r.id}">Open</a></td>
            </tr>`;
            }).join('')}</tbody></table>`;
    } catch (e) {
        document.getElementById('recentList').textContent = 'Unable to load history.';
    }
}

// Override api for structure (returns domains array) — loadRecent uses assessments which may not be array of plain objects with same helper
async function loadRecentSafe() {
    try {
        const res = await fetch('/api/rhspars/assessments', { headers: authHeaders() });
        if (res.status === 401) { logout(); return; }
        const rows = await res.json();
        const list = Array.isArray(rows) ? rows : [];
        const box = document.getElementById('recentList');
        if (!list.length) {
            box.textContent = 'No RH SPARS visits yet.';
            return;
        }
        box.innerHTML = `<table class="table table-sm mb-0"><thead><tr>
            <th>Facility</th><th>Date</th><th>Grand %</th><th>Level</th><th></th></tr></thead><tbody>
            ${list.slice(0, 10).map(r => {
                const level = r.performanceLevel || performanceLevelFromPct(Number(r.grandPercentage));
                return `<tr>
                <td>${r.facilityName}</td>
                <td>${(r.supervisionDate || r.createdAt || '').toString().slice(0,10)}</td>
                <td><strong>${Number(r.grandPercentage).toFixed(1)}%</strong></td>
                <td><span class="badge ${performanceBadgeClass(level)}">${level}</span></td>
                <td><a href="rhspars-view.html?id=${r.id}">Open</a></td>
            </tr>`;
            }).join('')}</tbody></table>`;
    } catch (e) {
        document.getElementById('recentList').textContent = 'Unable to load history.';
    }
}

if (!localStorage.getItem('token')) window.location.href = '/';
Promise.all([initGeo(), loadDomains(), loadRecentSafe()]).catch(console.error);
