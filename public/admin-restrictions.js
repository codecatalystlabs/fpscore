// Admin restrictions management - applies user admin area restrictions site-wide

// Get API_BASE from window if available, otherwise use default
// Use a function to avoid const declaration conflicts with inline scripts
function getApiBase() {
    return (typeof window !== 'undefined' && window.API_BASE) || '/api';
}

// Use getApiBase() function instead of const to avoid conflicts
let userRestrictions = null;
let restrictionsApplied = false;

// Expose restrictionsApplied globally so it can be reset by other scripts
if (typeof window !== 'undefined') {
    Object.defineProperty(window, 'restrictionsApplied', {
        get: () => restrictionsApplied,
        set: (value) => { restrictionsApplied = value; }
    });
}

// Load user restrictions on page load
async function loadUserRestrictions() {
    try {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        const response = await fetch(`${getApiBase()}/user/admin-areas`, { headers });
        if (!response.ok) {
            console.error('Failed to load user restrictions:', response.status, response.statusText);
            return null;
        }
        userRestrictions = await response.json();
        console.log('Loaded user restrictions:', userRestrictions);
        return userRestrictions;
    } catch (error) {
        console.error('Error loading user restrictions:', error);
        return null;
    }
}

// Apply restrictions to geographic dropdowns
async function applyAdminRestrictions() {
    // Prevent applying restrictions multiple times in quick succession
    if (restrictionsApplied) {
        return;
    }

    if (!userRestrictions) {
        userRestrictions = await loadUserRestrictions();
    }
    
    // Check if user is admin - if so, do NOT apply any restrictions
    if (!userRestrictions) {
        console.warn('Could not load user restrictions');
        restrictionsApplied = true;
        return;
    }
    
    // Check admin status - use strict equality and also check for truthy values
    const isAdmin = userRestrictions.isAdmin === true || userRestrictions.isAdmin === 1 || userRestrictions.isAdmin === "true";
    
    if (isAdmin) {
        // Admin users - no restrictions, ensure dropdowns are enabled
        console.log('User is admin - no restrictions applied. Admin status:', userRestrictions.isAdmin);
        restrictionsApplied = true;
        
        // Make sure all dropdowns are enabled for admin users
        const regionSelect = document.getElementById('regionSelect') || document.getElementById('filterRegion');
        const districtSelect = document.getElementById('districtSelect') || document.getElementById('filterDistrict');
        const subcountySelect = document.getElementById('subcountySelect') || document.getElementById('filterSubcounty');
        const facilitySelect = document.getElementById('facilitySelect') || document.getElementById('filterFacility');
        
        if (regionSelect) {
            regionSelect.disabled = false;
            regionSelect.classList.remove('bg-light');
        }
        if (districtSelect) {
            districtSelect.disabled = false;
            districtSelect.classList.remove('bg-light');
        }
        if (subcountySelect) {
            subcountySelect.disabled = false;
            subcountySelect.classList.remove('bg-light');
        }
        if (facilitySelect) {
            facilitySelect.disabled = false;
            facilitySelect.classList.remove('bg-light');
        }
        
        return;
    }
    
    console.log('User is NOT admin. Admin status:', userRestrictions.isAdmin, 'Type:', typeof userRestrictions.isAdmin);

    const restrictions = userRestrictions.restrictions;
    const level = restrictions.restrictionLevel;

    if (!level) {
        // No restrictions
        restrictionsApplied = true;
        return;
    }

    // Apply restrictions based on level
    if (level === 'facility') {
        await setRestrictedHierarchy(
            restrictions.regionId,
            restrictions.regionName,
            restrictions.districtId,
            restrictions.districtName,
            restrictions.subcountyId,
            restrictions.subcountyName,
            restrictions.facilityId,
            restrictions.facilityName,
            true, // regionDisabled
            true, // districtDisabled
            true, // subcountyDisabled
            true  // facilityDisabled
        );
    } else if (level === 'subcounty') {
        await setRestrictedHierarchy(
            restrictions.regionId,
            restrictions.regionName,
            restrictions.districtId,
            restrictions.districtName,
            restrictions.subcountyId,
            restrictions.subcountyName,
            null,
            null,
            true, // regionDisabled
            true, // districtDisabled
            true, // subcountyDisabled
            false // facilityEnabled
        );
    } else if (level === 'district') {
        await setRestrictedHierarchy(
            restrictions.regionId,
            restrictions.regionName,
            restrictions.districtId,
            restrictions.districtName,
            null,
            null,
            null,
            null,
            true, // regionDisabled
            true, // districtDisabled
            false, // subcountyEnabled
            false  // facilityEnabled
        );
    } else if (level === 'region') {
        await setRestrictedHierarchy(
            restrictions.regionId,
            restrictions.regionName,
            null,
            null,
            null,
            null,
            null,
            null,
            true, // regionDisabled
            false, // districtEnabled
            false, // subcountyEnabled
            false  // facilityEnabled
        );
    }
    
    restrictionsApplied = true;
}

// Set restricted hierarchy and disable appropriate dropdowns
async function setRestrictedHierarchy(regionId, regionName, districtId, districtName, subcountyId, subcountyName, facilityId, facilityName, regionDisabled, districtDisabled, subcountyDisabled, facilityDisabled) {
    // Try multiple possible element IDs for different pages
    const regionSelect = document.getElementById('regionSelect') || document.getElementById('filterRegion');
    const districtSelect = document.getElementById('districtSelect') || document.getElementById('filterDistrict');
    const subcountySelect = document.getElementById('subcountySelect') || document.getElementById('filterSubcounty');
    const facilitySelect = document.getElementById('facilitySelect') || document.getElementById('filterFacility');

    if (!regionSelect) {
        // Dropdowns not found on this page
        console.log('No region dropdown found on this page');
        return;
    }
    
    console.log('setRestrictedHierarchy called:', {
        regionId, districtId, subcountyId, facilityId,
        regionDisabled, districtDisabled, subcountyDisabled, facilityDisabled
    });

    // Helper function to get auth headers
    function getAuthHeaders() {
        const token = localStorage.getItem('token');
        const headers = {};
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        return headers;
    }

    // Set region - ALWAYS set this first, even if dropdown is empty
    if (regionId && regionName) {
        // Check if the option already exists
        const existingOption = Array.from(regionSelect.options).find(opt => opt.value === String(regionId));
        if (!existingOption) {
            // Option doesn't exist, add it
            const opt = document.createElement('option');
            opt.value = regionId;
            opt.textContent = regionName;
            regionSelect.appendChild(opt);
        }
        // Set the value
        regionSelect.value = regionId;
        
        if (regionDisabled) {
            regionSelect.disabled = true;
            regionSelect.classList.add('bg-light');
        } else {
            regionSelect.disabled = false;
            regionSelect.classList.remove('bg-light');
        }
        
        // Load districts for the region
        if (districtId) {
            try {
                const response = await fetch(`${getApiBase()}/districts/${regionId}`, { headers: getAuthHeaders() });
                if (response.ok) {
                    const districts = await response.json();
                    if (districtSelect) {
                        // Clear only if not already restricted
                        const isDistrictRestricted = districtSelect.disabled && districtSelect.classList.contains('bg-light');
                        if (!isDistrictRestricted) {
                            districtSelect.innerHTML = '';
                        }
                        districts.forEach(d => {
                            // Check if option already exists
                            const existingDistOption = Array.from(districtSelect.options).find(opt => opt.value === String(d.id));
                            if (!existingDistOption) {
                                const opt = document.createElement('option');
                                opt.value = d.id;
                                opt.textContent = d.name;
                                if (d.id === districtId) {
                                    opt.selected = true;
                                }
                                districtSelect.appendChild(opt);
                            } else if (d.id === districtId) {
                                existingDistOption.selected = true;
                            }
                        });
                        districtSelect.value = districtId;
                        
                        if (districtDisabled) {
                            districtSelect.disabled = true;
                            districtSelect.classList.add('bg-light');
                        } else {
                            districtSelect.disabled = false;
                            districtSelect.classList.remove('bg-light');
                        }
                        
                        // Load subcounties for the district (if subcountyId is provided OR if subcounty is enabled)
                        if (subcountyId || (!subcountyDisabled && subcountySelect)) {
                            console.log('Loading subcounties for district:', districtId, 'subcountyDisabled:', subcountyDisabled, 'subcountyId:', subcountyId);
                            const subcountyResponse = await fetch(`${getApiBase()}/subcounties/${districtId}`, { headers: getAuthHeaders() });
                            if (subcountyResponse.ok) {
                                const subcounties = await subcountyResponse.json();
                                console.log('Loaded subcounties:', subcounties.length);
                                if (subcountySelect) {
                                    // Clear only if not already restricted
                                    const isSubcountyRestricted = subcountySelect.disabled && subcountySelect.classList.contains('bg-light');
                                    if (!isSubcountyRestricted) {
                                        subcountySelect.innerHTML = subcountyId ? '' : '<option value="">Select Subcounty</option>';
                                    }
                                    subcounties.forEach(s => {
                                        // Check if option already exists
                                        const existingSubOption = Array.from(subcountySelect.options).find(opt => opt.value === String(s.id));
                                        if (!existingSubOption) {
                                            const opt = document.createElement('option');
                                            opt.value = s.id;
                                            opt.textContent = s.name;
                                            if (s.id === subcountyId) {
                                                opt.selected = true;
                                            }
                                            subcountySelect.appendChild(opt);
                                        } else if (s.id === subcountyId) {
                                            existingSubOption.selected = true;
                                        }
                                    });
                                    if (subcountyId) {
                                        subcountySelect.value = subcountyId;
                                    }
                                    
                                    if (subcountyDisabled) {
                                        subcountySelect.disabled = true;
                                        subcountySelect.classList.add('bg-light');
                                    } else {
                                        subcountySelect.disabled = false;
                                        subcountySelect.classList.remove('bg-light');
                                        console.log('Subcounty dropdown enabled');
                                    }
                                    
                                    // Load facilities for the subcounty (if subcountyId is set AND (facilityId is provided OR facility is enabled))
                                    if (subcountyId && (facilityId || (!facilityDisabled && facilitySelect))) {
                                        console.log('Loading facilities for subcounty:', subcountyId);
                                        const facilityResponse = await fetch(`${getApiBase()}/facilities/${subcountyId}`, { headers: getAuthHeaders() });
                                        if (facilityResponse.ok) {
                                            const facilities = await facilityResponse.json();
                                            console.log('Loaded facilities:', facilities.length);
                                            if (facilitySelect) {
                                                // Clear only if not already restricted
                                                const isFacilityRestricted = facilitySelect.disabled && facilitySelect.classList.contains('bg-light');
                                                if (!isFacilityRestricted) {
                                                    facilitySelect.innerHTML = facilityId ? '' : '<option value="">Select Facility</option>';
                                                }
                                                facilities.forEach(f => {
                                                    // Check if option already exists
                                                    const existingFacOption = Array.from(facilitySelect.options).find(opt => opt.value === String(f.id));
                                                    if (!existingFacOption) {
                                                        const opt = document.createElement('option');
                                                        opt.value = f.id;
                                                        opt.textContent = f.name;
                                                        if (f.id === facilityId) {
                                                            opt.selected = true;
                                                        }
                                                        facilitySelect.appendChild(opt);
                                                    } else if (f.id === facilityId) {
                                                        existingFacOption.selected = true;
                                                    }
                                                });
                                                if (facilityId) {
                                                    facilitySelect.value = facilityId;
                                                }
                                                
                                                if (facilityDisabled) {
                                                    facilitySelect.disabled = true;
                                                    facilitySelect.classList.add('bg-light');
                                                } else {
                                                    facilitySelect.disabled = false;
                                                    facilitySelect.classList.remove('bg-light');
                                                    console.log('Facility dropdown enabled');
                                                }
                                            }
                                        }
                                    }
                                }
                            } else {
                                console.error('Failed to load subcounties:', subcountyResponse.status, subcountyResponse.statusText);
                            }
                        } else {
                            console.log('Skipping subcounty load - subcountyId:', subcountyId, 'subcountyDisabled:', subcountyDisabled, 'subcountySelect exists:', !!subcountySelect);
                        }
                    }
                }
            } catch (e) {
                console.error('Error loading districts:', e);
            }
        } else if (districtSelect && !districtDisabled) {
            // Load all districts for the region if district is enabled
            try {
                const response = await fetch(`${getApiBase()}/districts/${regionId}`, { headers: getAuthHeaders() });
                if (response.ok) {
                    const districts = await response.json();
                    districtSelect.innerHTML = '<option value="">All Districts</option>';
                    districts.forEach(d => {
                        const opt = document.createElement('option');
                        opt.value = d.id;
                        opt.textContent = d.name;
                        districtSelect.appendChild(opt);
                    });
                    districtSelect.disabled = false;
                    districtSelect.classList.remove('bg-light');
                }
            } catch (e) {
                console.error('Error loading districts:', e);
            }
        } else if (districtSelect && districtDisabled) {
            // Disable district if it should be disabled
            districtSelect.disabled = true;
            districtSelect.classList.add('bg-light');
        }
    }

    // Disable dropdowns that should be disabled but don't have values set
    if (districtSelect && districtDisabled && !districtId) {
        districtSelect.disabled = true;
        districtSelect.classList.add('bg-light');
    }
    if (subcountySelect && subcountyDisabled && !subcountyId) {
        subcountySelect.disabled = true;
        subcountySelect.classList.add('bg-light');
    }
    if (facilitySelect && facilityDisabled && !facilityId) {
        facilitySelect.disabled = true;
        facilitySelect.classList.add('bg-light');
    }
}

// Apply restrictions immediately when page loads
async function initializeRestrictions() {
    // Load restrictions first
    if (!userRestrictions) {
        userRestrictions = await loadUserRestrictions();
    }
    
    // Check if user is admin - if so, don't apply any restrictions
    // Check for various possible admin status values
    const isAdmin = userRestrictions && (
        userRestrictions.isAdmin === true || 
        userRestrictions.isAdmin === 1 || 
        userRestrictions.isAdmin === "true" ||
        userRestrictions.isAdmin === "1"
    );
    
    if (isAdmin) {
        console.log('User is admin - skipping restriction initialization. Admin status:', userRestrictions.isAdmin);
        restrictionsApplied = true;
        
        // Ensure dropdowns are enabled for admin users
        const regionSelect = document.getElementById('regionSelect') || document.getElementById('filterRegion');
        const districtSelect = document.getElementById('districtSelect') || document.getElementById('filterDistrict');
        const subcountySelect = document.getElementById('subcountySelect') || document.getElementById('filterSubcounty');
        const facilitySelect = document.getElementById('facilitySelect') || document.getElementById('filterFacility');
        
        if (regionSelect) {
            regionSelect.disabled = false;
            regionSelect.classList.remove('bg-light');
        }
        if (districtSelect) {
            districtSelect.disabled = false;
            districtSelect.classList.remove('bg-light');
        }
        if (subcountySelect) {
            subcountySelect.disabled = false;
            subcountySelect.classList.remove('bg-light');
        }
        if (facilitySelect) {
            facilitySelect.disabled = false;
            facilitySelect.classList.remove('bg-light');
        }
        
        return; // Admin users - no restrictions, exit early
    }
    
    if (userRestrictions) {
        console.log('User is NOT admin. Admin status:', userRestrictions.isAdmin, 'Type:', typeof userRestrictions.isAdmin);
    }
    
    // If user has restrictions, wait for dropdowns to load first, then apply restrictions
    if (userRestrictions && !userRestrictions.isAdmin && userRestrictions.restrictions && userRestrictions.restrictions.restrictionLevel) {
        console.log('User has restrictions. Restriction level:', userRestrictions.restrictions.restrictionLevel);
        
        // Wait for regions/filters to load, then apply restrictions
        // Check if we're on a page with loadRegions or loadFilters
        const hasLoadRegions = typeof window.loadRegions === 'function';
        const hasLoadFilters = typeof window.loadFilters === 'function';
        
        if (hasLoadRegions || hasLoadFilters) {
            // Wait for the load functions to complete, then apply restrictions
            // The hooks will handle reapplying restrictions after loading
            setTimeout(async () => {
                await applyAdminRestrictions();
            }, 1000);
        } else {
            // No load functions, apply restrictions immediately
            await applyAdminRestrictions();
        }
        
        // Set up a watcher to reapply restrictions if dropdowns get cleared
        const checkInterval = setInterval(() => {
            // Double-check user is still not admin
            const isAdmin = userRestrictions && (
                userRestrictions.isAdmin === true || 
                userRestrictions.isAdmin === 1 || 
                userRestrictions.isAdmin === "true"
            );
            if (isAdmin) {
                clearInterval(checkInterval);
                return;
            }
            
            const regionSelect = document.getElementById('regionSelect') || document.getElementById('filterRegion');
            if (regionSelect && userRestrictions && userRestrictions.restrictions) {
                const restrictions = userRestrictions.restrictions;
                const expectedRegionId = String(restrictions.regionId);
                
                // If region dropdown doesn't have the restricted value, reapply
                if (restrictions.regionId && regionSelect.value !== expectedRegionId) {
                    const hasRestrictedOption = Array.from(regionSelect.options).some(
                        opt => opt.value === expectedRegionId
                    );
                    
                    if (!hasRestrictedOption || regionSelect.value !== expectedRegionId) {
                        // Reapply restrictions
                        restrictionsApplied = false;
                        applyAdminRestrictions();
                    }
                }
            }
        }, 500);
        
        // Stop checking after 15 seconds
        setTimeout(() => clearInterval(checkInterval), 15000);
    } else {
        restrictionsApplied = true;
    }
}

// Initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeRestrictions);
} else {
    initializeRestrictions();
}

// Also expose applyAdminRestrictions globally so it can be called from other scripts
if (typeof window !== 'undefined') {
    window.applyAdminRestrictions = applyAdminRestrictions;
    
    // Hook into loadRegions if it exists (or wait for it to be defined)
    function hookLoadRegions() {
        if (typeof window.loadRegions === 'function' && !window.loadRegions._hooked) {
            const originalLoadRegions = window.loadRegions;
            window.loadRegions = async function() {
                await originalLoadRegions();
                // Wait a bit for DOM to update
                await new Promise(resolve => setTimeout(resolve, 200));
                
                // Reapply restrictions after regions are loaded (only for non-admin users)
                if (!userRestrictions) {
                    userRestrictions = await loadUserRestrictions();
                }
                
                const isAdmin = userRestrictions && (
                    userRestrictions.isAdmin === true || 
                    userRestrictions.isAdmin === 1 || 
                    userRestrictions.isAdmin === "true"
                );
                
                if (!isAdmin && userRestrictions && userRestrictions.restrictions && userRestrictions.restrictions.restrictionLevel) {
                    console.log('Reapplying restrictions after loadRegions()');
                    restrictionsApplied = false;
                    await applyAdminRestrictions();
                } else if (isAdmin) {
                    // Ensure dropdowns are enabled for admin users
                    const regionSelect = document.getElementById('regionSelect') || document.getElementById('filterRegion');
                    if (regionSelect) {
                        regionSelect.disabled = false;
                        regionSelect.classList.remove('bg-light');
                    }
                }
            };
            window.loadRegions._hooked = true;
        } else if (typeof window.loadRegions !== 'function') {
            // Wait a bit and try again
            setTimeout(hookLoadRegions, 100);
        }
    }
    
    // Hook into loadFilters if it exists (for dashboard/reports pages)
    function hookLoadFilters() {
        if (typeof window.loadFilters === 'function' && !window.loadFilters._hooked) {
            const originalLoadFilters = window.loadFilters;
            window.loadFilters = async function() {
                await originalLoadFilters();
                // Wait a bit for DOM to update
                await new Promise(resolve => setTimeout(resolve, 200));
                
                // Reapply restrictions after filters are loaded (only for non-admin users)
                if (!userRestrictions) {
                    userRestrictions = await loadUserRestrictions();
                }
                
                const isAdmin = userRestrictions && (
                    userRestrictions.isAdmin === true || 
                    userRestrictions.isAdmin === 1 || 
                    userRestrictions.isAdmin === "true"
                );
                
                if (!isAdmin && userRestrictions && userRestrictions.restrictions && userRestrictions.restrictions.restrictionLevel) {
                    console.log('Reapplying restrictions after loadFilters()');
                    restrictionsApplied = false;
                    await applyAdminRestrictions();
                } else if (isAdmin) {
                    // Ensure dropdowns are enabled for admin users
                    const regionSelect = document.getElementById('regionSelect') || document.getElementById('filterRegion');
                    const districtSelect = document.getElementById('districtSelect') || document.getElementById('filterDistrict');
                    if (regionSelect) {
                        regionSelect.disabled = false;
                        regionSelect.classList.remove('bg-light');
                    }
                    if (districtSelect) {
                        districtSelect.disabled = false;
                        districtSelect.classList.remove('bg-light');
                    }
                }
            };
            window.loadFilters._hooked = true;
        } else if (typeof window.loadFilters !== 'function') {
            // Wait a bit and try again (loadFilters might be defined inline in script tags)
            setTimeout(hookLoadFilters, 100);
        }
    }
    
    // Try to hook immediately, or wait for functions to be defined
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => {
            hookLoadRegions();
            hookLoadFilters();
        });
    } else {
        hookLoadRegions();
        hookLoadFilters();
    }
    
    // Also try to hook after a delay to catch inline functions
    setTimeout(() => {
        hookLoadRegions();
        hookLoadFilters();
    }, 1000);
}
