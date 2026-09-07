// Assessment form JavaScript

const API_BASE = '/api';
function getAuthHeaders(contentType = 'application/json') {
    const headers = {};
    const token = localStorage.getItem('token');
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }
    if (contentType) {
        headers['Content-Type'] = contentType;
    }
    return headers;
}
let assessmentData = {
    typeId: null,
    facilityId: null,
    facilityName: null,
    healthWorkerId: null,
    thematicAreas: [],
    responses: {}
};
let isSubmittingAssessment = false;

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    assessmentData.typeId = urlParams.get('typeId');
    assessmentData.healthWorkerId = urlParams.get('healthWorkerId');
    assessmentData.facilityId = urlParams.get('facilityId');
    assessmentData.facilityName = decodeURIComponent(urlParams.get('facilityName') || '');
    
    if (window.EventLogger) {
        window.EventLogger.log('assessment', 'page_load', {
            typeId: assessmentData.typeId,
            healthWorkerId: assessmentData.healthWorkerId,
            facilityId: assessmentData.facilityId,
            timestamp: new Date().toISOString()
        });
    }
    
    document.getElementById('facilityName').textContent = assessmentData.facilityName;
    document.getElementById('assessmentTypeName').textContent = urlParams.get('typeCode') || '';
    document.getElementById('assessmentTypeDisplay').textContent = urlParams.get('typeCode') || '';
    
    if (assessmentData.typeId && assessmentData.healthWorkerId) {
        loadHealthWorkerInfo().then(() => {
            loadThematicAreas();
        });
    } else {
        alert('Missing required parameters. Please go back and select a health worker and assessment type.');
        window.location.href = 'home.html';
    }
    
    document.getElementById('submitAssessment').addEventListener('click', submitAssessment);
});         

// Load health worker info and set it in the form
async function loadHealthWorkerInfo() {
    try {
        const response = await fetch(`${API_BASE}/health-workers/${assessmentData.healthWorkerId}`, {
            headers: getAuthHeaders(null)
        });
        if (!response.ok) {
            throw new Error('Failed to load health worker information');
        }
        const healthWorker = await response.json();
        
        const select = document.getElementById('healthWorkerSelect');
        if (select) {
            // Set the health worker in the dropdown (read-only)
            select.innerHTML = '';
            const option = document.createElement('option');
            option.value = healthWorker.id;
            option.textContent = healthWorker.fullName + (healthWorker.email ? ` (${healthWorker.email})` : '') + (healthWorker.phoneNumber ? ` - ${healthWorker.phoneNumber}` : '');
            option.selected = true;
            select.appendChild(option);
            select.disabled = true; // Make it read-only since it was selected on home page
        }
        
        // Update facility name if not already set
        if (!assessmentData.facilityName && healthWorker.facilityName) {
            assessmentData.facilityName = healthWorker.facilityName;
            assessmentData.facilityId = healthWorker.facilityId;
            document.getElementById('facilityName').textContent = healthWorker.facilityName;
        }
    } catch (error) {
        console.error('Error loading health worker:', error);
        alert(`Error loading health worker information: ${error.message || 'Unknown error'}. Please try again.`);
    }
}

// Load thematic areas and questions
let isLoadingThematicAreas = false;
async function loadThematicAreas() {
    // Prevent concurrent calls
    if (isLoadingThematicAreas) {
        return;
    }
    
    isLoadingThematicAreas = true;
    try {
        const response = await fetch(`${API_BASE}/assessment-types/${assessmentData.typeId}/thematic-areas`, { headers: getAuthHeaders(null) });
        if (!response.ok) {
            let errorMsg = 'Failed to load thematic areas';
            if (response.status === 403) {
                errorMsg = 'You do not have permission to view this assessment type';
            } else {
                try {
                    const errorData = await response.json();
                    errorMsg = errorData.error || errorMsg;
                } catch (_) {
                    errorMsg = `${errorMsg} (${response.status})`;
                }
            }
            throw new Error(errorMsg);
        }
        const thematicAreas = await response.json();
        
        // Ensure we got an array
        if (!Array.isArray(thematicAreas)) {
            throw new Error('Invalid response format: expected array of thematic areas');
        }
        
        // Deduplicate thematic areas by logical identity (name + display order),
        // not just by id, to guard against duplicate rows with different ids.
        const seenAreaKeys = new Set();
        const uniqueThematicAreas = thematicAreas.filter(area => {
            const name = String(area.name || '').trim().toLowerCase();
            const order = typeof area.displayOrder === 'number' ? area.displayOrder : '';
            const key = `${name}|${order}`;
            if (seenAreaKeys.has(key)) {
                return false;
            }
            seenAreaKeys.add(key);
            return true;
        });
        
        assessmentData.thematicAreas = uniqueThematicAreas;
        
        const accordion = document.getElementById('thematicAreasAccordion');
        if (!accordion) {
            throw new Error('Thematic areas accordion element not found');
        }
        accordion.innerHTML = '';
        
        for (let i = 0; i < uniqueThematicAreas.length; i++) {
            const area = uniqueThematicAreas[i];
            const questions = await loadQuestions(area.id);
            // Ensure questions is always an array
            area.questions = Array.isArray(questions) ? questions : [];
            
            const accordionItem = createThematicAreaAccordion(area, i);
            accordion.appendChild(accordionItem);
        }
        
        updateProgress();
    } catch (error) {
        console.error('Error loading thematic areas:', error);
        // Reset to empty array on error
        assessmentData.thematicAreas = [];
        alert(`Error loading assessment data: ${error.message || 'Unknown error'}. Please try again.`);
        // Don't call updateProgress on error - it will show 0% which is misleading
    } finally {
        isLoadingThematicAreas = false;
    }
}

// Load questions for a thematic area
async function loadQuestions(thematicAreaId) {
    try {
        const response = await fetch(`${API_BASE}/thematic-areas/${thematicAreaId}/questions`, { headers: getAuthHeaders(null) });
        if (!response.ok) {
            if (response.status === 403) {
                throw new Error('Permission denied for questions');
            }
            let errorMsg = 'Failed to load questions';
            try {
                const errorData = await response.json();
                errorMsg = errorData.error || errorMsg;
            } catch (_) {
                errorMsg = `${errorMsg} (${response.status})`;
            }
            throw new Error(errorMsg);
        }
        const questions = await response.json();
        
        // Deduplicate questions by text within a thematic area, not just by id.
        if (Array.isArray(questions)) {
            const seenQuestionKeys = new Set();
            return questions.filter(q => {
                const text = String(q.text || '').trim().toLowerCase();
                const key = `${thematicAreaId}|${text}`;
                if (seenQuestionKeys.has(key)) {
                    return false;
                }
                seenQuestionKeys.add(key);
                return true;
            });
        }
        return questions;
    } catch (error) {
        console.error('Error loading questions:', error);
        return [];
    }
}

// Create accordion item for thematic area
function createThematicAreaAccordion(area, index) {
    if (!area) {
        return document.createElement('div');
    }
    if (!Array.isArray(area.questions)) {
        area.questions = [];
    }
    const item = document.createElement('div');
    item.className = 'accordion-item';
    
    const headingId = `heading${area.id}`;
    const collapseId = `collapse${area.id}`;
    
    item.innerHTML = `
        <h2 class="accordion-header" id="${headingId}">
            <button class="accordion-button ${index === 0 ? '' : 'collapsed'}" type="button" 
                    data-bs-toggle="collapse" data-bs-target="#${collapseId}">
                ${area.name}
                <span class="badge bg-secondary ms-2" id="badge-${area.id}">0%</span>
            </button>
        </h2>
        <div id="${collapseId}" class="accordion-collapse collapse ${index === 0 ? 'show' : ''}" 
             data-bs-parent="#thematicAreasAccordion">
            <div class="accordion-body">
                <div id="questions-${area.id}">
                    ${(area.questions && Array.isArray(area.questions) ? area.questions : []).map(q => createQuestionItem(q, area.id)).join('')}
                </div>
            </div>
        </div>
    `;
    
    // Add event listeners for response buttons
    item.querySelectorAll('.response-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            const questionId = parseInt(this.dataset.questionId);
            const response = this.dataset.response;
            assessmentData.responses[questionId] = response;
            
            if (window.EventLogger) {
                window.EventLogger.log('assessment', 'question_response', {
                    questionId: questionId,
                    response: response,
                    thematicAreaId: area.id,
                    timestamp: new Date().toISOString()
                });
            }
            
            // Update button states
            this.parentElement.querySelectorAll('.response-btn').forEach(b => b.classList.remove('active'));
            this.classList.add('active');
            
            updateProgress();
        });
    });
    
    return item;
}

// Create question item HTML
function createQuestionItem(question, thematicAreaId) {
    const questionId = question.id;
    const weight = getQuestionWeight(question);
    const badgeClass = question.isCritical ? 'critical-badge' : 
                       question.isImportant ? 'important-badge' : 'normal-badge';
    const badgeText = question.isCritical ? `Critical (${weight})` : 
                     question.isImportant ? `Important (${weight})` : `Normal (${weight})`;
    
    return `
        <div class="question-item" data-question-id="${questionId}">
            <div class="d-flex justify-content-between align-items-start mb-2">
                <div class="flex-grow-1">
                    <strong>${question.text}</strong>
                </div>
                <span class="badge ${badgeClass} score-badge ms-2">${badgeText}</span>
            </div>
            <div class="response-buttons">
                <button class="btn btn-sm response-btn" data-question-id="${questionId}" data-response="Yes">
                    Yes
                </button>
                <button class="btn btn-sm response-btn" data-question-id="${questionId}" data-response="No">
                    No
                </button>
                <button class="btn btn-sm response-btn" data-question-id="${questionId}" data-response="NA">
                    N/A
                </button>
            </div>
        </div>
    `;
}

// Update progress display
function updateProgress() {
    // Ensure thematicAreas is always an array
    if (!assessmentData.thematicAreas || !Array.isArray(assessmentData.thematicAreas)) {
        assessmentData.thematicAreas = [];
    }
    
    let totalPossible = 0;
    let totalPossibleAdjusted = 0;
    let totalAchieved = 0;
    const thematicProgress = {};
    
    // Calculate scores
    assessmentData.thematicAreas.forEach(area => {
        // Ensure area.questions is an array
        if (!area.questions || !Array.isArray(area.questions)) {
            area.questions = [];
        }
        let areaPossibleAdjusted = 0;
        let areaAchieved = 0;
        
        area.questions.forEach(question => {
            const weight = getQuestionWeight(question);
            if (weight <= 0) {
                return;
            }
            totalPossible += weight;
            const response = assessmentData.responses[question.id];
            
            if (response === 'Yes') {
                areaAchieved += weight;
                totalAchieved += weight;
                areaPossibleAdjusted += weight;
            } else if (response === 'NA') {
                // Exclude NA from adjusted totals
            } else {
                areaPossibleAdjusted += weight;
            }
        });
        
        const areaPercentage = areaPossibleAdjusted > 0 ? (areaAchieved / areaPossibleAdjusted) * 100 : 0;
        totalPossibleAdjusted += areaPossibleAdjusted;
        thematicProgress[area.id] = {
            possible: areaPossibleAdjusted,
            achieved: areaAchieved,
            percentage: areaPercentage
        };
    });
    
    // Update overall score
    const overallPercentage = totalPossibleAdjusted > 0 ? (totalAchieved / totalPossibleAdjusted) * 100 : 0;
    const overallScoreEl = document.getElementById('overallScore');
    const overallProgressBarEl = document.getElementById('overallProgressBar');
    if (overallScoreEl) {
        overallScoreEl.textContent = overallPercentage.toFixed(1) + '%';
    }
    if (overallProgressBarEl) {
        overallProgressBarEl.style.width = overallPercentage + '%';
        overallProgressBarEl.textContent = `${totalAchieved}/${totalPossibleAdjusted}`;
    }
    
    // Update thematic area badges
    assessmentData.thematicAreas.forEach(area => {
        const progress = thematicProgress[area.id];
        const badge = document.getElementById(`badge-${area.id}`);
        if (badge) {
            badge.textContent = progress.percentage.toFixed(1) + '%';
            if (progress.percentage >= 90) {
                badge.className = 'badge bg-success ms-2';
            } else if (progress.percentage >= 70) {
                badge.className = 'badge bg-warning ms-2';
            } else {
                badge.className = 'badge bg-danger ms-2';
            }
        }
    });
    
    // Update thematic progress sidebar
    const container = document.getElementById('thematicProgress');
    if (container) {
        container.innerHTML = '<h6 class="mb-2">By Area</h6>';
        assessmentData.thematicAreas.forEach(area => {
        const progress = thematicProgress[area.id];
        const div = document.createElement('div');
        div.className = 'mb-2';
        div.innerHTML = `
            <div class="d-flex justify-content-between mb-1">
                <small>${area.name.substring(0, 20)}${area.name.length > 20 ? '...' : ''}</small>
                <small>${progress.percentage.toFixed(0)}%</small>
            </div>
            <div class="progress" style="height: 8px;">
                <div class="progress-bar" role="progressbar" style="width: ${progress.percentage}%"></div>
            </div>
        `;
            container.appendChild(div);
        });
    }
}

// Submit assessment
async function submitAssessment(event) {
    if (event) {
        event.preventDefault();
        event.stopPropagation();
    }
    // Guard against double-clicks / duplicate listeners posting twice
    if (isSubmittingAssessment) {
        return;
    }

    if (window.EventLogger) {
        window.EventLogger.log('assessment', 'submit_start', {
            typeId: assessmentData.typeId,
            healthWorkerId: assessmentData.healthWorkerId,
            responseCount: Object.keys(assessmentData.responses).length,
            timestamp: new Date().toISOString()
        });
    }
    
    // Ensure thematicAreas is loaded
    if (!assessmentData.thematicAreas || !Array.isArray(assessmentData.thematicAreas) || assessmentData.thematicAreas.length === 0) {
        alert('Assessment data is not loaded. Please refresh the page and try again.');
        if (window.EventLogger) {
            window.EventLogger.log('assessment', 'submit_failed', {
                reason: 'data_not_loaded',
                timestamp: new Date().toISOString()
            });
        }
        return;
    }
    
    // Validate: all questions with weight > 2 must have a response
    let mandatoryUnanswered = 0;
    assessmentData.thematicAreas.forEach(area => {
        // Ensure area.questions is an array
        if (!area.questions || !Array.isArray(area.questions)) {
            return; // Skip this area if questions are not loaded
        }
        area.questions.forEach(question => {
            const weight = getQuestionWeight(question);
            if (weight > 2 && !assessmentData.responses[question.id]) {
                mandatoryUnanswered++;
            }
        });
    });

    if (mandatoryUnanswered > 0) {
        alert(`Please answer all mandatory questions (${mandatoryUnanswered} remaining). Mandatory questions are weight > 2.`);
        if (window.EventLogger) {
            window.EventLogger.log('assessment', 'submit_failed', {
                reason: 'mandatory_questions_unanswered',
                count: mandatoryUnanswered,
                timestamp: new Date().toISOString()
            });
        }
        return;
    }
    
    // Validate health worker selection
    if (!assessmentData.healthWorkerId) {
        alert('Health worker information is missing. Please go back and select a health worker.');
        return;
    }
    const healthWorkerId = parseInt(assessmentData.healthWorkerId, 10);
    if (!Number.isInteger(healthWorkerId) || healthWorkerId <= 0) {
        alert('Health worker information is invalid. Please go back and select a health worker again.');
        return;
    }
    
    const submitBtn = document.getElementById('submitAssessment');
    isSubmittingAssessment = true;
    submitBtn.disabled = true;
    submitBtn.setAttribute('aria-busy', 'true');
    submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm"></span> Submitting...';
    
    try {
        const assessorName = document.getElementById('assessorName').value;
        const notes = document.getElementById('notes').value;
        
        const response = await fetch(`${API_BASE}/assessments`, {
            method: 'POST',
            headers: getAuthHeaders('application/json'),
            body: JSON.stringify({
                healthWorkerId: healthWorkerId,
                facilityId: parseInt(assessmentData.facilityId),
                assessmentTypeId: parseInt(assessmentData.typeId),
                assessorName: assessorName,
                notes: notes,
                responses: assessmentData.responses
            })
        });
        
        if (!response.ok) {
            throw new Error('Failed to submit assessment');
        }
        
        const result = await response.json();
        
        if (window.EventLogger) {
            window.EventLogger.log('assessment', 'submit_success', {
                assessmentId: result.id,
                percentage: result.percentage,
                performanceLevel: result.performanceLevel,
                timestamp: new Date().toISOString()
            });
        }
        
        alert(`Assessment submitted successfully!\nScore: ${result.percentage.toFixed(1)}% - ${result.performanceLevel}`);
        window.location.href = `view-assessment.html?id=${result.id}`;
    } catch (error) {
        console.error('Error submitting assessment:', error);
        
        if (window.EventLogger) {
            window.EventLogger.log('assessment', 'submit_error', {
                error: error.message,
                timestamp: new Date().toISOString()
            });
        }
        
        alert('Error submitting assessment. Please try again.');
        isSubmittingAssessment = false;
        submitBtn.disabled = false;
        submitBtn.removeAttribute('aria-busy');
        submitBtn.innerHTML = '<i class="bi bi-check-circle"></i> Submit Assessment';
    }
}

function getQuestionWeight(question) {
    const weight = Number(question.scoreWeight ?? question.score_weight ?? 0);
    if (Number.isNaN(weight)) {
        return 0;
    }
    return weight;
}

