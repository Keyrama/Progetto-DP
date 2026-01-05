// Minimal JavaScript only for date validation (no fetch/AJAX)
document.addEventListener('DOMContentLoaded', () => {
    const dateInput = document.getElementById('date');
    
    if (dateInput) {
        // Set minimum date to today
        const today = new Date();
        const todayStr = today.toISOString().split('T')[0];
        dateInput.setAttribute('min', todayStr);
        
        // Set maximum date to 3 months from now
        const maxDate = new Date();
        maxDate.setMonth(maxDate.getMonth() + 3);
        const maxDateStr = maxDate.toISOString().split('T')[0];
        dateInput.setAttribute('max', maxDateStr);
    }
});