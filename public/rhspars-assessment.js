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

    (domain.thematicAreas || []).forEach((area, ai) => {
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

        item.innerHTML = `
            <h2 class="accordion-header">
                <button class="accordion-button ${ai === 0 ? '' : 'collapsed'}" type="button"
                    data-bs-toggle="collapse" data-bs-target="#${itemId}">
                    ${area.name}
                    <span class="badge bg-light text-muted ms-2">${qs.length}</span>
                </button>
            </h2>
            <div id="${itemId}" class="accordion-collapse collapse ${ai === 0 ? 'show' : ''}" data-bs-parent="#${accordionId}">
                <div class="accordion-body py-2">${body}</div>
            </div>`;
        acc.appendChild(item);
    });

    const card = document.createElement('div');
    card.className = 'moh-card';
    card.innerHTML = `<h5 class="mb-3">${domain.name} — thematic areas</h5>`;
    card.appendChild(acc);
    root.appendChild(card);

    domainQuestionIds.forEach(applySavedState);
    updateProgress();
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

    sessionStorage.setItem(DRAFT_KEY, JSON.stringify({ ...draft, responses: merged }));
    window.location.href = 'rhspars-home.html';
}

document.getElementById('saveBtn').addEventListener('click', saveAndReturn);

if (!localStorage.getItem('token')) window.location.href = '/';
loadForm().catch(err => {
    console.error(err);
    document.getElementById('formRoot').textContent = 'Error loading thematic areas.';
});
