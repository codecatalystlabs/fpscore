// Centralized Event Logging System
// Logs all events in the application for debugging and analytics

const EventLogger = {
    enabled: true,
    logLevel: 'all', // 'all', 'error', 'warning', 'info', 'debug'
    events: [],
    maxEvents: 1000, // Keep last 1000 events in memory
    pendingEvents: [], // Events waiting to be sent to server
    batchSize: 10, // Send events in batches of 10
    batchInterval: 5000, // Send every 5 seconds
    isSending: false, // Flag to prevent concurrent sends
    originalFetch: null, // Store original fetch to avoid recursion
    
    // Initialize event logging
    init() {
        if (!this.enabled) return;
        
        // Log page load
        this.log('page', 'load', {
            url: window.location.href,
            referrer: document.referrer,
            userAgent: navigator.userAgent,
            timestamp: new Date().toISOString()
        });
        
        // Set up global event listeners
        this.setupGlobalListeners();
        
        // NOTE: Frontend API logging disabled - now handled by backend middleware
        // this.interceptFetch();
        
        // Start batch sending interval
        this.startBatchSending();
        
        // Log when page is about to unload
        window.addEventListener('beforeunload', () => {
            this.log('page', 'unload', {
                url: window.location.href,
                timestamp: new Date().toISOString()
            });
            // Send any pending events before unload
            this.flushPendingEvents();
        });
    },
    
    // Start batch sending of events
    startBatchSending() {
        setInterval(() => {
            this.flushPendingEvents();
        }, this.batchInterval);
    },
    
    // Flush pending events to server
    async flushPendingEvents() {
        if (this.isSending || this.pendingEvents.length === 0) {
            return;
        }
        
        this.isSending = true;
        
        try {
            const eventsToSend = this.pendingEvents.splice(0, this.batchSize);
            
            const token = localStorage.getItem('token');
            const headers = {
                'Content-Type': 'application/json'
            };
            
            // Add auth header if token exists
            if (token) {
                headers['Authorization'] = `Bearer ${token}`;
            }
            
            // Use original fetch to bypass our interceptor and prevent infinite loop
            const fetchFunc = this.originalFetch || fetch;
            await fetchFunc('/api/events/log', {
                method: 'POST',
                headers: headers,
                body: JSON.stringify({ events: eventsToSend })
            }).catch(error => {
                // Silently fail - don't spam console
                // If send fails, don't retry to avoid building up queue
            });
        } catch (error) {
            // Silently fail
        } finally {
            this.isSending = false;
        }
    },
    
    // Setup global event listeners for common events
    setupGlobalListeners() {
        // Log all clicks (with delegation to avoid performance issues)
        document.addEventListener('click', (e) => {
            const target = e.target;
            const tagName = target.tagName.toLowerCase();
            const id = target.id || null;
            const className = target.className || null;
            const text = target.textContent?.trim().substring(0, 100) || null;
            
            // Skip logging for very frequent events
            if (target.closest('.spinner-border') || target.closest('.progress-bar')) {
                return;
            }
            
            this.log('interaction', 'click', {
                tag: tagName,
                id: id,
                className: className,
                text: text,
                href: target.href || null,
                dataset: Object.assign({}, target.dataset),
                timestamp: new Date().toISOString()
            });
        }, true);
        
        // Log form submissions
        document.addEventListener('submit', (e) => {
            const form = e.target;
            const formData = {};
            
            // Collect form data (excluding passwords)
            Array.from(form.elements).forEach(element => {
                if (element.name && element.type !== 'password') {
                    formData[element.name] = element.value;
                }
            });
            
            this.log('form', 'submit', {
                formId: form.id || null,
                formAction: form.action || null,
                formMethod: form.method || 'get',
                formData: formData,
                timestamp: new Date().toISOString()
            });
        });
        
        // Log input changes (debounced)
        let changeTimeout;
        document.addEventListener('input', (e) => {
            clearTimeout(changeTimeout);
            changeTimeout = setTimeout(() => {
                const target = e.target;
                if (target.type === 'password') return; // Don't log passwords
                
                this.log('form', 'input', {
                    elementId: target.id || null,
                    elementName: target.name || null,
                    elementType: target.type || null,
                    value: target.value?.substring(0, 200) || null, // Limit length
                    timestamp: new Date().toISOString()
                });
            }, 500); // Debounce by 500ms
        });
        
        // Log select/change events
        document.addEventListener('change', (e) => {
            const target = e.target;
            if (target.tagName.toLowerCase() === 'select' || target.type === 'checkbox' || target.type === 'radio') {
                this.log('form', 'change', {
                    elementId: target.id || null,
                    elementName: target.name || null,
                    elementType: target.type || null,
                    value: target.value || (target.checked ? 'checked' : 'unchecked'),
                    timestamp: new Date().toISOString()
                });
            }
        });
        
        // Log navigation (check less frequently to reduce overhead)
        let lastUrl = window.location.href;
        setInterval(() => {
            if (window.location.href !== lastUrl) {
                this.log('navigation', 'url_change', {
                    from: lastUrl,
                    to: window.location.href,
                    timestamp: new Date().toISOString()
                });
                lastUrl = window.location.href;
            }
        }, 3000); // Check every 3 seconds instead of 1
    },
    
    // Intercept fetch calls to log API requests/responses
    interceptFetch() {
        // Save original fetch for sending events later
        this.originalFetch = window.fetch;
        const originalFetch = this.originalFetch;
        const self = this;
        
        window.fetch = async function(...args) {
            const url = args[0];
            const options = args[1] || {};
            const method = options.method || 'GET';
            const startTime = Date.now();
            
            // CRITICAL: Skip logging for ALL event-related endpoints to prevent infinite loop
            if (url && url.includes('/api/events')) {
                return originalFetch.apply(this, args);
            }
            
            // Log request
            self.log('api', 'request', {
                url: url,
                method: method,
                headers: options.headers ? Object.keys(options.headers) : [],
                hasBody: !!options.body,
                bodySize: options.body ? String(options.body).length : 0,
                timestamp: new Date().toISOString()
            });
            
            try {
                const response = await originalFetch.apply(this, args);
                const endTime = Date.now();
                const duration = endTime - startTime;
                
                // Clone response to read body without consuming it
                const responseClone = response.clone();
                let responseData = null;
                
                try {
                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        responseData = await responseClone.json();
                    } else {
                        responseData = await responseClone.text();
                        if (responseData.length > 500) {
                            responseData = responseData.substring(0, 500) + '... (truncated)';
                        }
                    }
                } catch (e) {
                    // Couldn't parse response, that's okay
                }
                
                // Log response
                self.log('api', 'response', {
                    url: url,
                    method: method,
                    status: response.status,
                    statusText: response.statusText,
                    duration: duration,
                    success: response.ok,
                    dataSize: responseData ? JSON.stringify(responseData).length : 0,
                    timestamp: new Date().toISOString()
                });
                
                return response;
            } catch (error) {
                const endTime = Date.now();
                const duration = endTime - startTime;
                
                // Log error
                self.log('api', 'error', {
                    url: url,
                    method: method,
                    error: error.message,
                    duration: duration,
                    timestamp: new Date().toISOString()
                });
                
                throw error;
            }
        };
    },
    
    // Main logging function
    log(category, event, data = {}) {
        if (!this.enabled) return;
        
        const logEntry = {
            category: category, // 'page', 'interaction', 'form', 'api', 'assessment', 'navigation', 'error'
            event: event,
            data: data,
            timestamp: new Date().toISOString(),
            page: window.location.pathname,
            sessionId: this.getSessionId()
        };
        
        // Add to events array
        this.events.push(logEntry);
        
        // Keep only last maxEvents
        if (this.events.length > this.maxEvents) {
            this.events.shift();
        }
        
        // Output to console with appropriate level
        this.consoleLog(logEntry);
        
        // Send to server for database storage
        this.sendToServer(logEntry);
    },
    
    // Console logging with levels
    consoleLog(logEntry) {
        const { category, event, data, timestamp } = logEntry;
        const message = `[${timestamp}] ${category.toUpperCase()}.${event}`;
        
        // Format data for console
        const consoleData = { ...data };
        
        switch (category) {
            case 'api':
                if (event === 'error') {
                    console.error(message, consoleData);
                } else if (event === 'response' && !data.success) {
                    console.warn(message, consoleData);
                } else {
                    console.info(message, consoleData);
                }
                break;
            case 'error':
                console.error(message, consoleData);
                break;
            case 'form':
            case 'interaction':
                console.debug(message, consoleData);
                break;
            default:
                console.log(message, consoleData);
        }
    },
    
    // Get or create session ID
    getSessionId() {
        let sessionId = sessionStorage.getItem('eventLoggerSessionId');
        if (!sessionId) {
            sessionId = 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
            sessionStorage.setItem('eventLoggerSessionId', sessionId);
        }
        return sessionId;
    },
    
    // Send log entry to server (batched)
    async sendToServer(logEntry) {
        try {
            // Skip sending certain high-frequency events to server to reduce load
            const skipCategories = ['interaction', 'form'];
            const skipEvents = ['input', 'click'];
            
            // Only send important events or API events to the server
            if (skipCategories.includes(logEntry.category) && skipEvents.includes(logEntry.event)) {
                // Keep in memory but don't send to server
                return;
            }
            
            // Add to pending queue for batch sending
            this.pendingEvents.push(logEntry);
            
            // If queue is getting large, flush immediately
            if (this.pendingEvents.length >= this.batchSize) {
                this.flushPendingEvents();
            }
        } catch (error) {
            // Silently fail
        }
    },
    
    // Get all events
    getAllEvents() {
        return this.events;
    },
    
    // Get events by category
    getEventsByCategory(category) {
        return this.events.filter(e => e.category === category);
    },
    
    // Get events by event type
    getEventsByType(event) {
        return this.events.filter(e => e.event === event);
    },
    
    // Clear events
    clearEvents() {
        this.events = [];
    },
    
    // Export events as JSON
    exportEvents() {
        return JSON.stringify(this.events, null, 2);
    },
    
    // Download events as file
    downloadEvents() {
        const data = this.exportEvents();
        const blob = new Blob([data], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `events_${new Date().toISOString().replace(/:/g, '-')}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }
};

// Initialize when script loads
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => EventLogger.init());
} else {
    EventLogger.init();
}

// Expose globally for debugging
window.EventLogger = EventLogger;

