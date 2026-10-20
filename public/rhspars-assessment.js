sessionStorage.setItem('active_tool', 'rhspars');

const DRAFT_KEY = 'rhspars_draft';
const responses = {}; // local for this domain; merged into draft on save
const questionWeights = {};
let totalQuestions = 0;
let domainQuestionIds = [];

function getDraft() {
    try { return JSON.parse(sessionStorage.getItem(DRAFT_KEY) || '{}'); }
    catch (e) { return {}; }
}

function setResponse(qid, value, btn) {
    responses[qid] = value;
    const group = btn.parentElement;
    group.querySelectorAll('.btn').forEach(b => {
        b.classList.remove('active-yes', 'active-no', 'active-na');
    });
    btn.classList.add(value === 'Yes' ? 'active-yes' : value === 'No' ? 'active-no' : 'active-na');
    updateProgress();
}

function updateProgress() {
    const answered = Object.keys(responses).length;
    const pct = totalQuestions ? Math.round((answered / totalQuestions) * 100) : 0;
    document.getElementById('progressBar').style.width = pct + '%';
    document.getElementById('progressLabel').textContent = answered + '/' + totalQuestions + ' (' + pct + '%)';

    let possible = 0;
    let achieved = 0;
    domainQuestionIds.forEach(qid => {
        const weight = questionWeights[qid] || 1;
        const response = responses[qid];
        if (!response) return;
        const scored = scoreResponse(response, weight);
        if (scored.counts) {
            possible += weight;
            achieved += scored.points;
        }
    });
    const domainPct = percentageScore(achieved, possible);
    const level = performanceLevelFromPct(domainPct);
    const label = document.getElementById('domainScoreLabel');
    if (label) {
        label.textContent = possible
            ? `Domain score: ${domainPct.toFixed(1)}% (${achieved}/${possible}) · ${level}`
            : 'Domain score: —';
        label.className = 'moh-badge-soft ' + performanceBadgeClass(level);
    }
}

function applySavedState(qid) {
    const val = responses[qid];
    if (!val) return;
    const row = document.querySelector(`.q-row[data-qid="${qid}"]`);
    if (!row) return;
    const btn = [...row.querySelectorAll('.btn')].find(b => b.textContent.trim() === val);
    if (btn) {
        btn.classList.add(val === 'Yes' ? 'active-yes' : val === 'No' ? 'active-no' : 'active-na');
    }
}

function areaKind(area) {
    const name = (area.name || '').toLowerCase();
    if (name.includes('health information system')) return 'his';
    if (name.includes('infrastructure')) return 'infra';
    return 'normal';
}

function areasForDomain(domain, domains) {
    const own = domain.thematicAreas || [];
    if (domain.code === 'health_information') {
        const registers = [];
        domains.forEach(d => {
            if (d.code === 'anc' || d.code === 'maternity' || d.code === 'fp') {
                (d.thematicAreas || []).forEach(area => {
                    if (areaKind(area) === 'his') registers.push(area);
                });
            }
        });
        return registers.concat(own.filter(area => areaKind(area) !== 'his'));
    }
    if (domain.code === 'general_facility') {
        const infra = [];
        domains.forEach(d => {
            if (d.code === 'anc' || d.code === 'maternity' || d.code === 'fp') {
                (d.thematicAreas || []).forEach(area => {
                    if (areaKind(area) === 'infra') infra.push(area);
                });
            }
        });
        return infra.concat(own.filter(area => areaKind(area) !== 'infra'));
    }
    if (domain.code === 'anc' || domain.code === 'maternity' || domain.code === 'fp') {
        return own.filter(area => areaKind(area) === 'normal');
    }
    return own;
}

function displayAreaName(name) {
    if (/essential equipment/i.test(name)) {
        return 'ESSENTIAL EQUIPMENT: Assess for availability and functionality of the following';
    }
    if (/emergency medicines/i.test(name)) {
        return 'EMERGENCY MEDICINES: Assessing for availability of the following';
    }
    return name;
}

function guideFor(area) {
    const name = area.name || '';
    if (/job aid|protocol/i.test(name)) {
        return 'Check that each protocol or job aid is available at the workstation. Protocols such as PPH, hypertension in pregnancy, PPFP, PAC FP, and HIV EMTCT are each scored on their own. Observe use of the job aids, or ask the provider to demonstrate if no client is present, and write what you saw in the notes box.';
    }
    if (/essential equipment/i.test(name)) {
        return 'Score Yes only when the item is available and functional. Score NA if this service is not provided.';
    }
    if (/emergency medicines/i.test(name)) {
        return 'Score Yes if the medicine is available and No if it is not. Score NA if maternity is not provided.';
    }
    if (/committee/i.test(name)) {
        return 'This section assesses the presence and functionality of facility committees. Review minutes from the previous quarter for meetings held and for discussion of RH services.';
    }
    if (/infection/i.test(name)) {
        return 'Check that infection prevention measures are in place at this service point.';
    }
    if (/commodit|stock|drug/i.test(name)) {
        return 'This section assesses stock management practices, including availability of essential commodities, correct use of stock cards, and risk of expiry.';
    }
    if (/health information|hmis|register/i.test(name)) {
        return 'If the register is available, review ten random entries for completeness and check that monthly summaries are present. These register checks are asked here so they are not repeated under each service.';
    }
    if (/signage|poster/i.test(name)) {
        return 'Check signage in English and a local language, including the day and time of the service and directions inside the facility.';
    }
    if (/infrastructure/i.test(name)) {
        return 'These space and privacy questions are asked once here, for the service they belong to, so they are not repeated inside each program area.';
    }
    return '';
}

async function loadForm() {
    const draft = getDraft();
    if (!draft.facilityId) {
        alert('Select a facility first.');
        window.location.href = 'rhspars-home.html';
        return;
    }

    const domainId = new URLSearchParams(location.search).get('domainId');
    if (!domainId) {
        window.location.href = 'rhspars-home.html';
        return;
    }

    document.getElementById('facilityLabel').textContent = draft.facilityName || ('Facility #' + draft.facilityId);
    document.getElementById('visitMeta').textContent =
        (draft.supervisionType || '') + (draft.supervisionDate ? ' · ' + draft.supervisionDate : '');

    const res = await fetch('/api/rhspars/structure', { headers: authHeaders() });
    if (!res.ok) {
        document.getElementById('formRoot').textContent = 'Failed to load structure.';
        return;
    }
    const domains = await res.json();
    const domain = domains.find(d => String(d.id) === String(domainId));
    if (!domain) {
        document.getElementById('formRoot').textContent = 'Domain not found.';
        return;
    }

    document.getElementById('domainTitle').textContent = domain.name;
    document.title = domain.name + ' — RH SPARS';

    // Prefill responses already saved for this domain's questions
    const saved = draft.responses || {};
    const root = document.getElementById('formRoot');
    root.innerHTML = '';

    const accordionId = 'acc-domain';
    const acc = document.createElement('div');
    acc.className = 'accordion';
    acc.id = accordionId;

    areasForDomain(domain, domains).forEach((area, ai) => {
        const qs = area.questions || [];
        totalQuestions += qs.length;
        const itemId = accordionId + '-' + ai;
        const item = document.createElement('div');
        item.className = 'accordion-item';

        const body = qs.map(q => {
            domainQuestionIds.push(q.id);
            questionWeights[q.id] = q.scoreWeight || q.score_weight || 1;
            if (saved[String(q.id)]) responses[q.id] = saved[String(q.id)];
            else if (saved[q.id]) responses[q.id] = saved[q.id];
            return `
                <div class="q-row" data-qid="${q.id}">
                    <div class="q-text">${q.text}</div>
                    <div class="resp-btns">
                        <button type="button" class="btn btn-outline-success" onclick="setResponse(${q.id},'Yes',this)">Yes</button>
                        <button type="button" class="btn btn-outline-danger" onclick="setResponse(${q.id},'No',this)">No</button>
                        <button type="button" class="btn btn-outline-secondary" onclick="setResponse(${q.id},'NA',this)">NA</button>
                    </div>
                </div>`;
        }).join('');
        const guide = guideFor(area);
        const guideHtml = guide
            ? `<details class="guide-box rounded px-2 py-1 mb-2"><summary class="small">Guidance</summary><p class="small mb-0 mt-1">${guide}</p></details>`
            : '';

        item.innerHTML = `
            <h2 class="accordion-header">
                <button class="accordion-button ${ai === 0 ? '' : 'collapsed'}" type="button"
                    data-bs-toggle="collapse" data-bs-target="#${itemId}">
                    ${displayAreaName(area.name)}
                    <span class="badge bg-light text-muted ms-2">${qs.length}</span>
                </button>
            </h2>
            <div id="${itemId}" class="accordion-collapse collapse ${ai === 0 ? 'show' : ''}" data-bs-parent="#${accordionId}">
                <div class="accordion-body py-2">${guideHtml}${body}</div>
            </div>`;
        acc.appendChild(item);
    });

    const card = document.createElement('div');
    card.className = 'moh-card';
    card.innerHTML = `<h5 class="mb-3">${domain.name} — thematic areas</h5>`;
    card.appendChild(acc);
    const priorNote = (draft.observationNotes || {})[domain.id] || '';
    const notes = document.createElement('div');
    notes.className = 'mt-3';
    notes.innerHTML = `<label class="form-label small">Notes — observe use of job aids, or request a demonstration if no client is available</label>
        <textarea class="form-control form-control-sm" id="observationNote" rows="2">${priorNote.replace(/</g, '&lt;')}</textarea>`;
    card.appendChild(notes);
    root.appendChild(card);
    root.appendChild(domainFollowUpCard(domain));
    bindDomainFollowUp();

    domainQuestionIds.forEach(applySavedState);
    updateProgress();
}

function escAttr(value) {
    return String(value || '').replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;');
}

function personFields(person) {
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

function actionFields(item) {
    const a = item || {};
    return `<div class="action-row border rounded p-2 mb-2">
        <div class="row g-1">
            <div class="col-md-4"><input class="form-control form-control-sm ap-gap" placeholder="Gap identified" value="${escAttr(a.gapIdentified)}"></div>
            <div class="col-md-4"><input class="form-control form-control-sm ap-action" placeholder="Action / recommendation" value="${escAttr(a.actionRecommendation)}"></div>
            <div class="col-md-4"><input class="form-control form-control-sm ap-mov" placeholder="Means of verification" value="${escAttr(a.meansOfVerification)}"></div>
            <div class="col-md-4"><input class="form-control form-control-sm ap-person" placeholder="Responsible person" value="${escAttr(a.responsiblePerson)}"></div>
            <div class="col-md-4"><input class="form-control form-control-sm ap-time" placeholder="Time frame" value="${escAttr(a.timeFrame)}"></div>
            <div class="col-md-4 text-end"><button type="button" class="btn btn-sm btn-outline-danger remove-row">Remove</button></div>
        </div>
    </div>`;
}

function domainFollowUpCard(domain) {
    const extra = ((getDraft().domainExtras || {})[domain.id]) || {};
    const present = extra.present && extra.present.length ? extra.present : [{}];
    const supervisors = extra.supervisors && extra.supervisors.length ? extra.supervisors : [{}];
    const actions = extra.actionPlan && extra.actionPlan.length ? extra.actionPlan : [{}];
    const card = document.createElement('div');
    card.className = 'moh-card';
    card.id = 'domainFollowUp';
    card.innerHTML = `
        <h5 class="mb-1">For this domain — ${domain.name}</h5>
        <p class="small text-muted">Individuals present, supervisors, and the action plan are recorded for this domain.</p>
        <div class="row g-3">
            <div class="col-md-6">
                <div class="d-flex justify-content-between align-items-center mb-1">
                    <label class="form-label mb-0">Individuals present</label>
                    <button type="button" class="btn btn-sm btn-outline-secondary" id="addPresentBtn">Add</button>
                </div>
                <div id="presentList">${present.map(personFields).join('')}</div>
            </div>
            <div class="col-md-6">
                <div class="d-flex justify-content-between align-items-center mb-1">
                    <label class="form-label mb-0">Supervisors</label>
                    <button type="button" class="btn btn-sm btn-outline-secondary" id="addSupervisorBtn">Add</button>
                </div>
                <div id="supervisorList">${supervisors.map(personFields).join('')}</div>
            </div>
        </div>
        <div class="d-flex justify-content-between align-items-center mt-3 mb-1">
            <label class="form-label mb-0">Action plan</label>
            <button type="button" class="btn btn-sm btn-outline-secondary" id="addActionBtn">Add row</button>
        </div>
        <div id="actionPlanList">${actions.map(actionFields).join('')}</div>`;
    return card;
}

function bindDomainFollowUp() {
    document.getElementById('addPresentBtn').addEventListener('click', () => {
        document.getElementById('presentList').insertAdjacentHTML('beforeend', personFields());
    });
    document.getElementById('addSupervisorBtn').addEventListener('click', () => {
        document.getElementById('supervisorList').insertAdjacentHTML('beforeend', personFields());
    });
    document.getElementById('addActionBtn').addEventListener('click', () => {
        document.getElementById('actionPlanList').insertAdjacentHTML('beforeend', actionFields());
    });
    document.getElementById('domainFollowUp').addEventListener('click', (event) => {
        const btn = event.target.closest('.remove-row');
        if (!btn) return;
        const row = btn.closest('.person-row, .action-row');
        const list = row.parentElement;
        if (list.children.length > 1) row.remove();
        else row.querySelectorAll('input').forEach(input => { input.value = ''; });
    });
}

function readPeople(listId) {
    return [...document.querySelectorAll('#' + listId + ' .person-row')].map(row => ({
        fullName: row.querySelector('.p-name').value.trim(),
        cadrePosition: row.querySelector('.p-cadre').value.trim(),
        affiliation: row.querySelector('.p-aff').value.trim(),
        contact: row.querySelector('.p-contact').value.trim()
    })).filter(p => p.fullName || p.cadrePosition || p.affiliation || p.contact);
}

function readDomainFollowUp() {
    return {
        present: readPeople('presentList'),
        supervisors: readPeople('supervisorList'),
        actionPlan: [...document.querySelectorAll('#actionPlanList .action-row')].map(row => ({
            gapIdentified: row.querySelector('.ap-gap').value.trim(),
            actionRecommendation: row.querySelector('.ap-action').value.trim(),
            meansOfVerification: row.querySelector('.ap-mov').value.trim(),
            responsiblePerson: row.querySelector('.ap-person').value.trim(),
            timeFrame: row.querySelector('.ap-time').value.trim()
        })).filter(item => Object.values(item).some(Boolean))
    };
}

function saveAndReturn() {
    const draft = getDraft();
    const merged = { ...(draft.responses || {}) };

    // Replace this domain's answers
    domainQuestionIds.forEach(id => {
        delete merged[String(id)];
        delete merged[id];
    });
    Object.keys(responses).forEach(qid => {
        merged[String(qid)] = responses[qid];
    });
    const observationNotes = { ...(draft.observationNotes || {}) };
    const domainId = new URLSearchParams(location.search).get('domainId');
    const noteEl = document.getElementById('observationNote');
    if (noteEl) observationNotes[domainId] = noteEl.value;
    const domainExtras = { ...(draft.domainExtras || {}) };
    domainExtras[domainId] = readDomainFollowUp();

    sessionStorage.setItem(DRAFT_KEY, JSON.stringify({ ...draft, responses: merged, observationNotes, domainExtras }));
    window.location.href = 'rhspars-home.html';
}

document.getElementById('saveBtn').addEventListener('click', saveAndReturn);

if (!localStorage.getItem('token')) window.location.href = '/';
loadForm().catch(err => {
    console.error(err);
    document.getElementById('formRoot').textContent = 'Error loading thematic areas.';
});
