/**
 * Packs module for managing pack-related UI functionality
 */
const packsModule = (function() {
    // DOM elements
    const packsTableBody = document.getElementById('packs-table').querySelector('tbody');
    const addPackForm = document.getElementById('add-pack-form');
    const packAmountInput = document.getElementById('pack-amount');
    const updatePackModal = document.getElementById('update-pack-modal');
    const updatePackForm = document.getElementById('update-pack-form');
    const updateOldAmountInput = document.getElementById('update-old-amount');
    const updateNewAmountInput = document.getElementById('update-new-amount');
    const closeModalBtn = updatePackModal.querySelector('.close');

    /**
     * Initialize the packs module
     */
    function init() {
        // Load packs on page load
        loadPacks();

        // Add event listeners
        addPackForm.addEventListener('submit', handleAddPack);
        updatePackForm.addEventListener('submit', handleUpdatePack);
        closeModalBtn.addEventListener('click', closeModal);
        window.addEventListener('click', function(event) {
            if (event.target === updatePackModal) {
                closeModal();
            }
        });
    }

    /**
     * Load packs from the API and display them in the table
     */
    async function loadPacks() {
        try {
            const packs = await api.getPacks();
            renderPacksTable(packs);
        } catch (error) {
            showError('Failed to load packs: ' + error.message);
        }
    }

    /**
     * Render the packs table with the provided packs data
     * @param {Array} packs - Array of pack objects
     */
    function renderPacksTable(packs) {
        // Clear the table
        packsTableBody.innerHTML = '';

        // Sort packs by amount (descending)
        packs.sort((a, b) => b.amount - a.amount);

        if (packs.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="2" class="text-center">No packs available</td>';
            packsTableBody.appendChild(row);
            return;
        }

        // Add each pack to the table
        packs.forEach(pack => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${pack.amount}</td>
                <td>
                    <div class="action-buttons">
                        <button class="btn btn-edit btn-sm edit-pack" data-amount="${pack.amount}">Edit</button>
                        <button class="btn btn-danger btn-sm delete-pack" data-amount="${pack.amount}">Delete</button>
                    </div>
                </td>
            `;

            // Add event listeners for edit and delete buttons
            row.querySelector('.edit-pack').addEventListener('click', () => openUpdateModal(pack.amount));
            row.querySelector('.delete-pack').addEventListener('click', () => handleDeletePack(pack.amount));

            packsTableBody.appendChild(row);
        });
    }

    /**
     * Handle adding a new pack
     * @param {Event} event - The form submit event
     */
    async function handleAddPack(event) {
        event.preventDefault();
        
        const amount = parseInt(packAmountInput.value);
        if (!amount || amount <= 0) {
            showError('Please enter a valid amount');
            return;
        }

        try {
            await api.addPack(amount);
            packAmountInput.value = '';
            loadPacks();
            showSuccess(`Pack with amount ${amount} added successfully`);
        } catch (error) {
            showError('Failed to add pack: ' + error.message);
        }
    }

    /**
     * Open the update pack modal
     * @param {number} amount - The current amount of the pack
     */
    function openUpdateModal(amount) {
        updateOldAmountInput.value = amount;
        updateNewAmountInput.value = amount;
        updatePackModal.style.display = 'block';
    }

    /**
     * Close the update pack modal
     */
    function closeModal() {
        updatePackModal.style.display = 'none';
    }

    /**
     * Handle updating a pack
     * @param {Event} event - The form submit event
     */
    async function handleUpdatePack(event) {
        event.preventDefault();
        
        const oldAmount = parseInt(updateOldAmountInput.value);
        const newAmount = parseInt(updateNewAmountInput.value);
        
        if (!newAmount || newAmount <= 0) {
            showError('Please enter a valid amount');
            return;
        }

        try {
            await api.updatePack(oldAmount, newAmount);
            closeModal();
            loadPacks();
            showSuccess(`Pack updated from ${oldAmount} to ${newAmount} successfully`);
        } catch (error) {
            showError('Failed to update pack: ' + error.message);
        }
    }

    /**
     * Handle deleting a pack
     * @param {number} amount - The amount of the pack to delete
     */
    async function handleDeletePack(amount) {
        if (!confirm(`Are you sure you want to delete the pack with amount ${amount}?`)) {
            return;
        }

        try {
            await api.deletePack(amount);
            loadPacks();
            showSuccess(`Pack with amount ${amount} deleted successfully`);
        } catch (error) {
            showError('Failed to delete pack: ' + error.message);
        }
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
        loadPacks: loadPacks
    };
})();