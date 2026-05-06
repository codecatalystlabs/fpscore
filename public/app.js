// Main application JavaScript

const API_BASE = '/api';
const HOME_SELECTION_STORAGE_KEY = 'fp_home_selection_v1';

// Log application initialization
if (window.EventLogger) {
    window.EventLogger.log('application', 'init', {
        page: 'home',
        timestamp: new Date().toISOString()
    });
}

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

// Load regions on page load
document.addEventListener('DOMContentLoaded', function() {
    if (window.EventLogger) {
        window.EventLogger.log('application', 'dom_ready', {
            page: 'home',
            timestamp: new Date().toISOString()
        });
    }
    loadRegions();
    loadAssessmentTypes();
    setupCascadingDropdowns();
    loadAssessments();
    
    // Add filter button handlers
    const filterBtn = document.getElementById('filterAssessments');
    const clearBtn = document.getElementById('clearFilters');
    if (filterBtn) {
        filterBtn.addEventListener('click', loadAssessments);
    }
    if (clearBtn) {
        clearBtn.addEventListener('click', function() {
            const startDate = document.getElementById('startDate');
            const endDate = document.getElementById('endDate');
            if (startDate) startDate.value = '';
            if (endDate) endDate.value = '';
            loadAssessments();
        });
    }
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
        if (!select) return;
        
        // Check if this dropdown is restricted - if so, don't clear it
        const isRestricted = select.disabled && select.classList.contains('bg-light');
        if (!isRestricted) {
            select.innerHTML = '<option value="">All Regions</option>';
            regions.forEach(region => {
                const option = document.createElement('option');
                option.value = region.id;
                option.textContent = region.name;
                select.appendChild(option);
            });
        } else {
            // If restricted, only add options that aren't already there
            const existingValues = Array.from(select.options).map(opt => opt.value);
            regions.forEach(region => {
                if (!existingValues.includes(String(region.id))) {
                    const option = document.createElement('option');
                    option.value = region.id;
                    option.textContent = region.name;
                    select.appendChild(option);
                }
            });
        }
    } catch (error) {
        console.error('Error loading regions:', error);
    }
}

// Global variables for health worker selection
let allHealthWorkers = []; // Store all loaded health workers for filtering
let selectedHealthWorker = null;

// Setup cascading dropdowns and health worker selection
function setupCascadingDropdowns() {
    const regionSelect = document.getElementById('regionSelect');
    const districtSelect = document.getElementById('districtSelect');
    const subcountySelect = document.getElementById('subcountySelect');
    const facilitySelect = document.getElementById('facilitySelect');
    const healthWorkerSelect = document.getElementById('healthWorkerSelect');
    const healthWorkerSearch = document.getElementById('healthWorkerSearch');
    const healthWorkerDropdown = document.getElementById('healthWorkerDropdown');
    let isRestoringSelection = false;

    if (!regionSelect || !districtSelect || !facilitySelect || !healthWorkerSelect) {
        return; // Elements don't exist on this page
    }

    function saveHomeSelectionState() {
        if (isRestoringSelection) return;
        const state = {
            regionId: regionSelect.value || '',
            districtId: districtSelect.value || '',
            subcountyId: subcountySelect ? (subcountySelect.value || '') : '',
            facilityId: facilitySelect.value || '',
            selectedHealthWorker: selectedHealthWorker,
            healthWorkerSearch: healthWorkerSearch ? (healthWorkerSearch.value || '') : ''
        };
        localStorage.setItem(HOME_SELECTION_STORAGE_KEY, JSON.stringify(state));
    }

    function loadHomeSelectionState() {
        try {
            return JSON.parse(localStorage.getItem(HOME_SELECTION_STORAGE_KEY) || '{}');
        } catch (e) {
            return {};
        }
    }

    async function restoreHomeSelectionState() {
        const saved = loadHomeSelectionState();
        if (!saved || Object.keys(saved).length === 0) return;

        isRestoringSelection = true;
        try {
            if (saved.regionId) {
                regionSelect.value = saved.regionId;
                districtSelect.disabled = false;
                const districtsResp = await fetch(`${API_BASE}/districts/${saved.regionId}`, { headers: getAuthHeaders(null) });
                if (districtsResp.ok) {
                    const districts = await districtsResp.json();
                    districtSelect.innerHTML = '<option value="">All Districts</option>';
                    districts.forEach(district => {
                        const option = document.createElement('option');
                        option.value = district.id;
                        option.textContent = district.name;
                        districtSelect.appendChild(option);
                    });
                }
            }
            if (saved.districtId && districtSelect) {
                districtSelect.value = saved.districtId;
                if (subcountySelect) {
                    subcountySelect.disabled = false;
                    const subcountyResp = await fetch(`${API_BASE}/subcounties/${saved.districtId}`, { headers: getAuthHeaders(null) });
                    if (subcountyResp.ok) {
                        const subcounties = await subcountyResp.json();
                        subcountySelect.innerHTML = '<option value="">All Subcounties</option>';
                        subcounties.forEach(subcounty => {
                            const option = document.createElement('option');
                            option.value = subcounty.id;
                            option.textContent = subcounty.name;
                            subcountySelect.appendChild(option);
                        });
                    }
                }
            }
            if (saved.subcountyId && subcountySelect) {
                subcountySelect.value = saved.subcountyId;
                facilitySelect.disabled = false;
                const facilityResp = await fetch(`${API_BASE}/facilities/${saved.subcountyId}`, { headers: getAuthHeaders(null) });
                if (facilityResp.ok) {
                    const facilities = await facilityResp.json();
                    facilitySelect.innerHTML = '<option value="">All Facilities</option>';
                    facilities.forEach(facility => {
                        const option = document.createElement('option');
                        option.value = facility.id;
                        option.textContent = facility.name;
                        facilitySelect.appendChild(option);
                    });
                }
            }
            if (saved.facilityId) {
                facilitySelect.value = saved.facilityId;
            }

            await loadHealthWorkers();

            if (saved.selectedHealthWorker && saved.selectedHealthWorker.id) {
                selectedHealthWorker = saved.selectedHealthWorker;
                if (healthWorkerSearch) {
                    healthWorkerSearch.value = saved.selectedHealthWorker.fullName || saved.healthWorkerSearch || '';
                }
            } else if (saved.healthWorkerSearch && healthWorkerSearch) {
                healthWorkerSearch.value = saved.healthWorkerSearch;
            }
        } finally {
            isRestoringSelection = false;
            saveHomeSelectionState();
        }
    }

    regionSelect.addEventListener('change', async function() {
        const regionId = this.value;
        if (window.EventLogger) {
            window.EventLogger.log('filter', 'region_change', {
                regionId: regionId,
                timestamp: new Date().toISOString()
            });
        }
        if (regionId) {
            districtSelect.disabled = false;
            districtSelect.innerHTML = '<option value="">Loading...</option>';
            if (subcountySelect) {
                subcountySelect.disabled = true;
                subcountySelect.innerHTML = '<option value="">All Subcounties</option>';
            }
            facilitySelect.disabled = true;
            facilitySelect.innerHTML = '<option value="">All Facilities</option>';
            clearHealthWorkerSelection();
            
            try {
                const response = await fetch(`${API_BASE}/districts/${regionId}`, { headers: getAuthHeaders(null) });
                if (!response.ok) {
                    throw new Error('Failed to load districts');
                }
                const districts = await response.json();
                districtSelect.innerHTML = '<option value="">All Districts</option>';
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
            districtSelect.innerHTML = '<option value="">All Districts</option>';
            if (subcountySelect) {
                subcountySelect.disabled = true;
                subcountySelect.innerHTML = '<option value="">All Subcounties</option>';
            }
            facilitySelect.disabled = true;
            facilitySelect.innerHTML = '<option value="">All Facilities</option>';
            clearHealthWorkerSelection();
        }
        loadHealthWorkers();
        saveHomeSelectionState();
    });

    districtSelect.addEventListener('change', async function() {
        const districtId = this.value;
        if (window.EventLogger) {
            window.EventLogger.log('filter', 'district_change', {
                districtId: districtId,
                timestamp: new Date().toISOString()
            });
        }
        if (districtId) {
            if (subcountySelect) {
                subcountySelect.disabled = false;
                subcountySelect.innerHTML = '<option value="">Loading...</option>';
            }
            facilitySelect.disabled = true;
            facilitySelect.innerHTML = '<option value="">All Facilities</option>';
            clearHealthWorkerSelection();
            
            try {
                const response = await fetch(`${API_BASE}/subcounties/${districtId}`, { headers: getAuthHeaders(null) });
                if (!response.ok) {
                    throw new Error('Failed to load subcounties');
                }
                const subcounties = await response.json();
                if (subcountySelect) {
                    subcountySelect.innerHTML = '<option value="">All Subcounties</option>';
                    subcounties.forEach(subcounty => {
                        const option = document.createElement('option');
                        option.value = subcounty.id;
                        option.textContent = subcounty.name;
                        subcountySelect.appendChild(option);
                    });
                }
            } catch (error) {
                console.error('Error loading subcounties:', error);
                if (subcountySelect) {
                    subcountySelect.innerHTML = '<option value="">Error loading subcounties</option>';
                }
            }
        } else {
            if (subcountySelect) {
                subcountySelect.disabled = true;
                subcountySelect.innerHTML = '<option value="">All Subcounties</option>';
            }
            facilitySelect.disabled = true;
            facilitySelect.innerHTML = '<option value="">All Facilities</option>';
            clearHealthWorkerSelection();
        }
        loadHealthWorkers();
        saveHomeSelectionState();
    });

    if (subcountySelect) {
        subcountySelect.addEventListener('change', async function() {
            const subcountyId = this.value;
            if (subcountyId) {
                facilitySelect.disabled = false;
                facilitySelect.innerHTML = '<option value="">Loading...</option>';
                clearHealthWorkerSelection();
                
                try {
                    const response = await fetch(`${API_BASE}/facilities/${subcountyId}`, { headers: getAuthHeaders(null) });
                    if (!response.ok) {
                        throw new Error('Failed to load facilities');
                    }
                    const facilities = await response.json();
                    facilitySelect.innerHTML = '<option value="">All Facilities</option>';
                    facilities.forEach(facility => {
                        const option = document.createElement('option');
                        option.value = facility.id;
                        option.textContent = facility.name;
                        facilitySelect.appendChild(option);
                    });
                } catch (error) {
                    console.error('Error loading facilities:', error);
                    facilitySelect.innerHTML = '<option value="">Error loading facilities</option>';
                }
            } else {
                facilitySelect.disabled = true;
                facilitySelect.innerHTML = '<option value="">All Facilities</option>';
                clearHealthWorkerSelection();
            }
            loadHealthWorkers();
            saveHomeSelectionState();
        });
    }

    facilitySelect.addEventListener('change', function() {
        loadHealthWorkers();
        saveHomeSelectionState();
    });

    // Clear health worker selection
    function clearHealthWorkerSelection() {
        allHealthWorkers = [];
        selectedHealthWorker = null;
        if (healthWorkerSearch) {
            healthWorkerSearch.value = '';
        }
        if (healthWorkerDropdown) {
            healthWorkerDropdown.innerHTML = '';
            healthWorkerDropdown.classList.remove('show');
        }
        saveHomeSelectionState();
    }

    // Filter health workers based on search input
    function filterHealthWorkers() {
        if (!healthWorkerDropdown || !healthWorkerSearch) return;
        
        const searchTerm = healthWorkerSearch.value.toLowerCase().trim();
        
        if (searchTerm.length === 0) {
            healthWorkerDropdown.classList.remove('show');
            return;
        }
        
        const filtered = allHealthWorkers.filter(hw => {
            const name = hw.fullName.toLowerCase();
            const email = hw.email ? hw.email.toLowerCase() : '';
            const phone = hw.phoneNumber ? hw.phoneNumber.toLowerCase() : '';
            const facility = hw.facilityName ? hw.facilityName.toLowerCase() : '';
            return name.includes(searchTerm) || email.includes(searchTerm) || phone.includes(searchTerm) || facility.includes(searchTerm);
        });
        
        if (filtered.length === 0) {
            healthWorkerDropdown.innerHTML = '<div class="health-worker-option text-muted p-2">No health workers found</div>';
            healthWorkerDropdown.classList.add('show');
            return;
        }
        
        healthWorkerDropdown.innerHTML = '';
        filtered.forEach(hw => {
            const option = document.createElement('div');
            option.className = 'health-worker-option';
            const displayName = hw.fullName + (hw.email ? ` (${hw.email})` : '') + (hw.phoneNumber ? ` - ${hw.phoneNumber}` : '') + ` - ${hw.facilityName}`;
            option.textContent = displayName;
            option.dataset.hwId = hw.id;
            option.dataset.facilityId = hw.facilityId;
            option.dataset.facilityName = hw.facilityName;
            option.addEventListener('click', function() {
                selectedHealthWorker = {
                    id: hw.id,
                    fullName: hw.fullName,
                    email: hw.email,
                    phoneNumber: hw.phoneNumber,
                    facilityId: hw.facilityId,
                    facilityName: hw.facilityName
                };
                healthWorkerSearch.value = hw.fullName;
                healthWorkerDropdown.classList.remove('show');
                saveHomeSelectionState();
            });
            healthWorkerDropdown.appendChild(option);
        });
        healthWorkerDropdown.classList.add('show');
    }

    // Load health workers based on filters
    async function loadHealthWorkers() {
        if (!healthWorkerSelect) return;
        
        const params = new URLSearchParams();
        const regionId = regionSelect.value;
        const districtId = districtSelect.value;
        const subcountyId = subcountySelect ? subcountySelect.value : '';
        const facilityId = facilitySelect.value;
        
        if (regionId) params.append('regionId', regionId);
        if (districtId) params.append('districtId', districtId);
        if (subcountyId) params.append('subcountyId', subcountyId);
        if (facilityId) params.append('facilityId', facilityId);
        
        try {
            const response = await fetch(`${API_BASE}/health-workers?${params.toString()}`, { headers: getAuthHeaders(null) });
            if (!response.ok) {
                throw new Error('Failed to load health workers');
            }
            const healthWorkers = await response.json();
            
            // Store health workers with all needed properties
            allHealthWorkers = healthWorkers.map(hw => ({
                id: hw.id,
                fullName: hw.fullName,
                email: hw.email || '',
                phoneNumber: hw.phoneNumber || '',
                facilityId: hw.facilityId,
                facilityName: hw.facilityName
            }));
            
            // If there's a search term, filter immediately
            if (healthWorkerSearch && healthWorkerSearch.value.trim()) {
                filterHealthWorkers();
            }

            // Re-bind selected health worker after reload if still available
            if (selectedHealthWorker && selectedHealthWorker.id) {
                const match = allHealthWorkers.find(hw => hw.id === selectedHealthWorker.id);
                if (match) {
                    selectedHealthWorker = { ...match };
                    if (healthWorkerSearch && !healthWorkerSearch.value.trim()) {
                        healthWorkerSearch.value = match.fullName;
                    }
                }
            }
        } catch (error) {
            console.error('Error loading health workers:', error);
        }
    }

    // Search input handlers
    if (healthWorkerSearch) {
        let searchTimeout;
        healthWorkerSearch.addEventListener('input', function() {
            if (selectedHealthWorker && this.value !== selectedHealthWorker.fullName) {
                selectedHealthWorker = null;
            }
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                filterHealthWorkers();
                saveHomeSelectionState();
            }, 200);
        });
        
        healthWorkerSearch.addEventListener('focus', function() {
            if (this.value.trim() && allHealthWorkers.length > 0) {
                filterHealthWorkers();
            }
        });
        
        // Close dropdown when clicking outside
        document.addEventListener('click', function(e) {
            if (healthWorkerDropdown && healthWorkerSearch && 
                !healthWorkerDropdown.contains(e.target) && 
                !healthWorkerSearch.contains(e.target)) {
                healthWorkerDropdown.classList.remove('show');
            }
        });
    }

    // Initial load of health workers
    loadHealthWorkers();
    restoreHomeSelectionState();
}

// Load assessment types
let isLoadingAssessmentTypes = false;
async function loadAssessmentTypes() {
    // Prevent concurrent calls
    if (isLoadingAssessmentTypes) {
        return;
    }
    
    const container = document.getElementById('assessmentTypesContainer');
    if (!container) return;
    
    isLoadingAssessmentTypes = true;
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
                container.innerHTML = '<div class="alert alert-warning">You do not have permission to view assessment types. Please contact your administrator.</div>';
                return;
            }
            throw new Error('Failed to load assessment types');
        }
        const types = await response.json();
        
        // Deduplicate by id to prevent showing duplicates
        const seenIds = new Set();
        const uniqueTypes = types.filter(type => {
            if (seenIds.has(type.id)) {
                return false;
            }
            seenIds.add(type.id);
            return true;
        });
        
        container.innerHTML = '';
        
        uniqueTypes.forEach(type => {
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
                // Get selected health worker from the search/selection
                if (!selectedHealthWorker) {
                    alert('Please select a health worker first');
                    if (window.EventLogger) {
                        window.EventLogger.log('assessment', 'start_failed', {
                            reason: 'no_health_worker_selected',
                            timestamp: new Date().toISOString()
                        });
                    }
                    return;
                }
                
                const healthWorkerId = selectedHealthWorker.id;
                const facilityId = selectedHealthWorker.facilityId;
                const facilityName = selectedHealthWorker.facilityName;
                const typeId = this.dataset.typeId;
                const typeCode = this.dataset.typeCode;
                localStorage.setItem(HOME_SELECTION_STORAGE_KEY, JSON.stringify({
                    regionId: document.getElementById('regionSelect')?.value || '',
                    districtId: document.getElementById('districtSelect')?.value || '',
                    subcountyId: document.getElementById('subcountySelect')?.value || '',
                    facilityId: document.getElementById('facilitySelect')?.value || '',
                    selectedHealthWorker: selectedHealthWorker,
                    healthWorkerSearch: document.getElementById('healthWorkerSearch')?.value || ''
                }));
                
                if (window.EventLogger) {
                    window.EventLogger.log('assessment', 'start', {
                        typeId: typeId,
                        typeCode: typeCode,
                        healthWorkerId: healthWorkerId,
                        facilityId: facilityId,
                        facilityName: facilityName,
                        timestamp: new Date().toISOString()
                    });
                }
                
                window.location.href = `assessment.html?typeId=${typeId}&healthWorkerId=${healthWorkerId}&facilityId=${facilityId}&facilityName=${encodeURIComponent(facilityName)}&typeCode=${typeCode}`;
            });
        });
    } catch (error) {
        console.error('Error loading assessment types:', error);
    } finally {
        isLoadingAssessmentTypes = false;
    }
}

// Load health workers with assessments for history
let isLoadingAssessments = false;
async function loadAssessments() {
    // Prevent concurrent calls
    if (isLoadingAssessments) {
        return;
    }
    
    const container = document.getElementById('assessmentsList');
    if (!container) {
        console.error('assessmentsList container not found');
        return;
    }

    isLoadingAssessments = true;
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        // Get date range filters
        const startDate = document.getElementById('startDate') ? document.getElementById('startDate').value : '';
        const endDate = document.getElementById('endDate') ? document.getElementById('endDate').value : '';

        const params = new URLSearchParams();
        if (startDate) params.append('startDate', startDate);
        if (endDate) params.append('endDate', endDate);

        const response = await fetch(`${API_BASE}/health-workers-with-assessments?${params.toString()}`, { headers });
        if (!response.ok) {
            if (response.status === 403) {
                console.error('Permission denied: You do not have permission to view health workers');
                container.innerHTML = '<div class="alert alert-warning">You do not have permission to view health workers. Please contact your administrator.</div>';
                return;
            }
            // Try to get error message from response
            let errorMsg = 'Failed to load health workers';
            try {
                const errorData = await response.json();
                errorMsg = errorData.error || errorMsg;
            } catch (e) {
                errorMsg = `Failed to load health workers (${response.status})`;
            }
            container.innerHTML = `<div class="alert alert-danger">${errorMsg}</div>`;
            return;
        }
        const healthWorkers = await response.json();
        
        // Deduplicate by id
        const seenIds = new Set();
        const uniqueHealthWorkers = healthWorkers.filter(hw => {
            if (seenIds.has(hw.id)) {
                return false;
            }
            seenIds.add(hw.id);
            return true;
        });
        
        if (uniqueHealthWorkers.length === 0) {
            container.innerHTML = '<p class="text-center text-muted">No health workers with assessments found.</p>';
            return;
        }

        container.innerHTML = '<div class="table-responsive"><table class="table table-hover"><thead><tr><th>Health Worker</th><th>Email</th><th>Phone</th><th>Facility</th><th>Region</th><th>District</th><th>Subcounty</th><th>Assessments</th><th>First Assessment</th><th>Last Assessment</th><th>Actions</th></tr></thead><tbody></tbody></table></div>';
        const tbody = container.querySelector('tbody');

        uniqueHealthWorkers.forEach(hw => {
            const row = document.createElement('tr');
            const firstDate = hw.firstAssessment ? new Date(hw.firstAssessment).toLocaleDateString() : 'N/A';
            const lastDate = hw.lastAssessment ? new Date(hw.lastAssessment).toLocaleDateString() : 'N/A';
            
            row.innerHTML = `
                <td><strong>${hw.fullName}</strong></td>
                <td>${hw.email || 'N/A'}</td>
                <td>${hw.phoneNumber || 'N/A'}</td>
                <td>${hw.facilityName || 'N/A'}</td>
                <td>${hw.regionName || 'N/A'}</td>
                <td>${hw.districtName || 'N/A'}</td>
                <td>${hw.subcountyName || 'N/A'}</td>
                <td><span class="badge bg-primary">${hw.assessmentCount}</span></td>
                <td>${firstDate}</td>
                <td>${lastDate}</td>
                <td>
                    <button class="btn btn-sm btn-primary view-performance" data-id="${hw.id}" data-name="${hw.fullName}">
                        <i class="bi bi-graph-up"></i> View Performance
                    </button>
                </td>
            `;
            tbody.appendChild(row);
        });

        // Add click handlers for view performance buttons
        container.querySelectorAll('.view-performance').forEach(btn => {
            btn.addEventListener('click', function() {
                const healthWorkerId = this.dataset.id;
                const startDate = document.getElementById('startDate') ? document.getElementById('startDate').value : '';
                const endDate = document.getElementById('endDate') ? document.getElementById('endDate').value : '';
                const params = new URLSearchParams();
                if (startDate) params.append('startDate', startDate);
                if (endDate) params.append('endDate', endDate);
                
                if (window.EventLogger) {
                    window.EventLogger.log('navigation', 'view_performance', {
                        healthWorkerId: healthWorkerId,
                        startDate: startDate,
                        endDate: endDate,
                        timestamp: new Date().toISOString()
                    });
                }
                
                window.location.href = `health-worker-performance.html?id=${healthWorkerId}&${params.toString()}`;
            });
        });
    } catch (error) {
        console.error('Error loading health workers:', error);
        const container = document.getElementById('assessmentsList');
        if (container) {
            container.innerHTML = `<div class="alert alert-danger">Error loading health workers: ${error.message || 'Unknown error'}. Please try again or contact support.</div>`;
        }
    } finally {
        isLoadingAssessments = false;
    }
}
