// Health Worker Performance JavaScript

const API_BASE = '/api';
let currentHealthWorkerId = null;

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

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    // Check authentication
    const token = localStorage.getItem('token');
    if (!token) {
        alert('Please login to view health worker performance');
        window.location.href = '/';
        return;
    }

    const urlParams = new URLSearchParams(window.location.search);
    currentHealthWorkerId = urlParams.get('id');
    const startDate = urlParams.get('startDate');
    const endDate = urlParams.get('endDate');

    if (!currentHealthWorkerId) {
        alert('No health worker ID provided');
        window.location.href = 'home.html';
        return;
    }

    // Set date filters if provided
    if (startDate) {
        document.getElementById('startDate').value = startDate;
    }
    if (endDate) {
        document.getElementById('endDate').value = endDate;
    }

    // Set report date
    document.getElementById('reportDate').textContent = new Date().toLocaleDateString('en-US', { 
        year: 'numeric', 
        month: 'long', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });

    loadPerformance(currentHealthWorkerId);

    // Add filter button handlers
    document.getElementById('filterBtn').addEventListener('click', function() {
        loadPerformance(currentHealthWorkerId);
    });

    document.getElementById('clearBtn').addEventListener('click', function() {
        document.getElementById('startDate').value = '';
        document.getElementById('endDate').value = '';
        loadPerformance(currentHealthWorkerId);
    });
});

// Load health worker performance
async function loadPerformance(healthWorkerId) {
    try {
        const startDate = document.getElementById('startDate').value;
        const endDate = document.getElementById('endDate').value;

        const params = new URLSearchParams();
        if (startDate) params.append('startDate', startDate);
        if (endDate) params.append('endDate', endDate);

        const response = await fetch(`${API_BASE}/health-workers/${healthWorkerId}/performance?${params.toString()}`, {
            headers: getAuthHeaders(null)
        });

        if (!response.ok) {
            if (response.status === 404) {
                alert('Health worker not found');
                window.location.href = 'home.html';
                return;
            }
            throw new Error('Failed to load performance data');
        }

        const data = await response.json();
        const hw = data.healthWorker;
        const thematicScores = data.thematicScores;
        const avgScore = data.averageScore;
        const assessedCount = data.assessedCount;

        // Update health worker info
        document.getElementById('healthWorkerName').textContent = hw.fullName;
        document.getElementById('facilityName').textContent = hw.facilityName || '-';
        document.getElementById('regionName').textContent = hw.regionName || '-';
        document.getElementById('districtName').textContent = hw.districtName || '-';
        document.getElementById('subcountyName').textContent = hw.subcountyName || '-';
        document.getElementById('email').textContent = hw.email || '-';
        document.getElementById('phoneNumber').textContent = hw.phoneNumber || '-';

        // Update average score in the third column
        const avgScoreElement = document.getElementById('averageScore');
        const avgScoreContainer = avgScoreElement ? avgScoreElement.closest('.col-md-4') : null;
        
        // Update the background color based on score
        if (avgScoreContainer) {
            const scoreDiv = avgScoreContainer.querySelector('div');
            if (scoreDiv) {
                if (avgScore >= 90) {
                    scoreDiv.style.background = 'linear-gradient(135deg, #28a745 0%, #20c997 100%)';
                } else if (avgScore >= 70) {
                    scoreDiv.style.background = 'linear-gradient(135deg, #ffc107 0%, #fd7e14 100%)';
                } else {
                    scoreDiv.style.background = 'linear-gradient(135deg, #dc3545 0%, #c82333 100%)';
                }
            }
        }

        if (avgScoreElement) {
            avgScoreElement.textContent = avgScore.toFixed(1) + '%';
        }
        const assessedCountElement = document.getElementById('assessedCount');
        if (assessedCountElement) {
            assessedCountElement.textContent = `Based on ${assessedCount} thematic area${assessedCount !== 1 ? 's' : ''} assessed`;
        }

        // Update thematic scores - display as cards in grid
        const container = document.getElementById('thematicScoresContainer');
        if (thematicScores.length === 0) {
            container.innerHTML = '<div class="col-12"><p class="text-center text-muted py-5">No thematic area scores found for the selected date range.</p></div>';
            return;
        }

        container.innerHTML = '';
        thematicScores.forEach(score => {
            const div = document.createElement('div');
            div.className = 'thematic-score';
            div.dataset.thematicAreaId = score.id;
            const progressClass = score.percentage >= 90 ? 'success' : 
                                 score.percentage >= 70 ? 'warning' : 'danger';
            const progressIcon = score.percentage >= 90 ? 'bi-check-circle-fill' : 
                                score.percentage >= 70 ? 'bi-exclamation-circle-fill' : 'bi-x-circle-fill';
            const borderColor = score.percentage >= 90 ? '#28a745' : 
                               score.percentage >= 70 ? '#ffc107' : '#dc3545';
            
            div.innerHTML = `
                <div class="text-center mb-3">
                    <i class="bi ${progressIcon} text-${progressClass}" style="font-size: 2rem;"></i>
                </div>
                <h6 class="text-center mb-2" style="color: #495057; font-weight: 600; font-size: 0.95rem; min-height: 2.5rem;">
                    ${score.name}
                </h6>
                <div class="text-center mb-3">
                    <small class="text-muted d-block mb-2">
                        <i class="bi bi-tag"></i> ${score.assessmentType}
                    </small>
                </div>
                <div class="text-center mb-3">
                    <div style="font-size: 2rem; font-weight: 700; color: ${borderColor};">
                        ${score.percentage.toFixed(1)}%
                    </div>
                    <small class="text-muted d-block">${score.totalAchieved}/${score.totalPossible} points</small>
                    <small class="text-muted">
                        <i class="bi bi-clipboard-check"></i> ${score.assessmentCount} assessment${score.assessmentCount !== 1 ? 's' : ''}
                    </small>
                </div>
                <div class="progress mb-2" style="height: 25px; border-radius: 8px; overflow: hidden; box-shadow: inset 0 2px 4px rgba(0,0,0,0.1);">
                    <div class="progress-bar bg-${progressClass}" role="progressbar" 
                         style="width: ${score.percentage}%; font-weight: 600; font-size: 0.85rem; display: flex; align-items: center; justify-content: center;">
                        ${score.percentage.toFixed(1)}%
                    </div>
                </div>
                <div class="text-center mt-3">
                    <small class="text-muted" style="font-size: 0.75rem;">
                        <i class="bi bi-cursor"></i> Click for details
                    </small>
                </div>
            `;
            // Add click handler
            div.addEventListener('click', function() {
                loadThematicAreaDetails(currentHealthWorkerId, score.id);
            });
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error loading performance:', error);
        alert(`Error loading performance data: ${error.message || 'Unknown error'}`);
    }
}

// Load thematic area details
async function loadThematicAreaDetails(healthWorkerId, thematicAreaId) {
    try {
        const startDate = document.getElementById('startDate').value;
        const endDate = document.getElementById('endDate').value;

        const params = new URLSearchParams();
        if (startDate) params.append('startDate', startDate);
        if (endDate) params.append('endDate', endDate);

        const response = await fetch(`${API_BASE}/health-workers/${healthWorkerId}/thematic-areas/${thematicAreaId}/details?${params.toString()}`, {
            headers: getAuthHeaders(null)
        });

        if (!response.ok) {
            throw new Error('Failed to load thematic area details');
        }

        const data = await response.json();
        const ta = data.thematicArea;
        const performance = data.performance;
        const assessments = data.assessments || [];
        const yesQuestions = data.yesQuestions || [];
        const noQuestions = data.noQuestions || [];
        const naQuestions = data.naQuestions || [];
        const allQuestions = data.allQuestions || [];

        // Update modal title
        document.getElementById('thematicAreaModalTitle').textContent = ta.name;

        // Build modal content
        const content = document.getElementById('thematicAreaDetailsContent');
        const progressClass = performance.percentage >= 90 ? 'success' : 
                             performance.percentage >= 70 ? 'warning' : 'danger';

        let html = `
            <div class="mb-4 text-center">
                <h4 class="mb-2" style="color: #667eea; font-weight: 600;">${ta.name}</h4>
                <p class="text-muted mb-0"><i class="bi bi-tag"></i> ${ta.assessmentType}</p>
            </div>

            <div class="card mb-4" style="border: none; box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);">
                <div class="card-body p-4">
                    <h5 class="card-title mb-4" style="color: #495057; font-weight: 600;">
                        <i class="bi bi-graph-up-arrow"></i> Performance Summary
                    </h5>
                    <div class="row g-3 mb-4">
                        <div class="col-md-3">
                            <div class="performance-metric" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white;">
                                <h3 style="color: white;">${performance.percentage.toFixed(1)}%</h3>
                                <small style="color: rgba(255,255,255,0.9);">Score</small>
                            </div>
                        </div>
                        <div class="col-md-3">
                            <div class="performance-metric" style="background: #f8f9fa; border: 2px solid #e9ecef;">
                                <h3 style="color: #495057;">${performance.totalAchieved}/${performance.totalPossible}</h3>
                                <small class="text-muted">Points</small>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="performance-metric" style="background: #d4edda; border: 2px solid #c3e6cb;">
                                <h3 style="color: #155724;">${performance.yesCount}</h3>
                                <small style="color: #155724;">Yes</small>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="performance-metric" style="background: #f8d7da; border: 2px solid #f5c6cb;">
                                <h3 style="color: #721c24;">${performance.noCount}</h3>
                                <small style="color: #721c24;">No</small>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="performance-metric" style="background: #e2e3e5; border: 2px solid #d6d8db;">
                                <h3 style="color: #383d41;">${performance.naCount}</h3>
                                <small style="color: #383d41;">N/A</small>
                            </div>
                        </div>
                    </div>
                    <div class="progress" style="height: 35px; border-radius: 10px; overflow: hidden; box-shadow: inset 0 2px 4px rgba(0,0,0,0.1);">
                        <div class="progress-bar bg-${progressClass}" role="progressbar" 
                             style="width: ${performance.percentage}%; font-weight: 600; font-size: 1rem; display: flex; align-items: center; justify-content: center;">
                            ${performance.percentage.toFixed(1)}%
                        </div>
                    </div>
                </div>
            </div>
        `;

        // Show all assessments for this thematic area
        html += `
            <div class="card mb-4">
                <div class="section-header" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                    <i class="bi bi-clipboard-data"></i>
                    <span>All Assessments for This Thematic Area ${assessments && assessments.length > 0 ? `(${assessments.length})` : '(0)'}</span>
                </div>
                <div class="card-body">
        `;
        
        if (assessments && assessments.length > 0) {
            assessments.forEach((assessment, index) => {
                const assessmentDate = new Date(assessment.createdAt).toLocaleDateString('en-US', { 
                    year: 'numeric', 
                    month: 'short', 
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                });
                const assessmentProgressClass = assessment.percentageScore >= 90 ? 'success' : 
                                                assessment.percentageScore >= 70 ? 'warning' : 'danger';
                const scoreColor = assessment.percentageScore >= 90 ? '#28a745' : 
                                   assessment.percentageScore >= 70 ? '#ffc107' : '#dc3545';
                html += `
                    <div class="card mb-3" style="border: 2px solid #e9ecef; border-radius: 10px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); transition: all 0.2s;">
                        <div class="card-body p-3">
                            <div class="d-flex justify-content-between align-items-start mb-2">
                                <div class="flex-grow-1">
                                    <h6 class="mb-1" style="color: #495057; font-weight: 600;">
                                        <i class="bi bi-file-earmark-check text-primary"></i> Assessment #${assessments.length - index}
                                    </h6>
                                    <small class="text-muted d-block mb-2">
                                        <i class="bi bi-calendar3"></i> ${assessmentDate}
                                    </small>
                                    ${assessment.assessmentType ? `<small class="text-muted d-block mb-1"><i class="bi bi-tag"></i> ${assessment.assessmentType}</small>` : ''}
                                </div>
                                <div class="text-end ms-3">
                                    <div style="font-size: 1.75rem; font-weight: 700; color: ${scoreColor};">
                                        ${(assessment.percentageScore || 0).toFixed(1)}%
                                    </div>
                                    <small class="text-muted d-block">${assessment.achievedScore || 0}/${assessment.possibleScore || 0} points</small>
                                </div>
                            </div>
                            <div class="progress mb-3" style="height: 25px; border-radius: 8px; overflow: hidden; box-shadow: inset 0 2px 4px rgba(0,0,0,0.1);">
                                <div class="progress-bar bg-${assessmentProgressClass}" role="progressbar" 
                                     style="width: ${assessment.percentageScore || 0}%; font-weight: 600; font-size: 0.85rem; display: flex; align-items: center; justify-content: center;">
                                    ${(assessment.percentageScore || 0).toFixed(1)}%
                                </div>
                            </div>
                            <div class="row g-2">
                                ${assessment.assessorName ? `
                                    <div class="col-12">
                                        <small class="text-muted">
                                            <i class="bi bi-person-check"></i> <strong>Assessor:</strong> ${assessment.assessorName}
                                        </small>
                                    </div>
                                ` : ''}
                                ${assessment.notes ? `
                                    <div class="col-12">
                                        <small class="text-muted">
                                            <i class="bi bi-sticky"></i> <strong>Notes:</strong> ${assessment.notes}
                                        </small>
                                    </div>
                                ` : ''}
                            </div>
                        </div>
                    </div>
                `;
            });
        } else {
            html += `
                <div class="text-center py-4">
                    <i class="bi bi-inbox" style="font-size: 3rem; color: #dee2e6;"></i>
                    <p class="text-muted mt-3">No assessments found for this thematic area in the selected date range.</p>
                </div>
            `;
        }
        html += `</div></div>`;

        // Questions answered Yes
        if (yesQuestions.length > 0) {
            html += `
                <div class="card mb-4">
                    <div class="section-header bg-success">
                        <i class="bi bi-check-circle-fill"></i>
                        <span>Questions Answered Yes (${yesQuestions.length})</span>
                    </div>
                    <div class="card-body">
            `;
            yesQuestions.forEach(q => {
                const latestResponse = q.responses && q.responses.length > 0 ? q.responses[0] : null;
                html += `
                    <div class="question-item yes">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <div class="flex-grow-1">
                                <div class="question-text mb-2">
                                    <strong>${q.text}</strong>
                                </div>
                                <div class="d-flex flex-wrap gap-2">
                                    ${q.isCritical ? '<span class="badge bg-danger"><i class="bi bi-exclamation-triangle"></i> Critical</span>' : ''}
                                    ${q.isImportant ? '<span class="badge bg-warning text-dark"><i class="bi bi-star"></i> Important</span>' : ''}
                                    <span class="badge bg-secondary"><i class="bi bi-trophy"></i> ${q.scoreWeight} points</span>
                                </div>
                            </div>
                            ${latestResponse ? `<span class="response-badge yes ms-3"><i class="bi bi-check-circle"></i> ${latestResponse.response}</span>` : ''}
                        </div>
                        ${latestResponse ? `<div class="mt-2"><small class="text-muted"><i class="bi bi-calendar"></i> Last answered: ${new Date(latestResponse.date).toLocaleDateString()}</small></div>` : ''}
                    </div>
                `;
            });
            html += `</div></div>`;
        }

        // Questions answered No
        if (noQuestions.length > 0) {
            html += `
                <div class="card mb-4">
                    <div class="section-header bg-danger">
                        <i class="bi bi-x-circle-fill"></i>
                        <span>Questions Answered No (${noQuestions.length})</span>
                    </div>
                    <div class="card-body">
            `;
            noQuestions.forEach(q => {
                const latestResponse = q.responses && q.responses.length > 0 ? q.responses[0] : null;
                html += `
                    <div class="question-item no">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <div class="flex-grow-1">
                                <div class="question-text mb-2">
                                    <strong>${q.text}</strong>
                                </div>
                                <div class="d-flex flex-wrap gap-2">
                                    ${q.isCritical ? '<span class="badge bg-danger"><i class="bi bi-exclamation-triangle"></i> Critical</span>' : ''}
                                    ${q.isImportant ? '<span class="badge bg-warning text-dark"><i class="bi bi-star"></i> Important</span>' : ''}
                                    <span class="badge bg-secondary"><i class="bi bi-trophy"></i> ${q.scoreWeight} points</span>
                                </div>
                            </div>
                            ${latestResponse ? `<span class="response-badge no ms-3"><i class="bi bi-x-circle"></i> ${latestResponse.response}</span>` : ''}
                        </div>
                        ${latestResponse ? `<div class="mt-2"><small class="text-muted"><i class="bi bi-calendar"></i> Last answered: ${new Date(latestResponse.date).toLocaleDateString()}</small></div>` : ''}
                    </div>
                `;
            });
            html += `</div></div>`;
        }

        // Questions answered N/A
        if (naQuestions.length > 0) {
            html += `
                <div class="card mb-4">
                    <div class="section-header bg-secondary">
                        <i class="bi bi-dash-circle-fill"></i>
                        <span>Questions Answered N/A (${naQuestions.length})</span>
                    </div>
                    <div class="card-body">
            `;
            naQuestions.forEach(q => {
                const latestResponse = q.responses && q.responses.length > 0 ? q.responses[0] : null;
                html += `
                    <div class="question-item na">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <div class="flex-grow-1">
                                <div class="question-text mb-2">
                                    <strong>${q.text}</strong>
                                </div>
                                <div class="d-flex flex-wrap gap-2">
                                    ${q.isCritical ? '<span class="badge bg-danger"><i class="bi bi-exclamation-triangle"></i> Critical</span>' : ''}
                                    ${q.isImportant ? '<span class="badge bg-warning text-dark"><i class="bi bi-star"></i> Important</span>' : ''}
                                    <span class="badge bg-secondary"><i class="bi bi-trophy"></i> ${q.scoreWeight} points</span>
                                </div>
                            </div>
                            ${latestResponse ? `<span class="response-badge na ms-3"><i class="bi bi-dash-circle"></i> ${latestResponse.response}</span>` : ''}
                        </div>
                        ${latestResponse ? `<div class="mt-2"><small class="text-muted"><i class="bi bi-calendar"></i> Last answered: ${new Date(latestResponse.date).toLocaleDateString()}</small></div>` : ''}
                    </div>
                `;
            });
            html += `</div></div>`;
        }

        // All questions with all responses (only show questions that have responses)
        if (allQuestions && allQuestions.length > 0) {
            html += `
                <div class="card">
                    <div class="section-header" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                        <i class="bi bi-list-ul"></i>
                        <span>All Questions and Responses</span>
                    </div>
                    <div class="card-body">
            `;
            allQuestions.forEach((q, index) => {
                if (!q.responses || q.responses.length === 0) {
                    return; // Skip questions without responses
                }
                html += `
                    <div class="question-item mb-4" style="border-left: 4px solid #667eea;">
                        <div class="mb-3">
                            <div class="d-flex justify-content-between align-items-start mb-2">
                                <div class="flex-grow-1">
                                    <div class="question-text mb-2">
                                        <strong style="color: #495057;">${q.text}</strong>
                                    </div>
                                    <div class="d-flex flex-wrap gap-2">
                                        ${q.isCritical ? '<span class="badge bg-danger"><i class="bi bi-exclamation-triangle"></i> Critical</span>' : ''}
                                        ${q.isImportant ? '<span class="badge bg-warning text-dark"><i class="bi bi-star"></i> Important</span>' : ''}
                                        <span class="badge bg-secondary"><i class="bi bi-trophy"></i> ${q.scoreWeight} points</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="responses mt-3">
                            <h6 class="mb-3" style="color: #6c757d; font-size: 0.85rem; font-weight: 600;">
                                <i class="bi bi-clock-history"></i> Response History (${q.responses.length})
                            </h6>
                `;
                q.responses.forEach(r => {
                    const badgeClass = r.response === 'Yes' ? 'yes' : r.response === 'No' ? 'no' : 'na';
                    html += `
                        <div class="response-item">
                            <div class="d-flex align-items-center gap-3">
                                <span class="response-badge ${badgeClass}">
                                    <i class="bi ${r.response === 'Yes' ? 'bi-check-circle' : r.response === 'No' ? 'bi-x-circle' : 'bi-dash-circle'}"></i>
                                    ${r.response}
                                </span>
                                <span class="text-muted" style="font-size: 0.85rem;">
                                    <i class="bi bi-calendar3"></i> ${new Date(r.date).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })}
                                </span>
                            </div>
                            <div>
                                <span class="badge bg-info text-white">
                                    <i class="bi bi-star-fill"></i> ${r.pointsEarned} points
                                </span>
                            </div>
                        </div>
                    `;
                });
                html += `</div></div>`;
                if (index < allQuestions.length - 1) {
                    html += `<hr style="margin: 1.5rem 0; opacity: 0.3;">`;
                }
            });
            html += `</div></div>`;
        }

        content.innerHTML = html;

        // Show modal
        const modal = new bootstrap.Modal(document.getElementById('thematicAreaModal'));
        modal.show();
    } catch (error) {
        console.error('Error loading thematic area details:', error);
        alert(`Error loading thematic area details: ${error.message || 'Unknown error'}`);
    }
}

