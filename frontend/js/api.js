// API endpoints
const API_ENDPOINTS = {
    packs: `/packs`,
    orders: `/orders`,
    createOrder: `/orders/items`
};

/**
 * Generic function to handle API errors
 * @param {Response} response - The fetch response object
 * @returns {Promise} - Promise that resolves to the response JSON or rejects with an error
 */
async function handleResponse(response) {
    if (!response.ok) {
        const errorData = await response.json().catch(() => ({
            error: 'An unknown error occurred'
        }));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }
    return response.json();
}

/**
 * API client for interacting with the backend
 */
const api = {
    /**
     * Get all packs
     * @returns {Promise<Array>} - Promise that resolves to an array of packs
     */
    getPacks: async function() {
        try {
            const response = await fetch(API_ENDPOINTS.packs);
            return handleResponse(response);
        } catch (error) {
            console.error('Error fetching packs:', error);
            throw error;
        }
    },

    /**
     * Add a new pack
     * @param {number} amount - The amount of the pack
     * @returns {Promise<Object>} - Promise that resolves to the created pack
     */
    addPack: async function(amount) {
        try {
            const response = await fetch(`${API_ENDPOINTS.packs}/${amount}`, {
                method: 'POST'
            });
            return handleResponse(response);
        } catch (error) {
            console.error('Error adding pack:', error);
            throw error;
        }
    },

    /**
     * Update a pack
     * @param {number} oldAmount - The current amount of the pack
     * @param {number} newAmount - The new amount for the pack
     * @returns {Promise<Object>} - Promise that resolves to the updated pack
     */
    updatePack: async function(oldAmount, newAmount) {
        try {
            const response = await fetch(`${API_ENDPOINTS.packs}/${oldAmount}/${newAmount}`, {
                method: 'PUT'
            });
            return handleResponse(response);
        } catch (error) {
            console.error('Error updating pack:', error);
            throw error;
        }
    },

    /**
     * Delete a pack
     * @param {number} amount - The amount of the pack to delete
     * @returns {Promise<void>} - Promise that resolves when the pack is deleted
     */
    deletePack: async function(amount) {
        try {
            const response = await fetch(`${API_ENDPOINTS.packs}/${amount}`, {
                method: 'DELETE'
            });
            if (!response.ok) {
                const errorData = await response.json().catch(() => ({
                    error: 'An unknown error occurred'
                }));
                throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
            }
            return true;
        } catch (error) {
            console.error('Error deleting pack:', error);
            throw error;
        }
    },

    /**
     * Get all orders
     * @returns {Promise<Array>} - Promise that resolves to an array of orders
     */
    getOrders: async function() {
        try {
            const response = await fetch(API_ENDPOINTS.orders);
            return handleResponse(response);
        } catch (error) {
            console.error('Error fetching orders:', error);
            throw error;
        }
    },

    /**
     * Create a new order
     * @param {number} amount - The number of items to order
     * @returns {Promise<Object>} - Promise that resolves to the created order
     */
    createOrder: async function(amount) {
        try {
            const response = await fetch(`${API_ENDPOINTS.createOrder}/${amount}`, {
                method: 'POST'
            });
            return handleResponse(response);
        } catch (error) {
            console.error('Error creating order:', error);
            throw error;
        }
    }
};