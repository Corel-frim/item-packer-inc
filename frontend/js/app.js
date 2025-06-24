/**
 * Main application module
 */
(function() {
    // DOM elements
    const tabs = document.querySelectorAll('.tab');
    const tabContents = document.querySelectorAll('.tab-content');

    /**
     * Initialize the application
     */
    function init() {
        // Set up tab switching
        tabs.forEach(tab => {
            tab.addEventListener('click', () => switchTab(tab.dataset.tab));
        });

        // Initialize modules
        packsModule.init();
        ordersModule.init();
    }

    /**
     * Switch between tabs
     * @param {string} tabId - The ID of the tab to switch to
     */
    function switchTab(tabId) {
        // Update active tab
        tabs.forEach(tab => {
            if (tab.dataset.tab === tabId) {
                tab.classList.add('active');
            } else {
                tab.classList.remove('active');
            }
        });

        // Update active content
        tabContents.forEach(content => {
            if (content.id === `${tabId}-section`) {
                content.classList.add('active');
            } else {
                content.classList.remove('active');
            }
        });

        // Refresh data when switching tabs
        if (tabId === 'packs') {
            packsModule.loadPacks();
        } else if (tabId === 'orders') {
            ordersModule.loadOrders();
        }
    }

    // Initialize the application when the DOM is loaded
    document.addEventListener('DOMContentLoaded', init);
})();