// Assessment form JavaScript

const API_BASE = '/api';
let assessmentData = {
    typeId: null,
    facilityId: null,
    facilityName: null,
    thematicAreas: [],
    responses: {}
};

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    assessmentData.typeId = urlParams.get('typeId');
    assessmentData.facilityId = urlParams.get('facilityId');
    assessmentData.facilityName = decodeURIComponent(urlParams.get('facilityName') || '');
    
    document.getElementById('facilityName').textContent = assessmentData.facilityName;
    document.getElementById('assessmentTypeName').textContent = urlParams.get('typeCode') || '';
    document.getElementById('assessmentTypeDisplay').textContent = urlParams.get('typeCode') || '';
    
    if (assessmentData.typeId) {
        loadThematicAreas();
    }
    
    document.getElementById('submitAssessment').addEventListener('click', submitAssessment);
});

// Load thematic areas and questions
async function loadThematicAreas() {
    try {
        const response = await fetch(`${API_BASE}/assessment-types/${assessmentData.typeId}/thematic-areas`);
        const thematicAreas = await response.json();
        assessmentData.thematicAreas = thematicAreas;
        
        const accordion = document.getElementById('thematicAreasAccordion');
        accordion.innerHTML = '';
        
        for (let i = 0; i < thematicAreas.length; i++) {
            const area = thematicAreas[i];
            const questions = await loadQuestions(area.id);
            area.questions = questions;
            
            const accordionItem = createThematicAreaAccordion(area, i);
            accordion.appendChild(accordionItem);
        }
        
        updateProgress();
    } catch (error) {
        console.error('Error loading thematic areas:', error);
        alert('Error loading assessment data. Please try again.');
    }
}

// Load questions for a thematic area
async function loadQuestions(thematicAreaId) {
    try {
        const response = await fetch(`${API_BASE}/thematic-areas/${thematicAreaId}/questions`);
        const questions = await response.json();
        return questions;
    } catch (error) {
        console.error('Error loading questions:', error);
        return [];
    }
}

// Create accordion item for thematic area
function createThematicAreaAccordion(area, index) {
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
                    ${area.questions.map(q => createQuestionItem(q, area.id)).join('')}
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
    const badgeClass = question.isCritical ? 'critical-badge' : 
                       question.isImportant ? 'important-badge' : 'normal-badge';
    const badgeText = question.isCritical ? 'Critical (10)' : 
                     question.isImportant ? 'Important (5)' : 'Normal (2)';
    
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
    let totalPossible = 0;
    let totalAchieved = 0;
    const thematicProgress = {};
    
    // Calculate scores
    assessmentData.thematicAreas.forEach(area => {
        let areaPossible = 0;
        let areaAchieved = 0;
        
        area.questions.forEach(question => {
            const weight = question.scoreWeight;
            areaPossible += weight;
            totalPossible += weight;
            
            if (assessmentData.responses[question.id] === 'Yes') {
                areaAchieved += weight;
                totalAchieved += weight;
            }
        });
        
        const areaPercentage = areaPossible > 0 ? (areaAchieved / areaPossible) * 100 : 0;
        thematicProgress[area.id] = {
            possible: areaPossible,
            achieved: areaAchieved,
            percentage: areaPercentage
        };
    });
    
    // Update overall score
    const overallPercentage = totalPossible > 0 ? (totalAchieved / totalPossible) * 100 : 0;
    document.getElementById('overallScore').textContent = overallPercentage.toFixed(1) + '%';
    document.getElementById('overallProgressBar').style.width = overallPercentage + '%';
    document.getElementById('overallProgressBar').textContent = `${totalAchieved}/${totalPossible}`;
    
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

// Submit assessment
async function submitAssessment() {
    // Validate: all questions with weight > 2 must have a response
    let mandatoryUnanswered = 0;
    assessmentData.thematicAreas.forEach(area => {
        area.questions.forEach(question => {
            if (question.scoreWeight > 2 && !assessmentData.responses[question.id]) {
                mandatoryUnanswered++;
            }
        });
    });

    if (mandatoryUnanswered > 0) {
        alert(`Please answer all mandatory questions (${mandatoryUnanswered} remaining). Mandatory questions are weight > 2.`);
        return;
    }
    
    const submitBtn = document.getElementById('submitAssessment');
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm"></span> Submitting...';
    
    try {
        const assessorName = document.getElementById('assessorName').value;
        const clientName = document.getElementById('clientName').value;
        const notes = document.getElementById('notes').value;
        
        const response = await fetch(`${API_BASE}/assessments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                facilityId: parseInt(assessmentData.facilityId),
                assessmentTypeId: parseInt(assessmentData.typeId),
                assessorName: assessorName,
                clientName: clientName,
                notes: notes,
                responses: assessmentData.responses
            })
        });
        
        if (!response.ok) {
            throw new Error('Failed to submit assessment');
        }
        
        const result = await response.json();
        alert(`Assessment submitted successfully!\nScore: ${result.percentage.toFixed(1)}% - ${result.performanceLevel}`);
        window.location.href = `view-assessment.html?id=${result.id}`;
    } catch (error) {
        console.error('Error submitting assessment:', error);
        alert('Error submitting assessment. Please try again.');
        submitBtn.disabled = false;
        submitBtn.innerHTML = '<i class="bi bi-check-circle"></i> Submit Assessment';
    }
}

