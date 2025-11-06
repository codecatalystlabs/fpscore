// Navigation management - loads navigation items based on user permissions
async function loadNavigation() {
    try {
        const token = localStorage.getItem('token');
        if (!token) {
            window.location.href = '/';
            return;
        }

        const response = await fetch('/api/navigation', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            console.error('Failed to load navigation');
            return;
        }

        const items = await response.json();
        const navContainer = document.getElementById('navigationItems');
        if (!navContainer) return;

        navContainer.innerHTML = '';
        items.forEach(item => {
            const link = document.createElement('a');
            link.href = item.url;
            link.className = 'list-group-item list-group-item-action';
            
            // Check if current page matches
            const currentPage = window.location.pathname.split('/').pop() || 'home.html';
            if (item.url === currentPage || (currentPage === '' && item.url === 'home.html')) {
                link.classList.add('active');
            }

            link.innerHTML = `<i class="bi ${item.icon} me-2"></i> ${item.label}`;
            navContainer.appendChild(link);
        });

        // Add logout link
        const logoutLink = document.createElement('a');
        logoutLink.href = '#';
        logoutLink.className = 'list-group-item list-group-item-action';
        logoutLink.onclick = function(e) {
            e.preventDefault();
            logout();
        };
        logoutLink.innerHTML = '<i class="bi bi-box-arrow-right me-2"></i> Logout';
        navContainer.appendChild(logoutLink);

    } catch (error) {
        console.error('Error loading navigation:', error);
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/';
}

// Load navigation when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', loadNavigation);
} else {
    loadNavigation();
}

