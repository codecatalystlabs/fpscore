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

function areaKind(area) {
    const name = (area.name || '').toLowerCase();
    if (name.includes('health information system')) return 'his';
    if (name.includes('infrastructure')) return 'infra';
    return 'normal';
}

function areasShownFor(domain) {
    const own = domain.thematicAreas || [];
    if (domain.code === 'health_information') {
        const registers = [];
        domainsCache.forEach(d => {
            if (d.code === 'anc' || d.code === 'maternity' || d.code === 'fp') {
                (d.thematicAreas || []).forEach(area => {
                    if (areaKind(area) === 'his') registers.push(area);
                });
            }
        });
        return registers.concat(own);
    }
    if (domain.code === 'general_facility') {
        const infra = [];
        domainsCache.forEach(d => {
            if (d.code === 'anc' || d.code === 'maternity' || d.code === 'fp') {
                (d.thematicAreas || []).forEach(area => {
                    if (areaKind(area) === 'infra') infra.push(area);
                });
            }
        });
        return infra.concat(own);
    }
    if (domain.code === 'anc' || domain.code === 'maternity' || domain.code === 'fp') {
        return own.filter(area => areaKind(area) === 'normal');
    }
    return own;
}

function answeredCountForDomain(draft, domain) {
    const responses = draft.responses || {};
    let total = 0;
    let answered = 0;
    areasShownFor(domain).forEach(ta => {
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

function inferFacilityType(name) {
    const n = (name || '').toLowerCase();
    if (!n || n === 'select facility') return '';
    if (n.includes('national referral')) return 'National Referral Hospital';
    if (n.includes('regional')) return 'Regional Referral Hospital';
    if (n.includes('hospital')) return 'Hospital';
    if (/\bhc\s*iv\b|health centre iv|health center iv/.test(n)) return 'HC IV';
    if (/\bhc\s*iii\b|health centre iii|health center iii/.test(n)) return 'HC III';
    if (/\bhc\s*ii\b|health centre ii|health center ii/.test(n)) return 'HC II';
    if (n.includes('health centre') || n.includes('health center')) return 'Health Centre';
    if (n.includes('clinic')) return 'Clinic';
    return '';
}

function applyFacilityType() {
    const facilityEl = document.getElementById('facilitySelect');
    const typeEl = document.getElementById('facilityTypeDetail');
    const name = facilityEl.value ? (facilityEl.selectedOptions[0]?.text || '') : '';
    typeEl.value = inferFacilityType(name);
    typeEl.readOnly = !!typeEl.value;
}

function lettersOnly(value) {
    return String(value || '').replace(/[0-9]/g, '');
}

function digitsOnly(value) {
    return String(value || '').replace(/\D/g, '');
}

function isTextOnly(value) {
    const v = String(value || '').trim();
    return /[A-Za-z]/.test(v) && !/[0-9]/.test(v);
}

function isEmail(value) {
    return /^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$/.test(String(value || '').trim());
}

function syncDateLimits() {
    const today = todayISO();
    const completion = document.getElementById('assessmentCompletionDate');
    const next = document.getElementById('nextSupervisionDate');
    completion.min = today;
    if (completion.value && completion.value < today) completion.value = today;
    const base = completion.value || today;
    next.min = plusDaysISO(base, 1);
    if (next.value && next.value < next.min) next.value = '';
}

function todayISO() {
    const d = new Date();
    return d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0') + '-' + String(d.getDate()).padStart(2, '0');
}

function plusDaysISO(iso, days) {
    const d = new Date(iso + 'T00:00:00');
    d.setDate(d.getDate() + days);
    return d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0') + '-' + String(d.getDate()).padStart(2, '0');
}

function escAttr(value) {
    return String(value || '').replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;');
}

function personRow(person) {
    const p = person || {};
    return `<div class="person-row border rounded p-2 mb-2">
        <input class="form-control form-control-sm mb-1 p-name" placeholder="Name" value="${escAttr(p.fullName)}">
        <input class="form-control form-control-sm mb-1 p-cadre" placeholder="Cadre / position" value="${escAttr(p.cadrePosition)}">
        <input class="form-control form-control-sm mb-1 p-aff" placeholder="Affiliation" value="${escAttr(p.affiliation)}">
        <div class="d-flex gap-1">
            <input class="form-control form-control-sm p-contact" placeholder="Contact" value="${escAttr(p.contact)}">
            <button type="button" class="btn btn-sm btn-outline-danger remove-row">×</button>
        </div>
    </div>`;
}

function actionRow(item) {
    const a = item || {};
    return `<div class="row g-1 action-row mb-2">
        <div class="col-md-3"><input class="form-control form-control-sm ap-gap" placeholder="Gap identified" value="${escAttr(a.gapIdentified)}"></div>
        <div class="col-md-3"><input class="form-control form-control-sm ap-action" placeholder="Action" value="${escAttr(a.actionRecommendation)}"></div>
        <div class="col-md-2"><input class="form-control form-control-sm ap-mov" placeholder="Means of verification" value="${escAttr(a.meansOfVerification)}"></div>
        <div class="col-md-2"><input class="form-control form-control-sm ap-person" placeholder="Responsible person" value="${escAttr(a.responsiblePerson)}"></div>
        <div class="col-md-1"><input class="form-control form-control-sm ap-time" placeholder="Time frame" value="${escAttr(a.timeFrame)}"></div>
        <div class="col-md-1"><button type="button" class="btn btn-sm btn-outline-danger remove-row">×</button></div>
    </div>`;
}

function collectPeople(listId, personType) {
    return [...document.querySelectorAll('#' + listId + ' .person-row')].map(row => ({
        personType,
        fullName: row.querySelector('.p-name').value.trim(),
        cadrePosition: row.querySelector('.p-cadre').value.trim(),
        affiliation: row.querySelector('.p-aff').value.trim(),
        contact: row.querySelector('.p-contact').value.trim()
    })).filter(p => p.fullName || p.cadrePosition || p.affiliation || p.contact);
}

function flattenDomainAttendees(draft) {
    const extras = draft.domainExtras || {};
    const rows = [];
    Object.keys(extras).forEach(domainId => {
        const extra = extras[domainId] || {};
        (extra.present || []).forEach(p => {
            if (p.fullName) rows.push({ ...p, personType: 'present', domainId: parseInt(domainId, 10) });
        });
        (extra.supervisors || []).forEach(p => {
            if (p.fullName) rows.push({ ...p, personType: 'supervisor', domainId: parseInt(domainId, 10) });
        });
    });
    return rows;
}

function flattenDomainActionPlans(draft) {
    const extras = draft.domainExtras || {};
    const rows = [];
    Object.keys(extras).forEach(domainId => {
        ((extras[domainId] || {}).actionPlan || []).forEach(item => {
            if (Object.values(item).some(Boolean)) {
                rows.push({ ...item, domainId: parseInt(domainId, 10) });
            }
        });
    });
    return rows;
}

function collectActionPlan() {
    return [...document.querySelectorAll('#actionPlanList .action-row')].map(row => ({
        gapIdentified: row.querySelector('.ap-gap').value.trim(),
        actionRecommendation: row.querySelector('.ap-action').value.trim(),
        meansOfVerification: row.querySelector('.ap-mov').value.trim(),
        responsiblePerson: row.querySelector('.ap-person').value.trim(),
        timeFrame: row.querySelector('.ap-time').value.trim()
    })).filter(a => Object.values(a).some(Boolean));
}

function renderPeople() {
    const draft = getDraft();
    const people = draft.attendees || [];
    const present = people.filter(p => p.personType === 'present');
    const supervisors = people.filter(p => p.personType === 'supervisor');
    document.getElementById('presentList').innerHTML = (present.length ? present : [{}]).map(personRow).join('');
    document.getElementById('supervisorList').innerHTML = (supervisors.length ? supervisors : [{}]).map(personRow).join('');
    const actions = draft.actionPlan && draft.actionPlan.length ? draft.actionPlan : [{}];
    document.getElementById('actionPlanList').innerHTML = actions.map(actionRow).join('');
}

function syncVisitFieldsToDraft() {
    const facilityEl = document.getElementById('facilitySelect');
    const observation = getDraft().observationNotes || {};
    const observationText = Object.values(observation).filter(Boolean).join('\n');
    saveDraft({
        regionId: document.getElementById('regionSelect').value || '',
        districtId: document.getElementById('districtSelect').value || '',
        subcountyId: document.getElementById('subcountySelect').value || '',
        facilityId: facilityEl.value ? parseInt(facilityEl.value, 10) : null,
        facilityName: facilityEl.value ? (facilityEl.selectedOptions[0]?.text || '') : '',
        regionName: document.getElementById('regionSelect').value ? (document.getElementById('regionSelect').selectedOptions[0]?.text || '') : '',
        districtName: document.getElementById('districtSelect').value ? (document.getElementById('districtSelect').selectedOptions[0]?.text || '') : '',
        supervisionType: document.getElementById('supervisionType').value,
        supervisionDate: document.getElementById('supervisionDate').value,
        nextSupervisionDate: document.getElementById('nextSupervisionDate').value,
        assessmentCompletionDate: document.getElementById('assessmentCompletionDate').value,
        rhSupplier: document.getElementById('rhSupplier').value,
        facilityTypeDetail: document.getElementById('facilityTypeDetail').value,
        facilityOwnership: document.getElementById('facilityOwnership').value,
        primaryContactName: document.getElementById('primaryContactName').value,
        primaryContactPhone: document.getElementById('primaryContactPhone').value,
        primaryContactEmail: document.getElementById('primaryContactEmail').value,
        completedByName: document.getElementById('completedByName').value,
        completedByDesignation: document.getElementById('completedByDesignation').value,
        completedByInstitution: document.getElementById('completedByInstitution').value,
        notes: observationText
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

    const domainRank = { anc: 1, maternity: 2, fp: 3, commodities: 4, health_information: 5, general_facility: 6 };
    const orderedDomains = [...domainsCache].sort((a, b) => (domainRank[a.code] || 99) - (domainRank[b.code] || 99));
    let domainsDone = 0;
    box.innerHTML = orderedDomains.map(d => {
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
                    <p class="small text-muted mb-2">${areasShownFor(d).length} thematic areas</p>
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

    const today = todayISO();
    const supervisionDate = document.getElementById('supervisionDate');
    supervisionDate.value = today;
    supervisionDate.readOnly = true;
    if (!draft.assessmentCompletionDate) document.getElementById('assessmentCompletionDate').value = today;
    syncDateLimits();

    ['supervisionType','nextSupervisionDate','assessmentCompletionDate','rhSupplier','facilityOwnership','primaryContactName','primaryContactPhone','primaryContactEmail','completedByName','completedByDesignation','completedByInstitution']
        .forEach(id => {
            const el = document.getElementById(id);
            if (el && draft[id]) el.value = draft[id];
        });
    syncDateLimits();

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

    facilitySelect.addEventListener('change', () => {
        applyFacilityType();
        syncVisitFieldsToDraft();
    });
    document.getElementById('assessmentCompletionDate').addEventListener('change', () => {
        syncDateLimits();
        syncVisitFieldsToDraft();
    });
    ['supervisionType','supervisionDate','nextSupervisionDate','rhSupplier','facilityOwnership','primaryContactName','primaryContactPhone','primaryContactEmail','completedByName','completedByDesignation','completedByInstitution']
        .forEach(id => document.getElementById(id).addEventListener('change', syncVisitFieldsToDraft));
    ['primaryContactName','completedByName','completedByDesignation','completedByInstitution'].forEach(id => {
        const el = document.getElementById(id);
        el.addEventListener('input', () => { el.value = lettersOnly(el.value); });
    });
    const phone = document.getElementById('primaryContactPhone');
    phone.addEventListener('input', () => { phone.value = digitsOnly(phone.value); });

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
                applyFacilityType();
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
    const missing = [];
    [
        ['facilityId', 'Facility'],
        ['supervisionType', 'Supervision type'],
        ['supervisionDate', 'Supervision date'],
        ['nextSupervisionDate', 'Next supervision date'],
        ['assessmentCompletionDate', 'Assessment completion date'],
        ['rhSupplier', 'RH supplier / warehouse'],
        ['facilityTypeDetail', 'Facility type'],
        ['facilityOwnership', 'Facility ownership'],
        ['primaryContactName', 'Primary contact person at site'],
        ['primaryContactPhone', 'Telephone'],
        ['primaryContactEmail', 'Email'],
        ['completedByName', 'Completed by name'],
        ['completedByDesignation', 'Designation'],
        ['completedByInstitution', 'Institution / organisation']
    ].forEach(([key, label]) => { if (!draft[key]) missing.push(label); });
    const today = todayISO();
    if (draft.assessmentCompletionDate && draft.assessmentCompletionDate < today) {
        missing.push('Assessment completion date cannot be before today');
    }
    if (draft.nextSupervisionDate && draft.assessmentCompletionDate && draft.nextSupervisionDate <= draft.assessmentCompletionDate) {
        missing.push('Next supervision date must be after the assessment completion date');
    }
    [
        ['primaryContactName', 'Primary contact person at site'],
        ['completedByName', 'Completed by name'],
        ['completedByDesignation', 'Designation'],
        ['completedByInstitution', 'Institution / organisation']
    ].forEach(([key, label]) => {
        if (draft[key] && !isTextOnly(draft[key])) missing.push(label + ' must be text only');
    });
    if (draft.primaryContactPhone && !/^\d{7,15}$/.test(draft.primaryContactPhone)) {
        missing.push('Telephone must be numbers only');
    }
    if (draft.primaryContactEmail && !isEmail(draft.primaryContactEmail)) {
        missing.push('Email must include letters and a valid address');
    }
    if (missing.length) {
        alert('Complete the visit details before submitting:\n• ' + missing.join('\n• '));
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
            primaryContactEmail: draft.primaryContactEmail,
            completedByName: draft.completedByName || '',
            completedByDesignation: draft.completedByDesignation || '',
            completedByInstitution: draft.completedByInstitution || '',
            assessmentCompletionDate: draft.assessmentCompletionDate || '',
            notes: draft.notes || '',
            responses,
            actionPlan: flattenDomainActionPlans(draft),
            attendees: flattenDomainAttendees(draft)
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
