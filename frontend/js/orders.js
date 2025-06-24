/**
 * Orders module for managing order-related UI functionality
 */
const ordersModule = (function() {
    // DOM elements
    const createOrderForm = document.getElementById('create-order-form');
    const orderAmountInput = document.getElementById('order-amount');
    const orderResult = document.getElementById('order-result');
    const requestedItemsEl = document.getElementById('requested-items');
    const overpackedItemsEl = document.getElementById('overpacked-items');
    const totalItemsEl = document.getElementById('total-items');
    const orderPacksList = document.getElementById('order-packs');
    const ordersTableBody = document.getElementById('orders-table').querySelector('tbody');

    /**
     * Initialize the orders module
     */
    function init() {
        // Load orders on page load
        loadOrders();

        // Add event listeners
        createOrderForm.addEventListener('submit', handleCreateOrder);
    }

    /**
     * Load orders from the API and display them in the table
     */
    async function loadOrders() {
        try {
            const orders = await api.getOrders();
            renderOrdersTable(orders);
        } catch (error) {
            showError('Failed to load orders: ' + error.message);
        }
    }

    /**
     * Render the orders table with the provided orders data
     * @param {Array} orders - Array of order objects
     */
    function renderOrdersTable(orders) {
        // Clear the table
        ordersTableBody.innerHTML = '';

        if (orders.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="4" class="text-center">No orders available</td>';
            ordersTableBody.appendChild(row);
            return;
        }

        // Add each order to the table
        orders.forEach(order => {
            const row = document.createElement('tr');
            
            // Format packs for display
            const packsDisplay = order.packs.map(pack => 
                `${pack.quantity}x${pack.pack.amount}`
            ).join(', ');
            
            row.innerHTML = `
                <td>${order.requestedItems}</td>
                <td>${order.overpackedItems}</td>
                <td>${order.totalItems}</td>
                <td>${packsDisplay}</td>
            `;
            
            ordersTableBody.appendChild(row);
        });
    }

    /**
     * Handle creating a new order
     * @param {Event} event - The form submit event
     */
    async function handleCreateOrder(event) {
        event.preventDefault();
        
        const amount = parseInt(orderAmountInput.value);
        if (!amount || amount <= 0) {
            showError('Please enter a valid amount');
            return;
        }

        try {
            const order = await api.createOrder(amount);
            displayOrderResult(order);
            orderAmountInput.value = '';
            loadOrders();
            showSuccess('Order created successfully');
        } catch (error) {
            showError('Failed to create order: ' + error.message);
        }
    }

    /**
     * Display the order result in the UI
     * @param {Object} order - The order object
     */
    function displayOrderResult(order) {
        // Update order details
        requestedItemsEl.textContent = order.requestedItems;
        overpackedItemsEl.textContent = order.overpackedItems;
        totalItemsEl.textContent = order.totalItems;
        
        // Clear and update packs list
        orderPacksList.innerHTML = '';
        order.packs.forEach(pack => {
            const li = document.createElement('li');
            li.textContent = `${pack.quantity} x ${pack.pack.amount}`;
            orderPacksList.appendChild(li);
        });
        
        // Show the result
        orderResult.classList.remove('hidden');
    }

    /**
     * Show a success message
     * @param {string} message - The success message to display
     */
    function showSuccess(message) {
        alert(message); // Simple alert for now, can be replaced with a better UI component
    }

    /**
     * Show an error message
     * @param {string} message - The error message to display
     */
    function showError(message) {
        alert('Error: ' + message); // Simple alert for now, can be replaced with a better UI component
    }

    // Public API
    return {
        init: init,
        loadOrders: loadOrders
    };
})();