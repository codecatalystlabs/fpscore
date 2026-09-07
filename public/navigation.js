// Navigation — tool-scoped sidebar (set data-nav-tool on <body>)
async function loadNavigation() {
    try {
        const token = localStorage.getItem('token');
        if (!token) {
            window.location.href = '/';
            return;
        }

        const navContainer = document.getElementById('navigationItems');
        if (!navContainer) return;

        const tool = document.body.getAttribute('data-nav-tool') ||
            sessionStorage.getItem('active_tool') || '';
        if (!tool) {
            navContainer.innerHTML = '';
            return;
        }
        sessionStorage.setItem('active_tool', tool);

        const response = await fetch('/api/navigation?tool=' + encodeURIComponent(tool), {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) {
            console.error('Failed to load navigation');
            return;
        }

        const items = await response.json();
        const currentPage = window.location.pathname.split('/').pop() || 'tools.html';
        navContainer.innerHTML = '';
        const useMoh = document.body.classList.contains('moh-body') ||
            !!navContainer.closest('.moh-sidebar');

        items.forEach(item => {
            const link = document.createElement('a');
            link.href = item.url;
            link.className = useMoh ? 'moh-nav-item' : 'list-group-item list-group-item-action';
            if (item.url === currentPage ||
                (currentPage === 'assessment.html' && item.url === 'home.html') ||
                (currentPage === 'view-assessment.html' && item.url === 'home.html') ||
                (currentPage === 'health-worker-performance.html' && item.url === 'health-workers.html') ||
                (currentPage === 'rhspars-assessment.html' && item.url === 'rhspars-home.html') ||
                (currentPage === 'rhspars-view.html' && item.url === 'rhspars-history.html')) {
                link.classList.add('active');
            }
            if (useMoh) {
                link.innerHTML = `<i class="bi ${item.icon}"></i><span>${item.label}</span>`;
            } else {
                link.innerHTML = `<i class="bi ${item.icon} me-2"></i> ${item.label}`;
            }
            navContainer.appendChild(link);
        });

        const logoutLink = document.createElement('a');
        logoutLink.href = '#';
        logoutLink.className = useMoh ? 'moh-nav-item' : 'list-group-item list-group-item-action';
        logoutLink.onclick = function(e) {
            e.preventDefault();
            logout();
        };
        logoutLink.innerHTML = useMoh
            ? '<i class="bi bi-box-arrow-right"></i><span>Logout</span>'
            : '<i class="bi bi-box-arrow-right me-2"></i> Logout';
        navContainer.appendChild(logoutLink);
    } catch (error) {
        console.error('Error loading navigation:', error);
    }
}

function logout() {
    if (window.EventLogger) {
        window.EventLogger.log('auth', 'logout', { timestamp: new Date().toISOString() });
    }
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    sessionStorage.removeItem('rhspars_draft');
    window.location.href = '/';
}

function authHeaders() {
    return {
        'Authorization': 'Bearer ' + localStorage.getItem('token'),
        'Content-Type': 'application/json'
    };
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', loadNavigation);
} else {
    loadNavigation();
}
