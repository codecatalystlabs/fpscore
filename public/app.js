// Main application JavaScript

const API_BASE = '/api';

// Load regions on page load
document.addEventListener('DOMContentLoaded', function() {
    loadRegions();
    loadAssessmentTypes();
    setupCascadingDropdowns();
    loadAssessments();
});

// Load regions
async function loadRegions() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        const response = await fetch(`${API_BASE}/regions`, { headers });
        if (!response.ok) {
            if (response.status === 403) {
                console.error('Permission denied: You do not have permission to view regions');
                return;
            }
            throw new Error('Failed to load regions');
        }
        const regions = await response.json();
        const select = document.getElementById('regionSelect');
        select.innerHTML = '<option value="">Select Region</option>';
        regions.forEach(region => {
            const option = document.createElement('option');
            option.value = region.id;
            option.textContent = region.name;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading regions:', error);
    }
}

// Setup cascading dropdowns
function setupCascadingDropdowns() {
    const regionSelect = document.getElementById('regionSelect');
    const districtSelect = document.getElementById('districtSelect');
    const subcountySelect = document.getElementById('subcountySelect');
    const facilitySelect = document.getElementById('facilitySelect');

    regionSelect.addEventListener('change', async function() {
        const regionId = this.value;
        if (regionId) {
            districtSelect.disabled = false;
            districtSelect.innerHTML = '<option value="">Loading...</option>';
            subcountySelect.disabled = true;
            subcountySelect.innerHTML = '<option value="">Select Subcounty</option>';
            facilitySelect.disabled = true;
            facilitySelect.innerHTML = '<option value="">Select Facility</option>';
            
            try {
                const token = localStorage.getItem('token');
                const headers = {};
                if (token) {
                    headers['Authorization'] = `Bearer ${token}`;
                }
                const response = await fetch(`${API_BASE}/districts/${regionId}`, { headers });
                if (!response.ok) {
                    if (response.status === 403) {
                        console.error('Permission denied');
                        districtSelect.innerHTML = '<option value="">Permission denied</option>';
                        return;
                    }
                    throw new Error('Failed to load districts');
                }
                const districts = await response.json();
                districtSelect.innerHTML = '<option value="">Select District</option>';
                districts.forEach(district => {
                    const option = document.createElement('option');
                    option.value = district.id;
                    option.textContent = district.name;
                    districtSelect.appendChild(option);
                });
            } catch (error) {
                console.error('Error loading districts:', error);
            }
        } else {
            districtSelect.disabled = true;
            districtSelect.innerHTML = '<option value="">Select District</option>';
        }
    });

    districtSelect.addEventListener('change', async function() {
        const districtId = this.value;
        if (districtId) {
            subcountySelect.disabled = false;
            subcountySelect.innerHTML = '<option value="">Loading...</option>';
            facilitySelect.disabled = true;
            facilitySelect.innerHTML = '<option value="">Select Facility</option>';
            
            try {
                const token = localStorage.getItem('token');
                const headers = {};
                if (token) {
                    headers['Authorization'] = `Bearer ${token}`;
                }
                const response = await fetch(`${API_BASE}/subcounties/${districtId}`, { headers });
                if (!response.ok) {
                    if (response.status === 403) {
                        console.error('Permission denied');
                        subcountySelect.innerHTML = '<option value="">Permission denied</option>';
                        return;
                    }
                    throw new Error('Failed to load subcounties');
                }
                const subcounties = await response.json();
                subcountySelect.innerHTML = '<option value="">Select Subcounty</option>';
                subcounties.forEach(subcounty => {
                    const option = document.createElement('option');
                    option.value = subcounty.id;
                    option.textContent = subcounty.name;
                    subcountySelect.appendChild(option);
                });
            } catch (error) {
                console.error('Error loading subcounties:', error);
            }
        } else {
            subcountySelect.disabled = true;
            subcountySelect.innerHTML = '<option value="">Select Subcounty</option>';
        }
    });

    subcountySelect.addEventListener('change', async function() {
        const subcountyId = this.value;
        if (subcountyId) {
            facilitySelect.disabled = false;
            facilitySelect.innerHTML = '<option value="">Loading...</option>';
            
            try {
                const response = await fetch(`${API_BASE}/facilities/${subcountyId}`);
                const facilities = await response.json();
                facilitySelect.innerHTML = '<option value="">Select Facility</option>';
                facilities.forEach(facility => {
                    const option = document.createElement('option');
                    option.value = facility.id;
                    option.textContent = facility.name;
                    facilitySelect.appendChild(option);
                });
            } catch (error) {
                console.error('Error loading facilities:', error);
            }
        } else {
            facilitySelect.disabled = true;
            facilitySelect.innerHTML = '<option value="">Select Facility</option>';
        }
    });
}

// Load assessment types
async function loadAssessmentTypes() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        const response = await fetch(`${API_BASE}/assessment-types`, { headers });
        if (!response.ok) {
            if (response.status === 403) {
                console.error('Permission denied: You do not have permission to view assessment types');
                const container = document.getElementById('assessmentTypesContainer');
                if (container) {
                    container.innerHTML = '<div class="alert alert-warning">You do not have permission to view assessment types. Please contact your administrator.</div>';
                }
                return;
            }
            throw new Error('Failed to load assessment types');
        }
        const types = await response.json();
        const container = document.getElementById('assessmentTypesContainer');
        if (!container) return;
        container.innerHTML = '';
        
        types.forEach(type => {
            const col = document.createElement('div');
            col.className = 'col-md-6 col-lg-4 mb-3';
            col.innerHTML = `
                <div class="card assessment-type-card h-100" data-type-id="${type.id}" data-type-code="${type.code}">
                    <div class="card-body">
                        <h6 class="card-title">${type.name}</h6>
                        <button class="btn btn-primary btn-sm w-100">
                            Start Assessment
                        </button>
                    </div>
                </div>
            `;
            container.appendChild(col);
        });

        // Add click handlers
        container.querySelectorAll('.assessment-type-card').forEach(card => {
            card.addEventListener('click', function() {
                const facilityId = document.getElementById('facilitySelect').value;
                if (!facilityId) {
                    alert('Please select a facility first');
                    return;
                }
                const typeId = this.dataset.typeId;
                const typeCode = this.dataset.typeCode;
                const facilityName = document.getElementById('facilitySelect').selectedOptions[0].textContent;
                window.location.href = `assessment.html?typeId=${typeId}&facilityId=${facilityId}&facilityName=${encodeURIComponent(facilityName)}&typeCode=${typeCode}`;
            });
        });
    } catch (error) {
        console.error('Error loading assessment types:', error);
    }
}

// Load assessments for history
async function loadAssessments() {
    const container = document.getElementById('assessmentsList');
    if (!container) {
        console.error('assessmentsList container not found');
        return;
    }

    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        const response = await fetch(`${API_BASE}/assessments`, { headers });
        if (!response.ok) {
            if (response.status === 403) {
                console.error('Permission denied: You do not have permission to view assessments');
                container.innerHTML = '<div class="alert alert-warning">You do not have permission to view assessments. Please contact your administrator.</div>';
                return;
            }
            // Try to get error message from response
            let errorMsg = 'Failed to load assessments';
            try {
                const errorData = await response.json();
                errorMsg = errorData.error || errorMsg;
            } catch (e) {
                errorMsg = `Failed to load assessments (${response.status})`;
            }
            container.innerHTML = `<div class="alert alert-danger">${errorMsg}</div>`;
            return;
        }
        const assessments = await response.json();
        
        if (assessments.length === 0) {
            container.innerHTML = '<p class="text-center text-muted">No assessments found.</p>';
            return;
        }

        container.innerHTML = '<div class="table-responsive"><table class="table table-hover"><thead><tr><th>Date</th><th>Facility</th><th>Type</th><th>Score</th><th>Level</th><th>Actions</th></tr></thead><tbody></tbody></table></div>';
        const tbody = container.querySelector('tbody');

        assessments.forEach(assessment => {
            const row = document.createElement('tr');
            // Handle both percentage and percentageScore fields
            const percentage = assessment.percentage !== undefined 
                ? parseFloat(assessment.percentage).toFixed(1)
                : (assessment.percentageScore !== undefined 
                    ? parseFloat(assessment.percentageScore).toFixed(1) 
                    : '0.0');
            const levelClass = assessment.performanceLevel === 'Proficient' ? 'success' :
                              assessment.performanceLevel === 'Competent' ? 'warning' : 'danger';
            
            row.innerHTML = `
                <td>${new Date(assessment.createdAt || assessment.created_at || Date.now()).toLocaleDateString()}</td>
                <td>${assessment.facilityName || assessment.facility_name || 'N/A'}</td>
                <td>${assessment.assessmentType || assessment.assessment_type || 'N/A'}</td>
                <td><strong>${percentage}%</strong></td>
                <td><span class="badge bg-${levelClass}">${assessment.performanceLevel || assessment.performance_level || 'N/A'}</span></td>
                <td>
                    <button class="btn btn-sm btn-primary view-assessment" data-id="${assessment.id}">
                        <i class="bi bi-eye"></i> View
                    </button>
                </td>
            `;
            tbody.appendChild(row);
        });

        // Add click handlers for view buttons
        container.querySelectorAll('.view-assessment').forEach(btn => {
            btn.addEventListener('click', function() {
                const assessmentId = this.dataset.id;
                window.location.href = `view-assessment.html?id=${assessmentId}`;
            });
        });
    } catch (error) {
        console.error('Error loading assessments:', error);
        const container = document.getElementById('assessmentsList');
        if (container) {
            container.innerHTML = `<div class="alert alert-danger">Error loading assessments: ${error.message || 'Unknown error'}. Please try again or contact support.</div>`;
        }
    }
}

