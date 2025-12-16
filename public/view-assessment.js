// View assessment JavaScript
console.log('view-assessment.js script loaded');

const API_BASE = '/api';

function getAuthHeaders(contentType = 'application/json') {
    const headers = {};
    const token = localStorage.getItem('token');
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    } else {
        console.error('No token found in localStorage');
    }
    if (contentType) {
        headers['Content-Type'] = contentType;
    }
    return headers;
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    console.log('view-assessment.js: DOMContentLoaded fired');
    
    // Check authentication
    const token = localStorage.getItem('token');
    console.log('Token exists:', !!token);
    
    if (!token) {
        alert('Please login to view assessments');
        window.location.href = '/';
        return;
    }
    
    const urlParams = new URLSearchParams(window.location.search);
    const assessmentId = urlParams.get('id');
    console.log('Assessment ID from URL:', assessmentId);
    
    if (!assessmentId) {
        alert('No assessment ID provided');
        // Try to go back to previous page or home
        if (document.referrer) {
            window.history.back();
        } else {
            window.location.href = 'home.html';
        }
        return;
    }
    
    console.log('Calling loadAssessment and loadAssessmentSummary');
    loadAssessment(assessmentId);
    loadAssessmentSummary(assessmentId);
});

// Load assessment details
async function loadAssessment(assessmentId) {
    try {
        const headers = getAuthHeaders(null);
        console.log('Loading assessment with headers:', Object.keys(headers));
        const response = await fetch(`${API_BASE}/assessments/${assessmentId}`, {
            headers: headers
        });
        if (!response.ok) {
            if (response.status === 401) {
                alert('Your session has expired. Please login again.');
                window.location.href = '/';
                return;
            }
            if (response.status === 403) {
                throw new Error('You do not have permission to view this assessment');
            }
            throw new Error('Assessment not found');
        }
        
        const data = await response.json();
        
        if (!data || !data.assessment) {
            throw new Error('Invalid response format from server');
        }
        
        const assessment = data.assessment;
        const responses = data.responses || [];
        const thematicScores = data.thematicScores || [];
        
        // Populate assessment info
        const healthWorkerNameEl = document.getElementById('healthWorkerName');
        const facilityNameEl = document.getElementById('facilityName');
        const assessmentTypeEl = document.getElementById('assessmentType');
        const assessmentDateEl = document.getElementById('assessmentDate');
        const assessorNameEl = document.getElementById('assessorName');
        const clientNameEl = document.getElementById('clientName');
        const notesEl = document.getElementById('notes');
        
        // Helper to extract value from sql.NullString or plain string
        const getStringValue = (value) => {
            if (!value) return null;
            if (typeof value === 'string') return value;
            if (value.String) return value.String;
            return null;
        };
        
        if (healthWorkerNameEl) healthWorkerNameEl.textContent = assessment.HealthWorkerName || 'N/A';
        if (facilityNameEl) facilityNameEl.textContent = assessment.FacilityName || 'N/A';
        if (assessmentTypeEl) assessmentTypeEl.textContent = assessment.AssessmentType || 'N/A';
        if (assessmentDateEl) assessmentDateEl.textContent = new Date(assessment.CreatedAt || Date.now()).toLocaleString();
        if (assessorNameEl) assessorNameEl.textContent = getStringValue(assessment.AssessorName) || 'N/A';
        if (clientNameEl) clientNameEl.textContent = getStringValue(assessment.ClientName) || 'N/A';
        if (notesEl) notesEl.textContent = getStringValue(assessment.Notes) || 'None';
        
        // Populate overall score
        const percentage = assessment.Percentage || 0;
        const percentageScoreEl = document.getElementById('percentageScore');
        const achievedScoreEl = document.getElementById('achievedScore');
        const totalPossibleScoreEl = document.getElementById('totalPossibleScore');
        const performanceLevelEl = document.getElementById('performanceLevel');
        
        if (percentageScoreEl) percentageScoreEl.textContent = percentage.toFixed(1) + '%';
        if (achievedScoreEl) achievedScoreEl.textContent = assessment.Achieved || 0;
        if (totalPossibleScoreEl) totalPossibleScoreEl.textContent = assessment.TotalPossible || 0;
        if (performanceLevelEl) performanceLevelEl.textContent = assessment.PerformanceLevel || 'N/A';
        
        // Set score card styling
        const scoreCard = document.getElementById('overallScoreCard');
        if (scoreCard) {
            scoreCard.className = 'score-card ';
            if (percentage > 90) {
                scoreCard.className += 'score-proficient';
            } else if (percentage >= 70) {
                scoreCard.className += 'score-competent';
            } else {
                scoreCard.className += 'score-not-acceptable';
            }
        }
        
        // Populate thematic scores
        const thematicContainer = document.getElementById('thematicScoresContainer');
        if (thematicContainer) {
            thematicContainer.innerHTML = '';
            
            // Deduplicate thematic scores by id
            const seenIds = new Set();
            const uniqueThematicScores = thematicScores.filter(score => {
                if (seenIds.has(score.id)) {
                    return false;
                }
                seenIds.add(score.id);
                return true;
            });
            
            uniqueThematicScores.forEach(score => {
            const div = document.createElement('div');
            div.className = 'thematic-score';
            const progressClass = score.percentage >= 90 ? 'success' : 
                                 score.percentage >= 70 ? 'warning' : 'danger';
            div.innerHTML = `
                <div class="d-flex justify-content-between mb-1">
                    <strong>${score.name}</strong>
                    <span>${score.percentage.toFixed(1)}% (${score.achieved}/${score.possible})</span>
                </div>
                <div class="progress" style="height: 25px;">
                    <div class="progress-bar bg-${progressClass}" role="progressbar" 
                         style="width: ${score.percentage}%">
                        ${score.percentage.toFixed(1)}%
                    </div>
                </div>
            `;
                thematicContainer.appendChild(div);
            });
        }
        
        // Populate detailed responses by thematic area
        const responsesAccordion = document.getElementById('responsesAccordion');
        if (!responsesAccordion) {
            console.error('responsesAccordion element not found');
            return;
        }
        responsesAccordion.innerHTML = '';
        
        // Deduplicate responses by question id, then group by thematic area
        const seenQuestionIds = new Set();
        const uniqueResponses = responses.filter(response => {
            const questionId = response.questionId || response.question_id;
            if (questionId && seenQuestionIds.has(questionId)) {
                return false;
            }
            if (questionId) {
                seenQuestionIds.add(questionId);
            }
            return true;
        });
        
        // Group responses by thematic area
        const responsesByArea = {};
        uniqueResponses.forEach(response => {
            const areaId = response.thematicAreaId;
            if (!responsesByArea[areaId]) {
                responsesByArea[areaId] = {
                    name: response.thematicAreaName,
                    responses: []
                };
            }
            responsesByArea[areaId].responses.push(response);
        });
        
        // Create accordion items
        Object.keys(responsesByArea).forEach((areaId, index) => {
            const area = responsesByArea[areaId];
            const item = document.createElement('div');
            item.className = 'accordion-item';
            
            const headingId = `heading${areaId}`;
            const collapseId = `collapse${areaId}`;
            
            item.innerHTML = `
                <h2 class="accordion-header" id="${headingId}">
                    <button class="accordion-button ${index === 0 ? '' : 'collapsed'}" type="button" 
                            data-bs-toggle="collapse" data-bs-target="#${collapseId}">
                        ${area.name}
                    </button>
                </h2>
                <div id="${collapseId}" class="accordion-collapse collapse ${index === 0 ? 'show' : ''}" 
                     data-bs-parent="#responsesAccordion">
                    <div class="accordion-body">
                        <table class="table table-sm">
                            <thead>
                                <tr>
                                    <th>Question</th>
                                    <th>Response</th>
                                    <th>Points</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${area.responses.map(r => `
                                    <tr>
                                        <td>${r.questionText}</td>
                                        <td>
                                            <span class="badge bg-${r.response === 'Yes' ? 'success' : r.response === 'No' ? 'danger' : 'secondary'}">
                                                ${r.response}
                                            </span>
                                        </td>
                                        <td>${r.pointsEarned}/${r.scoreWeight}</td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </div>
            `;
            responsesAccordion.appendChild(item);
        });
        
    } catch (error) {
        console.error('Error loading assessment:', error);
        const errorMsg = error.message || 'Unknown error occurred';
        alert(`Error loading assessment details: ${errorMsg}`);
        
        // Try to go back to previous page or home
        if (document.referrer && document.referrer.includes(window.location.hostname)) {
            window.history.back();
        } else {
            window.location.href = 'home.html';
        }
    }
}

// Load assessment summary (good/bad contributions)
async function loadAssessmentSummary(assessmentId) {
    console.log('loadAssessmentSummary called with ID:', assessmentId);
    try {
        const headers = getAuthHeaders(null);
        console.log('Loading summary with headers:', Object.keys(headers));
        console.log('Fetching URL:', `${API_BASE}/assessments/${assessmentId}/summary`);
        
        const response = await fetch(`${API_BASE}/assessments/${assessmentId}/summary`, {
            headers: headers
        });
        
        console.log('Response status:', response.status, response.statusText);
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Response error:', errorText);
            if (response.status === 401) {
                console.error('Unauthorized: Token may be expired');
                return;
            }
            if (response.status === 403) {
                console.error('Permission denied for summary');
                return;
            }
            throw new Error(`Summary not found: ${response.status} - ${errorText}`);
        }
        
        const data = await response.json();
        console.log('Summary data received:', data);
        console.log('Yes responses:', data.yesResponses);
        console.log('No responses:', data.noResponses);
        console.log('N/A responses:', data.naResponses);
        
        // Helper function to group responses by thematic area
        function groupByThematicArea(responses) {
            if (!responses || !Array.isArray(responses)) {
                console.log('Invalid responses array:', responses);
                return {};
            }
            const grouped = {};
            responses.forEach(item => {
                const area = item.thematicArea || item.thematic_area || 'Unknown';
                if (!grouped[area]) {
                    grouped[area] = [];
                }
                grouped[area].push(item);
            });
            return grouped;
        }
        
        // Helper function to render responses grouped by thematic area
        function renderGroupedResponses(container, grouped, responseType) {
            container.innerHTML = '';
            
            if (!grouped || Object.keys(grouped).length === 0) {
                container.innerHTML = `<p class="text-muted">No questions answered ${responseType}.</p>`;
                return;
            }
            
            Object.keys(grouped).forEach(thematicArea => {
                // Thematic area header
                const areaHeader = document.createElement('div');
                areaHeader.className = 'mb-2 mt-3';
                areaHeader.style.borderBottom = '2px solid #dee2e6';
                areaHeader.style.paddingBottom = '0.5rem';
                const headerText = document.createElement('strong');
                headerText.textContent = thematicArea;
                headerText.style.fontSize = '0.9rem';
                headerText.style.color = '#495057';
                areaHeader.appendChild(headerText);
                container.appendChild(areaHeader);
                
                // Questions in this thematic area
                grouped[thematicArea].forEach(item => {
                    const div = document.createElement('div');
                    div.className = 'contribution-item';
                    const questionText = item.question || item.questionText || 'Unknown question';
                    const points = item.points || item.scoreWeight || 0;
                    
                    if (responseType === 'Yes') {
                        div.className += ' contribution-good';
                        div.innerHTML = `
                            <strong>+${points} points</strong>
                            <p class="mb-0 mt-1">${questionText}</p>
                        `;
                    } else if (responseType === 'No') {
                        div.className += ' contribution-bad';
                        div.innerHTML = `
                            <strong>-${points} points</strong>
                            <p class="mb-0 mt-1">${questionText}</p>
                        `;
                    } else {
                        div.className += ' contribution-na';
                        div.innerHTML = `
                            <strong>${points} points (N/A)</strong>
                            <p class="mb-0 mt-1">${questionText}</p>
                        `;
                    }
                    container.appendChild(div);
                });
            });
        }
        
        // Populate Yes responses
        const yesContainer = document.getElementById('yesResponses');
        if (yesContainer) {
            const yesGrouped = groupByThematicArea(data.yesResponses);
            renderGroupedResponses(yesContainer, yesGrouped, 'Yes');
        }
        
        // Populate No responses
        const noContainer = document.getElementById('noResponses');
        if (noContainer) {
            const noGrouped = groupByThematicArea(data.noResponses);
            renderGroupedResponses(noContainer, noGrouped, 'No');
        }
        
        // Populate N/A responses
        const naContainer = document.getElementById('naResponses');
        if (naContainer) {
            const naGrouped = groupByThematicArea(data.naResponses);
            renderGroupedResponses(naContainer, naGrouped, 'N/A');
        }
        
    } catch (error) {
        console.error('Error loading summary:', error);
        console.error('Error stack:', error.stack);
        console.error('Error message:', error.message);
    }
}

