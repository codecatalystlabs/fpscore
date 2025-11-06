// View assessment JavaScript

const API_BASE = '/api';

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    const assessmentId = urlParams.get('id');
    
    if (assessmentId) {
        loadAssessment(assessmentId);
        loadAssessmentSummary(assessmentId);
    } else {
        alert('No assessment ID provided');
        window.location.href = 'index.html';
    }
});

// Load assessment details
async function loadAssessment(assessmentId) {
    try {
        const response = await fetch(`${API_BASE}/assessments/${assessmentId}`);
        if (!response.ok) {
            throw new Error('Assessment not found');
        }
        
        const data = await response.json();
        const assessment = data.assessment;
        const responses = data.responses;
        const thematicScores = data.thematicScores;
        
        // Populate assessment info
        document.getElementById('facilityName').textContent = assessment.FacilityName;
        document.getElementById('assessmentType').textContent = assessment.AssessmentType;
        document.getElementById('assessmentDate').textContent = new Date(assessment.CreatedAt).toLocaleString();
        document.getElementById('assessorName').textContent = assessment.AssessorName.String || 'N/A';
        document.getElementById('clientName').textContent = assessment.ClientName.String || 'N/A';
        document.getElementById('notes').textContent = assessment.Notes.String || 'None';
        
        // Populate overall score
        const percentage = assessment.Percentage;
        document.getElementById('percentageScore').textContent = percentage.toFixed(1) + '%';
        document.getElementById('achievedScore').textContent = assessment.Achieved;
        document.getElementById('totalPossibleScore').textContent = assessment.TotalPossible;
        document.getElementById('performanceLevel').textContent = assessment.PerformanceLevel;
        
        // Set score card styling
        const scoreCard = document.getElementById('overallScoreCard');
        scoreCard.className = 'score-card ';
        if (percentage > 90) {
            scoreCard.className += 'score-proficient';
        } else if (percentage >= 70) {
            scoreCard.className += 'score-competent';
        } else {
            scoreCard.className += 'score-not-acceptable';
        }
        
        // Populate thematic scores
        const thematicContainer = document.getElementById('thematicScoresContainer');
        thematicContainer.innerHTML = '';
        thematicScores.forEach(score => {
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
        
        // Populate detailed responses by thematic area
        const responsesAccordion = document.getElementById('responsesAccordion');
        responsesAccordion.innerHTML = '';
        
        // Group responses by thematic area
        const responsesByArea = {};
        responses.forEach(response => {
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
        alert('Error loading assessment details');
        window.location.href = 'index.html';
    }
}

// Load assessment summary (good/bad contributions)
async function loadAssessmentSummary(assessmentId) {
    try {
        const response = await fetch(`${API_BASE}/assessments/${assessmentId}/summary`);
        if (!response.ok) {
            throw new Error('Summary not found');
        }
        
        const data = await response.json();
        
        // Populate good contributions
        const goodContainer = document.getElementById('goodContributions');
        goodContainer.innerHTML = '';
        if (data.goodContributions.length === 0) {
            goodContainer.innerHTML = '<p class="text-muted">No positive contributions recorded.</p>';
        } else {
            data.goodContributions.forEach(contrib => {
                const div = document.createElement('div');
                div.className = 'contribution-item contribution-good';
                div.innerHTML = `
                    <strong>+${contrib.points} points</strong>
                    <p class="mb-0 mt-1">${contrib.question}</p>
                    <small class="text-muted">${contrib.thematicArea}</small>
                `;
                goodContainer.appendChild(div);
            });
        }
        
        // Populate bad contributions
        const badContainer = document.getElementById('badContributions');
        badContainer.innerHTML = '';
        if (data.badContributions.length === 0) {
            badContainer.innerHTML = '<p class="text-muted">No areas for improvement.</p>';
        } else {
            data.badContributions.forEach(contrib => {
                const div = document.createElement('div');
                div.className = 'contribution-item contribution-bad';
                div.innerHTML = `
                    <strong>-${contrib.points} points (${contrib.response})</strong>
                    <p class="mb-0 mt-1">${contrib.question}</p>
                    <small class="text-muted">${contrib.thematicArea}</small>
                `;
                badContainer.appendChild(div);
            });
        }
        
    } catch (error) {
        console.error('Error loading summary:', error);
    }
}

